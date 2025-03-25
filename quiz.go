package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type QA struct {
	question   string
	goodAnswer string
}

func askQuestion(qa QA) bool {

	var proposedAnswer string
	_, err := fmt.Scan(&proposedAnswer)
	if err != nil {
		log.Println("invalid answer")
		return false
	}
	if proposedAnswer == qa.goodAnswer {
		return true
	}

	return false
}

func parseCSV(data []byte) []QA {

	// attempt parsing the CSV file and fill a slice of QA
	r := csv.NewReader(bytes.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln(err.Error() + " --> malformed csv, abort.")
	}

	var qa = make([]QA, len(records))

	// sanitize unecessary spaces in csv fields
	for i := range qa {
		qa[i].question = strings.TrimSpace(records[i][0])
		qa[i].goodAnswer = strings.TrimSpace(records[i][1])

	}

	return qa

}
func main() {

	var fileName string
	var timeLimit int

	// read flags
	flag.StringVar(&fileName, "csv", "problems.csv", "a csv file in the format of 'question, answer'")
	flag.IntVar(&timeLimit, "limit", 30, "time limit in seconds")
	flag.Parse()

	// open csv file
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// parse csv file and sanitize records
	qa := parseCSV(data)

	// run the quizz
	goodAnswerCount := 0
	for i, q := range qa {
		fmt.Printf("Question %d:	%s ?", i, q.question)
		ok := askQuestion(q)

		if ok {
			goodAnswerCount++
		}
	}
	fmt.Printf("Correct Answers/Total %d/%d \n", goodAnswerCount, len(qa))
}
