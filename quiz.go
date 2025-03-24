package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type QA struct {
	question      string
	goodAnswer    string
	questionIndex int
}

func askQuestion(qa QA) (bool, error) {
	fmt.Printf("Question %d:	%s ?", qa.questionIndex, qa.question)
	var proposedAnswer string
	_, err := fmt.Scan(&proposedAnswer)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if proposedAnswer == qa.goodAnswer {

		return true, nil
	}

	return false, nil
}


func main() {

	var fileName string
	var timeLimit int

	// read flags
	flag.StringVar(&fileName, "csv", "problems.csv", "a csv file in the format of 'question, answer'")
	flag.IntVar(&timeLimit, "limit", 30, "time limit in seconds")
	flag.Parse()
	

	// read csv file
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(data))
	questionCount := 0
	goodAnswerCount := 0

	for {
		record, err := r.Read()

		// end of file check
		if err == io.EOF {
			break
		}

		//invalid CSV record
		if err != nil {
			log.Print(err)
		}
		// increment question counter
		questionCount++
		// fill question/answer struct
		qa := QA{record[0], record[1], questionCount}

		// ask question and validate answer

		answer, err := askQuestion(qa); 
		
		if err != nil {
			log.Println("invalid answer, doesn't count")
			continue
		}
		if answer {
			goodAnswerCount++
			fmt.Println("Good answer")
			continue
		}
		
		//neither error nor good answer
		fmt.Println("Bad answer")

	}
	fmt.Printf("Good Answers/Total %d/%d \n", goodAnswerCount, questionCount)
}

