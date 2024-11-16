package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func backup() {
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

	fmt.Println("Listening for key pressess....")
	for event := range reader.Read() {
		if event.Type == keylogger.EvKey && event.KeyPress() {
			fmt.Println("key Pressed:", event.KeyString())

			go playSound(soundFile)
		}
	}
}

func playSound(filepath string) {
	f, err := os.Open(filepath)
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
		log.Println("Failed to Initialize speaker:", err)
	}

	speaker.Play(streamer)

	select {}
}
