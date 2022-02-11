package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format 'Question,Answer'")
	timeLimit := flag.Int("limit", 30, "The time limit (seconds) for quiz.")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s\n", *csvFileName))
		os.Exit(1)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse csv file.")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		answer := askQuestion(i, p, timer)
		if answer == "true" {
			correct++
		} else if answer == "timeComplete" {
			fmt.Printf("\nTime's up! Your score is %d of %d!\n", correct, len(problems))
			return
		}

	}

	fmt.Printf("Your score is %d of %d!\n", correct, len(problems))

}

func askQuestion(i int, p problem, timer *time.Timer) string {
	var ret string
	fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
	answerCh := make(chan string)
	go func() {
		var reader = bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answerCh <- answer
	}()
	select {
	case <-timer.C:
		return "timeComplete"
	case answer := <-answerCh:
		if strings.ToLower(strings.TrimSpace(answer)) == p.a {
			ret = "true"
		} else {
			ret = "false"
		}
	}

	return ret
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.ToLower(strings.TrimSpace(line[1])),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
