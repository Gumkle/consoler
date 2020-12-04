package main

import (
	"consoler/consoler"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test := []int{0,1,2,3,4,5}
	fmt.Println(test[2:])
	logger := consoler.NewLogger()
	go randomProgress(logger.NewTask("proces numer 1"))
	go randomProgress(logger.NewTask("proces numer 3"))
	go randomProgress(logger.NewTask("proces numer 2"))
	logger.PrintWarning("No warning")
	logger.PrintInfo("NO info")
	logger.PrintError("No error")
	logger.PrintSuccess("No success")
	go func() {
		time.Sleep(3 * time.Second)
		go randomProgress(logger.NewTask("proces numer 4"))
		time.Sleep(1 * time.Second)
		logger.PrintWarning("Uwaga zaraz dopierdoli")
		go randomProgress(logger.NewTask("proces numer 5"))
		time.Sleep(1 * time.Second)
		go randomProgress(logger.NewTask("proces numer 6"))
		go randomProgress(logger.NewTask("proces numer 7"))
		go randomProgress(logger.NewTask("proces numer 8"))
		go randomProgress(logger.NewTask("proces numer 9"))
		go randomProgress(logger.NewTask("proces numer 10"))
		go randomProgress(logger.NewTask("proces numer 11"))
		go randomProgress(logger.NewTask("proces numer 12"))
		go randomProgress(logger.NewTask("proces numer 13"))
		go randomProgress(logger.NewTask("proces numer 14"))
	}()

	var input string
	fmt.Scanln(&input)
}

func randomProgress(task *consoler.Task) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	time.Sleep(time.Duration(r1.Intn(5)) * time.Second)
	if r1.Intn(10) < 1 {
		task.SetFailed()
	} else {
		task.SetDone()
	}
}
