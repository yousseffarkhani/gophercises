package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Quiz struct {
	problems []Problem
}

type Problem struct {
	question, answer string
}

func main() {
	CSVFilename := flag.String("CSV", "./test.csv", "CSV Filename")
	timerLimit := flag.Int("timer", 2, "Timer limit")
	flag.Parse()
	var correctAnswers int
	quiz := CreateQuiz(*CSVFilename)
	fmt.Println("Press enter to start the game")
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(*timerLimit) * time.Second)
	AskQuestions(quiz, &correctAnswers, timer)
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateQuiz(filename string) *Quiz {
	file, err := ioutil.ReadFile(filename)
	check(err)
	csv := csv.NewReader(strings.NewReader(string(file)))
	records, err := csv.ReadAll()
	check(err)
	var quiz Quiz
	for _, problem := range records {
		quiz.problems = append(quiz.problems, Problem{
			question: problem[0],
			answer:   problem[1],
		})
	}
	return &quiz
}

func AskQuestions(quiz *Quiz, score *int, timer *time.Timer) {
problemLoop:
	for i, problem := range quiz.problems {
		fmt.Printf("Question n° %d : %s\n", i, problem.question)
		answerCh := make(chan string)
		go func() {
			var userInput string
			fmt.Scanln(&userInput)
			answerCh <- userInput
		}()
		select {
		case <-timer.C:
			break problemLoop
		case answer := <-answerCh:
			fmt.Println(answer)
			if assertCorrectAnswer(problem.answer, answer) {
				*score++
				fmt.Printf("Good Answer. Your score is %d / %d\n", *score, len(quiz.problems))
			} else {
				fmt.Printf("Wrong Answer. Your score is %d / %d\n", *score, len(quiz.problems))
			}
		}
	}
	fmt.Printf("\nFinal result : Your score is %d / %d\n", *score, len(quiz.problems))
}

func assertCorrectAnswer(answer, input string) bool {
	if answer == input {
		return true
	}
	return false
}
