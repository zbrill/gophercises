package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	fileName, timerFlag := getArgs()
	duration := time.Duration(timerFlag) * time.Second

	fmt.Print("Press any key to begin the timer...")
	fmt.Scanln()

	timer := time.NewTimer(duration)

	file := openFile(fileName)
	problems := parseFile(file)
	rightCount := 0
	var answer string
	for _, problem := range problems {
		result := askQuestion(problem, timer.C, answer)
		if result == -1 {
			break
		}
		rightCount += result
	}
	fmt.Println("You got", rightCount, "out of", len(problems), "questions right!")
}

func getArgs() (string, int) {
	fileName := flag.String("csv", "problems.csv", "name of your quiz csv file")
	timer := flag.Int("timer", 0, "how long you have to complete the quiz")
	flag.Parse()
	return *fileName, *timer
}

func openFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error while opening the file", err)
	}
	return file
}

func parseFile(file *os.File) []Problem {
	reader := csv.NewReader(file)
	problems, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Couldn't read all rows of csv", err)
	}
	problemsSlice := make([]Problem, len(problems))
	for i, problem := range problems {
		problemsSlice[i] = Problem{problem[0], problem[1]}
	}
	return problemsSlice
}

func askQuestion(question Problem, t <-chan time.Time, answer string) int {
	select {
	case <-t:
		fmt.Println("You ran out of time!")
		return -1
	default:
		fmt.Println(question.question)
		fmt.Print("Answer: ")
		fmt.Scanln(&answer)
		if answer == strings.TrimSpace(question.answer) {
			return 1
		}
	}
	return 0
}
