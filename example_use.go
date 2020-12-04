package main

import (
	"fmt"
	"helion_prices_monitor/logger"
	"math/rand"
	"time"
)

func main() {
	logg := logger.NewLogger()
	go randomProgress(logg.NewLoading("ladowanie1"), logg)
	go randomProgress(logg.NewLoading("ladowanie3"), logg)
	go randomProgress(logg.NewLoading("ladowanie2"), logg)
	logg.PrintWarning("No warning")
	logg.PrintInfo("NO info")
	logg.PrintError("No error")
	logg.PrintSuccess("No success")

	var input string
	fmt.Scanln(&input)
}

func randomProgress(loading, logg) {
	for initialProgress := float32(0.0); initialProgress < 100; initialProgress += float32(0.5) {
		loading.SetProgress(initialProgress)
		time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
	}
	logg.PrintSuccess("Finished")
}
