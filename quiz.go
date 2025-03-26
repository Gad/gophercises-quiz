package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
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

func shuffleQuestions(qa *[]QA) {

	rand.Shuffle(len(*qa), func(i, j int) {
		(*qa)[i], (*qa)[j] = (*qa)[j], (*qa)[i]
	})

}
func main() {

	var fileName *string
	var timeLimit *int
	var shuffle *bool

	// read flags
	fileName = flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit = flag.Int("limit", 30, "time limit in seconds")
	shuffle = flag.Bool("shuffle", false, "shuffle questions from questions set")
	flag.Parse()

	// open csv file
	data, err := os.ReadFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	// parse csv file and sanitize records
	qa := parseCSV(data)

	if *shuffle {
		shuffleQuestions(&qa)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeLimit)*time.Second)
	defer cancel()
	goodAnswerCount := 0
	done := make(chan bool)

	go runQuiz(qa, &goodAnswerCount, done)

	select {
	case <-ctx.Done():
		fmt.Println("time out !")
		fmt.Printf("Correct Answers/Total %d/%d \n", goodAnswerCount, len(qa))
	case <-done:
		fmt.Printf("Correct Answers/Total %d/%d \n", goodAnswerCount, len(qa))
	}
}

func runQuiz(qa []QA, goodAnswerCount *int, done chan<- bool) {

	for i, q := range qa {

		fmt.Printf("Question %d:	%s ?", i+1, q.question)
		ok := askQuestion(q)

		if ok {
			*goodAnswerCount++
		}
	}
	done <- true
}
