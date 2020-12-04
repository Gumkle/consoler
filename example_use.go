package main

import (
	"console_helper/consoler"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	logger := consoler.NewLogger()
	go randomProgress(logger.NewLoading("ladowanie1"), logger)
	go randomProgress(logger.NewLoading("ladowanie3"), logger)
	go randomProgress(logger.NewLoading("ladowanie2"), logger)
	logger.PrintWarning("No warning")
	logger.PrintInfo("NO info")
	logger.PrintError("No error")
	logger.PrintSuccess("No success")

	go func() {
		time.Sleep(3 * time.Second)
		go randomProgress(logger.NewLoading("ladowanie4"), logger)
		time.Sleep(1 * time.Second)
		go randomProgress(logger.NewLoading("ladowanie5"), logger)
	}()

	var input string
	fmt.Scanln(&input)
}

func randomProgress(loading *consoler.Loading, logger *consoler.Logger) {
	for progress := float32(0.0); progress <= float32(1.1); progress += float32(0.05) {
		loading.SetProgress(progress)
		time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
	}
}
