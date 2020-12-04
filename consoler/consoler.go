package consoler

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Task struct {
	name     string
	logger   *Logger
	done bool
}

func (task *Task) SetDone() {
	if !task.done {
		task.logger.PrintSuccess(fmt.Sprintf("Finished: %s", task.name))
		err := task.logger.RemoveTask(task)
		if err != nil {
			task.logger.PrintError(fmt.Sprintf("Task %s already done!", task.name))
		}
		task.done = true
	}
}

type Logger struct {
	errorPrepend   string
	infoPrepend    string
	warningPrepend string
	successPrepend string
	tasks          []*Task
	infoChannel    chan string
}

func NewLogger() *Logger {
	var loggerInstance *Logger = new(Logger)
	loggerInstance.infoChannel = make(chan string)
	loggerInstance.SetErrorPrepend("\033[38;2;255;50;50m[ERROR]\033[0m")
	loggerInstance.SetInfoPrepend("\033[38;2;100;100;255m[INFO]\033[0m")
	loggerInstance.SetSuccessPrepend("\033[38;2;50;255;50m[SUCCESS]\033[0m")
	loggerInstance.SetWarningPrepend("\033[38;2;255;255;50m[WARNING]\033[0m")
	go loggerInstance.processInput()
	return loggerInstance
}

func (logger *Logger) processInput() {
	go setupCloseCleanup(logger)
	for {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		select {
		 case text := <-logger.infoChannel:
			 fmt.Printf("\033[K%s\n\u001B[K", text)
		 default:
			 if len(logger.tasks) > 0 {
				 fmt.Print("===================================================================================\n")
				 for _, task := range logger.tasks {
					 fmt.Printf("\033[K%s...\n", task.name)
				 }
				 fmt.Printf("\033[%dA", len(logger.tasks)+1)
				 fmt.Print("\r")
			 }
		}
	}
}

func setupCloseCleanup(logger *Logger) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print(strings.Repeat("\n", len(logger.tasks)+1))
		os.Exit(0)
	}()
}

func (logger *Logger) SetErrorPrepend(prepend string) {
	logger.errorPrepend = prepend
}

func (logger *Logger) SetInfoPrepend(prepend string) {
	logger.infoPrepend = prepend
}

func (logger *Logger) SetWarningPrepend(prepend string) {
	logger.warningPrepend = prepend
}

func (logger *Logger) SetSuccessPrepend(prepend string) {
	logger.successPrepend = prepend
}

func (logger *Logger) PrintError(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.errorPrepend, message)
}

func (logger *Logger) PrintInfo(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.infoPrepend, message)
}

func (logger *Logger) PrintWarning(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.warningPrepend, message)
}

func (logger *Logger) PrintSuccess(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.successPrepend, message)
}

func (logger *Logger) NewTask(title string) *Task {
	task := new(Task)
	task.name = title
	task.logger = logger
	task.done = false
	logger.tasks = append(logger.tasks, task)
	return task
}

func (logger *Logger) RemoveTask(searchedTask *Task) error{
	var lock sync.Mutex
	lock.Lock()
	indexOfSearchedTask := -1
	for index, task := range logger.tasks {
		if task == searchedTask {
			indexOfSearchedTask = index
			break
		}
	}
	if indexOfSearchedTask == -1 {
		return errors.New("No such task in logger")
	}
	logger.tasks[indexOfSearchedTask] = logger.tasks[len(logger.tasks)-1]
	logger.tasks = logger.tasks[:len(logger.tasks)-1]
	lock.Unlock()
	return nil
}
