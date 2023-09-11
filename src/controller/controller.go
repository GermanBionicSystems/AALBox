package controller

import (
	"AALBox/src/repo"
	"AALBox/src/view"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Controller struct {
	SongRepo repo.Repository
}

func NewController() Controller {
	repository := repo.NewRepository()
	return Controller{
		SongRepo: repository,
	}
}

func (c *Controller) Start() {
	for {
		fmt.Println("1. Enter a word to play a song")
		fmt.Println("2. Add a new song to the database")
		fmt.Print("Choose an option: ")
		choice := getInput()
		switch choice {
		case "1":
			word := getInputWord()
			songPath := c.SongRepo.GetSongPath(word)
			if songPath != "" {
				view.PlaySong(songPath)
			} else {
				view.DisplayMessage("Song not found for the given word.")
			}
		case "2":
			fmt.Print("Enter the word for the song: ")
			word := getInput()
			fmt.Print("Enter the path for the song: ")
			path := getInput()
			err := c.SongRepo.AddSong(word, path)
			if err != nil {
				view.DisplayMessage("Error adding song: " + err.Error())
			} else {
				view.DisplayMessage("Song added successfully!")
			}
		default:
			view.DisplayMessage("Invalid option. Please try again.")
		}
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getInputWord() string {
	fmt.Print("Enter a word: ")
	return getInput()
}
