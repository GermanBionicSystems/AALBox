// Package controller provides functionalities to manage songs and NFC interactions.
package controller

import (
	"AALBox/src/controller/chip"
	"AALBox/src/repo"
	"AALBox/src/view"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
					tagId := c.Args().First()
					songPath := control.SongRepo.GetSongPath(tagId)
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
					tagId := c.Args().Get(0)
					path := c.Args().Get(1)
					if err := control.SongRepo.AddSong(tagId, path); err != nil {
						fmt.Println("Error adding song:", err)
						return err
					}
					fmt.Println("Song added successfully!")
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Delete a song from the database",
				Action: func(c *cli.Context) error {
					tagId := c.Args().Get(0)
					if err := control.SongRepo.RemoveSong(tagId); err != nil {
						fmt.Println("Error deleting song:", err)
						return err
					}
					fmt.Println("Song removed successfully!")
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

					// Ensure the sounds folder exists in the user's home directory and move the file if necessary
					homeSoundsDir := filepath.Join(os.Getenv("HOME"), "sounds")
					if _, err := os.Stat(homeSoundsDir); os.IsNotExist(err) {
						if err := os.Mkdir(homeSoundsDir, os.ModePerm); err != nil {
							log.Fatalf("Failed to create directory in home: %s", err)
						}
					}

					// Define the source and destination paths for the MP3 file
					defaultSrcPath := "../sounds/wat-wer-bist-du-denn.mp3"
					defaultDestPath := filepath.Join(homeSoundsDir, "wat-wer-bist-du-denn.mp3")

					// Check if the file already exists in the destination; if not, copy it
					if _, err := os.Stat(defaultDestPath); os.IsNotExist(err) {
						if err := copyFile(defaultSrcPath, defaultDestPath); err != nil {
							log.Printf("Failed to move file: %s", err)
						}
					}

					for {
						select {
						case tagId := <-rfidReader.TagChannel:
							fmt.Println("This is your id:", tagId)
							songPath := control.SongRepo.GetSongPath(tagId)
							if songPath != "" {
								go view.PlaySong(songPath)
							} else {
								// Play the default song from the home directory if no specific song is found
								go view.PlaySong(defaultDestPath)
								fmt.Println("No song associated with this tag, playing default sound.")
							}
						case <-quitChannel:
							rfidReader.Cleanup()
						default:
							time.Sleep(time.Millisecond * 10)
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

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}
