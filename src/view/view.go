// Package view provides functionalities to play and stop songs using the mpg123 command.
package view

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// currentCmd holds the currently playing song's command.
// mu is a mutex to ensure safe concurrent access to currentCmd.
var (
	currentCmd *exec.Cmd
	mu         sync.Mutex
)

// PlaySong plays the song located at the provided songPath.
// If a song is already playing, it stops the current song and starts the new one.
// It uses the mpg123 command to play the song.
func PlaySong(songPath string) {
	StopSong()

	mu.Lock()
	currentCmd = exec.Command("mpg123", songPath)
	currentCmd.Stdout = os.Stdout
	currentCmd.Stderr = os.Stderr
	mu.Unlock()

	err := currentCmd.Start()
	if err != nil {
		fmt.Println("Error playing song:", err)
		return
	}
	// Wait for the command to finish
	if err := currentCmd.Wait(); err != nil {
		fmt.Println("Error playing song:", err)
	}
}

// StopSong stops the currently playing song.
// If no song is playing, it does nothing.
func StopSong() {
	mu.Lock()
	defer mu.Unlock()
	if currentCmd != nil && currentCmd.Process != nil {
		if err := currentCmd.Process.Kill(); err != nil {
			fmt.Println("Failed to stop song:", err)
		}
		currentCmd = nil
	}

}
