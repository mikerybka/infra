package main

import (
	"fmt"
	"time"
)

func main() {
	count := 0
	for {
		fmt.Println(count)
		count++
		time.Sleep(1 * time.Hour)
		runJobs()
	}
}

func runJobs() {
}
