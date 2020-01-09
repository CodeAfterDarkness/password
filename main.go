package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

type WordScore struct {
	Word  string
	Score float64 `json:"-"`
}

var scores []string

const maxLength int = 14

func dothething(wordIn string, scoreChan chan WordScore) {

	for _, word := range words {

		if len(wordIn+word) == maxLength {
			newScore := WordScore{
				Word:  wordIn + word,
				Score: wordCheck(wordIn + word),
			}
			scoreChan <- newScore
		}

		if len(wordIn+word) > maxLength {
			continue
		}

		if len(wordIn+word) < maxLength {
			dothething(wordIn+word, scoreChan)
		}
	}

}

func main() {
	log.Print("Starting")

	initializePhyKeys()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Printf("Words: %d", len(words))

	for i := len(words) - 1; i >= 0; i-- {
		if wordCheck(words[i]) < float64(len(words[i])) {
			words = append(words[:i], words[i+1:]...)
		}
	}

	log.Printf("Words: %d", len(words))

	scoreChan := make(chan WordScore, 10)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, word := range words[:int(len(words)/4)] {
			fmt.Printf("1: %s\n", word)
			dothething(word, scoreChan)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, word := range words[int(len(words)/4):int(len(words)/2)] {
			fmt.Printf("2: %s\n", word)
			dothething(word, scoreChan)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, word := range words[int(len(words)/2) : int(len(words)/4)*3] {
			fmt.Printf("3: %s\n", word)
			dothething(word, scoreChan)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, word := range words[int(len(words)/4)*3:] {
			fmt.Printf("4: %s\n", word)
			dothething(word, scoreChan)

		}

	}()

	t := time.NewTicker(time.Second * 10)

	go func() {
		highScoresFound := 0
		for {
			select {
			case <-t.C:
				log.Printf("Found %d high scores so far", highScoresFound)
			case newScore := <-scoreChan:

				// if highScoresFound%1000000 == 0 {
				// 	log.Printf("Found '%s' score: %.1f", newScore.Word, newScore.Score)
				// }

				if newScore.Score > 13 {
					scores = append(scores, newScore.Word)
					highScoresFound++
				}

				if highScoresFound > 30000000 {

					jsonBytes, err := json.MarshalIndent(scores, "", "\t")
					if err != nil {
						log.Print(err)
					} else {
						err := ioutil.WriteFile("scores.json", jsonBytes, 0644)
						if err != nil {
							log.Print(err)
						}
					}

					os.Exit(0)
				}

			}
		}
	}()

	wg.Wait()

	// var highScores []WordScore
	// for _, score := range scores {
	// 	//log.Printf("Score: '%s': %f", score.Word, score.Score)
	// 	if score.Score == 14 && len(score.Word) == 14 {
	// 		highScores = append(highScores, score)
	// 	}
	// }

	// for _, hs := range highScores {
	// 	log.Printf("Found '%s' score %.0f", hs.Word, hs.Score)
	// }
	log.Print("Finished")
}

//autofail if false
func fingerCheck(testword string) int {
	for i := 0; i < len(testword)-1; i++ {
		key1 := phyKeys[string(testword[i])]
		key2 := phyKeys[string(testword[i+1])]

		//log.Printf("Finger1 %d, Finger2 %d", key1.Finger, key2.Finger)
		if key1.Finger == key2.Finger {
			return 0
		}
	}
	return 1
}

func handCheck(testword string) (factor float64) {
	factor = float64(len(testword))
	for i := 0; i < len(testword)-1; i++ {
		key1 := phyKeys[string(testword[i])]
		key2 := phyKeys[string(testword[i+1])]
		if key1.Hand == key2.Hand {
			factor *= 0.5
		}
	}
	return factor
}

func wordCheck(testword string) float64 {
	return handCheck(testword) * float64(fingerCheck(testword))
}
