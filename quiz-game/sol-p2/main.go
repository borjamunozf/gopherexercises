package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {

	filename := flag.String("csv", "problems.csv", "Specifies the csv file with a format 'Question,Answer'")
	limit := flag.Int("limit", 30, "Specifies seconds limit to finish quiz.")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV File: %s\n", *filename))
	}

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse CSV")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)

	correct := 0
problemloop:
	for i, p := range problems {
		fmt.Printf("Question number #%d: %s = \n", i+1, p.q)

		//doing channel, scanf is blocking
		//anonymous goroutine
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You scored %d correct answers of %d questions\n", correct, len(problems))
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d correct answers of %d questions\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}
func exit(message string) {
	fmt.Errorf(message)
	os.Exit(-1)
}
