package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer func(fd *os.File) {
		err = fd.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	reader := csv.NewReader(file)
	records, err2 := reader.ReadAll()
	if err2 != nil {
		log.Fatal(err2)
	}

	totalQuestions := len(records)
	answeredCorrect := 0

	after := time.After(10 * time.Second)

	go func() {
		for _, record := range records {
			fmt.Println("Type the answer for:", record[0])
			var givenAnswer int
			_, _ = fmt.Scan(&givenAnswer)
			correctAnswer, err3 := strconv.Atoi(record[1])
			if err3 != nil {
				log.Fatal(err3)
			}

			if givenAnswer == correctAnswer {
				answeredCorrect++
			}
		}
	}()

	select {
	case <-after:
		fmt.Println("Time is over!")
	}

	fmt.Println("You answered", answeredCorrect, "questions out of", totalQuestions, "correctly.")
}
