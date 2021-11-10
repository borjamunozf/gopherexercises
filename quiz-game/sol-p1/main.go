package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Answer struct {
	id       int
	response string
	input    string
	valid    bool
}

func main() {

	var filename string
	flag.StringVar(&filename, "CSV Quiz filename", "problems.csv", "Specify full filename for Quiz problems.")
	flag.Parse()

	fmt.Println("Filename specified: ", filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	var bingo bool
	var answers []Answer
	for i, record := range records {
		fmt.Println(record[0])
		var userAnswer string
		fmt.Scanln(&userAnswer)
		validateInput(userAnswer)
		if userAnswer == record[1] {
			bingo = true
		} else {
			bingo = false
		}
		answer := Answer{id: i, response: record[1], input: userAnswer, valid: bingo}
		answers = append(answers, answer)
	}

	fmt.Println(len(answers))
	var correct, fail int
	correct = 0
	fail = 0
	for _, ans := range answers {
		if ans.valid == true {
			correct += 1
		} else {
			fail += 1
			fmt.Printf("fail %v", fail)
		}

	}

	fmt.Printf("Total correct answers: %d", correct)
	fmt.Printf("Total failed answers: %d", fail)
}

func validateInput(input string) string {
	strings.TrimSpace(input)
	return input
}
