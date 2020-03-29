package main

import (
	"fmt"
	"strings"

	"github.com/alwindoss/morse"
)

func text2morse(text string) (string, error) {
	h := morse.NewHacker()
	morseCode, err := h.Encode(strings.NewReader(text))
	if err != nil {
		return "", err
	}
	return string(morseCode), nil
}

func morse2sequence(text string) (result []bool) {
	words := strings.Split(text, "/")
	for iword, word := range words {
		chars := strings.Fields(word)
		for ichar, char := range chars {
			for isignal, signal := range char {
				switch signal {
				case '.':
					result = append(result, true)
				case '-':
					result = append(result, true, true, true)
				}
				if isignal < len(char)-1 {
					result = append(result, false)
				}
			}
			if ichar < len(chars)-1 {
				// End of char
				result = append(result, false, false, false)
			}
		}
		if iword < len(words)-1 {
			// End of word
			result = append(result, false, false, false, false, false, false, false)
		}
	}
	fmt.Println()
	return
}

type event struct {
	When    int
	HowLong int
}

func morse2events(sequence []bool) (events []event) {
	var on bool
	var start int
	for i, s := range sequence {
		switch {
		case s && !on: //start
			start = i
			on = true
		case !s && on: //stop
			events = append(events, event{
				When:    start,
				HowLong: i - start,
			})
			on = false
		}
	}
	return
}

func main() {
	m, _ := text2morse("MORSE CODE")
	seq := morse2sequence(m)
	for _, sig := range seq {
		if sig {
			fmt.Print("+")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Println("===.===...===.===.===...=.===.=...=.=.=...=.......===.=.===.=...===.===.===...===.=.=...=")

	events := morse2events(seq)
	for _, e := range events {
		fmt.Println("Event", e)
	}
}
