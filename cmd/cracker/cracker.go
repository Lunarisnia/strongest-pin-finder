package main

import (
	"fmt"
	"sync"
	"time"
)

type TestCase struct {
	Target string
}

var testBench []TestCase

func createTestBench() {
	for i := 0; i <= 999999; i++ {
		testBench = append(testBench, TestCase{
			Target: fmt.Sprintf("%06d", i),
		})
	}
}

func crackFromLow(testCase TestCase) {
	start := time.Now()
	guess := "000000"
	attemptCount := 0
	for i := 0; i <= 999999; i++ {
		if guess == testCase.Target {
			break
		}
		attemptCount++
		guess = fmt.Sprintf("%06d", i)
	}
	elapsed := time.Since(start)
	elapsed = elapsed + (time.Duration(attemptCount/3*5) * time.Second)
	fmt.Printf("(%v) found in: %v\n", testCase.Target, elapsed)
}

func main() {
	createTestBench()
	var crackers sync.WaitGroup
	for i := 0; i < len(testBench); i++ {
		crackers.Add(1)
		go func(index int) {
			defer crackers.Done()
			crackFromLow(testBench[index])
		}(i)
	}
	crackers.Wait()
}
