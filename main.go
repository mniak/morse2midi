package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"flag"

	"github.com/algoGuy/EasyMIDI/smf"
	"github.com/algoGuy/EasyMIDI/smfio"
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

type event struct {
	Delta    int
	Duration int
}

func morse2events(morseCode string, norests bool) <-chan event {
	out := make(chan event)
	morseCode = strings.ReplaceAll(morseCode, " / ", "/")
	go func() {
		log.Printf("Evaluating morse code: %v", morseCode)
		var lastnote, spaces int
		for i, x := range morseCode {
			if norests {
				switch x {
				case '/':
					spaces = 7 - 1 - 1
				case ' ':
					spaces = 3 - 1 - 1
				case '.':
					out <- event{
						Delta:    i - lastnote + spaces,
						Duration: 2,
					}
					lastnote = i + 1
					spaces = 0
				case '-':
					out <- event{
						Delta:    i - lastnote + spaces,
						Duration: 4,
					}
					lastnote = i + 1
					spaces = 0
				}
			} else {
				switch x {
				case '/':
					spaces = 7 - 1
				case ' ':
					spaces = 3 - 1
				case '.':
					out <- event{
						Delta:    i - lastnote + spaces,
						Duration: 1,
					}
					lastnote = i + 1
					spaces = 1
				case '-':
					out <- event{
						Delta:    i - lastnote + spaces,
						Duration: 3,
					}
					lastnote = i + 1
					spaces = 1
				}
			}
		}
		close(out)
	}()
	return out
}
func parseArgs() (filepath, text string, norests bool) {
	pathPtr := flag.String("o", "morse.mid", "Where to save the midi file")
	restsPtr := flag.Bool("r", false, "Put rests between notes")
	flag.Parse()

	text = strings.Join(flag.Args(), " ")

	return *pathPtr, text, !*restsPtr
}

func main() {
	filepath, text, norests := parseArgs()

	m, err := text2morse(text)
	checkErr(err)
	events := morse2events(m, norests)

	// Create division
	ppqn := 960
	division, err := smf.NewDivision(uint16(ppqn), smf.NOSMTPE)
	checkErr(err)

	// Create new midi struct
	midi, err := smf.NewSMF(smf.Format0, *division)
	checkErr(err)

	// Create track struct
	track := &smf.Track{}

	// Add track to new midi struct
	err = midi.AddTrack(track)
	checkErr(err)

	const note = 60
	const velocity = 100
	baseDuration := ppqn / 2

	// Create some midi and meta events
	for event := range events {
		log.Println(event)
		startDelta := uint32(event.Delta * baseDuration)
		start, err := smf.NewMIDIEvent(startDelta, smf.NoteOnStatus, 0x00, note, velocity)
		checkErr(err)
		err = track.AddEvent(start)
		checkErr(err)

		endDelta := uint32(event.Duration * baseDuration)
		end, err := smf.NewMIDIEvent(endDelta, smf.NoteOffStatus, 0x00, note, 0)
		checkErr(err)
		err = track.AddEvent(end)
		checkErr(err)
	}
	metaEventOne, err := smf.NewMetaEvent(0, smf.MetaEndOfTrack, []byte{})
	checkErr(err)
	err = track.AddEvent(metaEventOne)
	checkErr(err)

	// Save to new midi source file
	outputMidi, err := os.Create(filepath)
	checkErr(err)
	defer outputMidi.Close()

	// Create buffering stream
	writer := bufio.NewWriter(outputMidi)
	smfio.Write(writer, midi)
	writer.Flush()
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
