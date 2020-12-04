package consoler

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Loading struct {
	name     string
	progress float32
	logger   *Logger
	done bool
}

func (loading *Loading) SetProgress(progress float32) {
	loading.progress = progress
	if progress >= float32(1.0) && !loading.done {
		loading.logger.PrintSuccess(fmt.Sprintf("Loading %s finished", loading.name))
		loading.logger.RemoveLoading(loading)
		loading.done = true
	}
}

type Logger struct {
	errorPrepend   string
	infoPrepend    string
	warningPrepend string
	successPrepend string
	loadings       []*Loading
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
	for {
		time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
		select {
		 case text := <-logger.infoChannel:
			 fmt.Printf("\033[K%s\n\u001B[K", text)
		 default:
			 if len(logger.loadings) > 0 {
				 fmt.Print("====================\n")
				 for _, loading := range logger.loadings {
					 fmt.Printf("Loading %s: %.2f\n", loading.name, loading.progress)
				 }
				 fmt.Print(strings.Repeat("\033[1A", len(logger.loadings)+1))
				 fmt.Print("\r")
			 }
		}
	}
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

func (logger *Logger) NewLoading(title string) *Loading {
	loading := new(Loading)
	loading.name = title
	loading.progress = 0
	loading.logger = logger
	logger.loadings = append(logger.loadings, loading)
	return loading
}



func (logger *Logger) RemoveLoading(searchedLoading *Loading) error{
	var lock sync.Mutex
	lock.Lock()
	indexOfSearchedLoading := -1
	for index, loading := range logger.loadings {
		if loading == searchedLoading {
			indexOfSearchedLoading = index
			break
		}
	}
	if indexOfSearchedLoading == -1 {
		return errors.New("No such loading in logger")
	}
	logger.loadings[indexOfSearchedLoading] = logger.loadings[len(logger.loadings)-1]
	logger.loadings = logger.loadings[:len(logger.loadings)-1]
	lock.Unlock()
	return nil
}
