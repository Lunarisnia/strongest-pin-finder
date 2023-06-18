package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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

// Ease of guess criteria
// +1 if same number are used more than 2 times
// +1 for each more repetition of the same number
// +4 if repetition are used more than once globally
// +4 if the same number are used consecutively
// +5 if its on the common password list

func analyzeEaseOfGuess(testCase TestCase) int {
	splitted := strings.Split(testCase.Target, "")
	var hashMap map[string]int = map[string]int{}
	globalRepetition := 0
	easeOfGuessScore := 0

	for i, numStr := range splitted {
		hashMap[numStr]++
		if hashMap[numStr] > 2 {
			globalRepetition++
			easeOfGuessScore++
		}
		if i > 0 && (numStr == splitted[i-1]) {
			easeOfGuessScore += 4
		}
	}
	if globalRepetition > 1 {
		easeOfGuessScore += globalRepetition * 4
	}
	if _, found := commonPasswordMap[testCase.Target]; found {
		easeOfGuessScore += 5
	}
	return easeOfGuessScore
}

var records [][]string = [][]string{
	{"pin", "crack_from_low", "crack_from_up", "ease_of_guess"},
}

var commonPasswordMap map[string]int = map[string]int{}

func mapCommonPassword() {
	jsonFile, err := os.Open("./commonPasswords.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteVal, _ := ioutil.ReadAll(jsonFile)
	commonPasswords := []string{}
	json.Unmarshal(byteVal, &commonPasswords)
	for _, password := range commonPasswords {
		commonPasswordMap[password] = 1
	}
}

func main() {
	mapCommonPassword()
	createTestBench()
	var crackers sync.WaitGroup
	for i := 0; i < len(testBench); i++ {
		crackers.Add(1)
		go func(index int) {
			defer crackers.Done()
			fromLow := crackFromLow(testBench[index])
			fromUp := crackFromUp(testBench[index])
			easeOfGuess := analyzeEaseOfGuess(testBench[index])
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
