package view

import (
	"fmt"
	"os"
	"os/exec"
)

func DisplayMessage(message string) {
	fmt.Println(message)
}

func PlaySong(songPath string) {
	cmd := exec.Command("mpg123", songPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error playing song:", err)
	}
}
