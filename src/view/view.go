package view

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	currentCmd *exec.Cmd
	mu         sync.Mutex
)

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
