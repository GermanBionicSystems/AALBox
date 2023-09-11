package controller

import (
	"AALBox/src/repo"
	"AALBox/src/view"
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
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

func (control *Controller) Start() {

	app := &cli.App{
		Name:  "AALBox",
		Usage: "Manage and play songs",
		Commands: []*cli.Command{
			{
				Name:  "play",
				Usage: "Play a song by word",
				Action: func(c *cli.Context) error {
					word := c.Args().First()
					songPath := control.SongRepo.GetSongPath(word)
					if songPath != "" {
						view.PlaySong(songPath)
					} else {
						fmt.Println("Song not found for the given word.")
					}
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "Add a new song to the database",
				Action: func(c *cli.Context) error {
					word := c.Args().Get(0)
					path := c.Args().Get(1)
					err := control.SongRepo.AddSong(word, path)
					if err != nil {
						fmt.Println("Error adding song:", err)
					} else {
						fmt.Println("Song added successfully!")
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
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
