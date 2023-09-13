package controller

import (
	"AALBox/src/controller/chip"
	"AALBox/src/repo"
	"AALBox/src/view"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clausecker/nfc/v2"
	"github.com/urfave/cli/v2"
)

// Controller struct to manage songs and NFC interactions
type Controller struct {
	SongRepo repo.Repository
}

// NewController initializes a new Controller with a song repository
func NewController() Controller {
	return Controller{
		SongRepo: repo.NewRepository(),
	}
}

// Start initializes the CLI app and handles user commands
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
					if err := control.SongRepo.AddSong(word, path); err != nil {
						fmt.Println("Error adding song:", err)
						return err
					}
					fmt.Println("Song added successfully!")
					return nil
				},
			},
			{
				Name:  "start-nfc",
				Usage: "Start nfc listener",
				Action: func(c *cli.Context) error {
					fmt.Printf("Using NFC library version: %s\n", nfc.Version())

					devices, err := nfc.ListDevices()
					if err != nil {
						log.Fatalf("Failed to list NFC devices: %v", err)
					}
					if len(devices) == 0 {
						log.Fatal("No NFC devices found.")
					}

					fmt.Println("Available NFC devices:")
					for _, device := range devices {
						fmt.Println("-", device)
					}

					rfidChannel := make(chan string)
					quitChannel := make(chan os.Signal, 1)

					// Create an abstraction of the Reader, DeviceConnection string is empty -> the library to autodetect reader
					rfidReader := chip.NewTagReader("", rfidChannel, 19)

					// Listen for an RFID/NFC tag in another goroutine
					go rfidReader.ListenForTags()

					for {
						fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
						select {
						case tagId := <-rfidReader.TagChannel:
							fmt.Println("this is your id:", tagId)
							songPath := control.SongRepo.GetSongPath(tagId)
							if songPath != "" {
								view.PlaySong(songPath)
							} else {
								fmt.Println("No song associated with this tag.")
							}
						case <-quitChannel:
							rfidReader.Cleanup()
						default:
							time.Sleep(time.Millisecond * 300)
						}
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
