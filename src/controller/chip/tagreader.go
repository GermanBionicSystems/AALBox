// Package chip provides functionalities to interact with NFC (Near Field Communication) devices.
package chip

import (
	"encoding/hex"
	"github.com/clausecker/nfc/v2"
	"github.com/warthog618/gpiod"
	"log"
	"time"
)

// TagReader represents an NFC tag reader.
type TagReader struct {
	// TagChannel is a channel to send detected tag UIDs.
	TagChannel chan string

	// reader is the NFC device instance.
	reader *nfc.Device

	// ResetPin is the GPIO pin used for resetting the NFC reader.
	ResetPin int

	// DeviceConnection is the connection string for the NFC device.
	DeviceConnection string
}

// NewTagReader creates a new TagReader instance.
// It takes in a device connection string, a channel for tag UIDs, and a reset pin.
func NewTagReader(deviceConnection string, tagChannel chan string, resetPin int) *TagReader {
	return &TagReader{
		DeviceConnection: deviceConnection,
		TagChannel:       tagChannel,
		ResetPin:         resetPin,
	}
}

// init initializes the NFC reader. If there's an error in communication with the device, it resets the reader.
func (reader *TagReader) init() {
	dev, err := nfc.Open(reader.DeviceConnection)
	if err != nil {
		reader.Reset()
		log.Printf("Cannot communicate with the device: %s \n", err)
		return
	}
	reader.reader = &dev
	if err := reader.reader.InitiatorInit(); err != nil {
		log.Fatal("Failed to initialize")
	}
}

// Reset implements the hardware reset by pulling the ResetPin low and then releasing it.
// This is achieved using the GPIO library.
func (reader *TagReader) Reset() {
	log.Println("Resetting the reader..")
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	pin, err := c.RequestLine(reader.ResetPin, gpiod.AsOutput(0))
	if err != nil {
		log.Println(err)
		return
	}
	defer pin.Close()

	if err := pin.SetValue(1); err != nil {
		log.Println(err)
	}
	time.Sleep(400 * time.Millisecond)

	if err := pin.SetValue(0); err != nil {
		log.Println(err)
	}
	time.Sleep(400 * time.Millisecond)

	if err := pin.SetValue(1); err != nil {
		log.Println(err)
	}
	time.Sleep(100 * time.Millisecond)
}

// Cleanup closes the NFC reader connection.
func (reader *TagReader) Cleanup() {
	if err := reader.reader.Close(); err != nil {
		log.Println("Error closing NFC reader:", err)
	}
}

// ListenForTags initializes the reader and then continuously listens for NFC tags.
// When a tag is detected, its UID is sent to the TagChannel.
func (reader *TagReader) ListenForTags() {

	//Initialize the reader
	reader.init()

	//Listen for all the modulations specified
	var modulations = []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106},
		{Type: nfc.Felica, BaudRate: nfc.Nbr212},
		{Type: nfc.Felica, BaudRate: nfc.Nbr424},
		{Type: nfc.Jewel, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106},
	}

	for {
		// Poll for 300ms
		tagCount, target, err := reader.reader.InitiatorPollTarget(modulations, 1, 300*time.Millisecond)
		if err != nil {
			log.Println("Error polling the reader:", err)
			continue
		}
		// Check if any tag was detected
		if tagCount > 0 {
			UID := extractUID(target)
			// Send the UID of the tag to controller goroutine
			if UID != "" {
				reader.TagChannel <- UID
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// extractUID extracts the UID from the given NFC target.
// It supports various NFC card types, including ISO14443a, ISO14443b, Felica, Jewel, and ISO14443biClass.
func extractUID(target nfc.Target) string {
	var UID string
	switch target.Modulation() {
	case nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}:
		var card = target.(*nfc.ISO14443aTarget)
		var UIDLen = card.UIDLen
		var ID = card.UID
		UID = hex.EncodeToString(ID[:UIDLen])
		break
	case nfc.Modulation{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106}:
		var card = target.(*nfc.ISO14443bTarget)
		var UIDLen = len(card.ApplicationData)
		var ID = card.ApplicationData
		UID = hex.EncodeToString(ID[:UIDLen])
		break
	case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr212}:
		var card = target.(*nfc.FelicaTarget)
		var UIDLen = card.Len
		var ID = card.ID
		UID = hex.EncodeToString(ID[:UIDLen])
		break
	case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr424}:
		var card = target.(*nfc.FelicaTarget)
		var UIDLen = card.Len
		var ID = card.ID
		UID = hex.EncodeToString(ID[:UIDLen])
		break
	case nfc.Modulation{Type: nfc.Jewel, BaudRate: nfc.Nbr106}:
		var card = target.(*nfc.JewelTarget)
		var ID = card.ID
		var UIDLen = len(ID)
		UID = hex.EncodeToString(ID[:UIDLen])
		break
	case nfc.Modulation{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106}:
		var card = target.(*nfc.ISO14443biClassTarget)
		var ID = card.UID
		var UIDLen = len(ID)
		UID = hex.EncodeToString(ID[:UIDLen])
		break
	}
	return UID
}
