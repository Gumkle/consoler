package logger

import (
	"fmt"
	"strings"
)

type loading struct {
	name     string
	progress float32
	logger   *logger
}

func (loading *loading) SetProgress(progress float32) {
	loading.progress = progress
}

type logger struct {
	errorPrepend   string
	infoPrepend    string
	warningPrepend string
	successPrepend string
	loadings       []*loading
	infoChannel    chan string
}

func NewLogger() *logger {
	var loggerInstance *logger = new(logger)
	loggerInstance.infoChannel = make(chan string)
	loggerInstance.SetErrorPrepend("\033[38;2;255;50;50m[ERROR]\033[0m")
	loggerInstance.SetInfoPrepend("\033[38;2;100;100;255m[INFO]\033[0m")
	loggerInstance.SetSuccessPrepend("\033[38;2;50;255;50m[SUCCESS]\033[0m")
	loggerInstance.SetWarningPrepend("\033[38;2;255;255;50m[WARNING]\033[0m")
	go loggerInstance.processInput()
	return loggerInstance
}

func (logger *logger) processInput() {
	for {
		text := <-logger.infoChannel
		fmt.Printf("\033[K%s\n", text)
		if len(logger.loadings) > 0 {
			fmt.Print("\033[K====================\n")
			for _, loading := range logger.loadings {
				fmt.Printf("Loading %s: %s", loading.name, strings.Repeat("#", int(loading.progress*100)))
			}
			fmt.Print(strings.Repeat("\033[1A", len(logger.loadings)+1))
			fmt.Print("\r")
		}
	}
}

func (logger *logger) SetErrorPrepend(prepend string) {
	logger.errorPrepend = prepend
}

func (logger *logger) SetInfoPrepend(prepend string) {
	logger.infoPrepend = prepend
}

func (logger *logger) SetWarningPrepend(prepend string) {
	logger.warningPrepend = prepend
}

func (logger *logger) SetSuccessPrepend(prepend string) {
	logger.successPrepend = prepend
}

func (logger *logger) PrintError(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.errorPrepend, message)
}

func (logger *logger) PrintInfo(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.infoPrepend, message)
}

func (logger *logger) PrintWarning(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.warningPrepend, message)
}

func (logger *logger) PrintSuccess(message string) {
	logger.infoChannel <- fmt.Sprintf("%s %s", logger.successPrepend, message)
}

func (logger *logger) NewLoading(title string) *loading {
	loading := new(loading)
	loading.name = title
	loading.progress = 0
	loading.logger = logger
	logger.loadings = append(logger.loadings, loading)
	return loading
}
