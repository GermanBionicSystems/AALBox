package view

import (
	"fmt"
	"os/exec"
)

func DisplayMessage(message string) {
	fmt.Println(message)
}

func PlaySong(songPath string) {
	cmd := exec.Command("mpg123", songPath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error playing song:", err)
	}
}
