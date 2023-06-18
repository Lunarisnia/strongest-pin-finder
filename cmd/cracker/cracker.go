package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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

func crackFromUp(testCase TestCase) time.Duration {
	start := time.Now()
	guess := "999999"
	attemptCount := 0
	for i := 999999; i >= 0; i-- {
		if guess == testCase.Target {
			break
		}
		attemptCount++
		guess = fmt.Sprintf("%06d", i)
	}
	elapsed := time.Since(start)
	elapsed = elapsed + (time.Duration(attemptCount/3*5) * time.Second)
	return elapsed
}
func crackFromLow(testCase TestCase) time.Duration {
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
	return elapsed
}

var records [][]string = [][]string{
	{"pin", "crack_from_low", "crack_from_up", "ease_of_guess"},
}

func main() {
	createTestBench()
	var crackers sync.WaitGroup
	for i := 0; i < 10; i++ {
		crackers.Add(1)
		go func(index int) {
			defer crackers.Done()
			fromLow := crackFromLow(testBench[index])
			fromUp := crackFromUp(testBench[index])
			easeOfGuess := 0
			records = append(records, []string{fmt.Sprintf("\"%v\"", testBench[index].Target), strconv.FormatInt(fromLow.Milliseconds(), 10), strconv.FormatInt(fromUp.Milliseconds(), 10), fmt.Sprint(easeOfGuess)})
		}(i)
	}
	crackers.Wait()

	f, err := os.Create("reports.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
