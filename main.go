package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var mu sync.Mutex // To synchronize access to the speaker

func main() {
	keyboard := keylogger.FindKeyboardDevice()
	if keyboard == "" {
		log.Fatal("No Keyboard device found :/")
	}
	fmt.Println("Using device:", keyboard)

	reader, err := keylogger.New(keyboard)
	if err != nil {
		log.Fatal("Failed to open keyboard device:", err.Error())
	}
	defer reader.Close()

	soundFile := "mech-keyboard-02-102918.mp3"

	// Initialize speaker globally
	f, err := os.Open(soundFile)
	if err != nil {
		log.Fatal("failed to open audio file:", err.Error())
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal("Failed to decode audio file:", err.Error())
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal("Failed to initialize speaker:", err)
	}

	// Create a buffer to replay the audio
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	fmt.Println("Listening for key presses...")
	for event := range reader.Read() {
		if event.Type == keylogger.EvKey && event.KeyPress() {
			fmt.Println("Key Pressed:", event.KeyString())

			// Play the buffered sound
			go func() {
				mu.Lock()
				defer mu.Unlock()

				speaker.Play(buffer.Streamer(0, buffer.Len()))
			}()
		}
	}
}
