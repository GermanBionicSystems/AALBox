package main

import (
	"fmt"
	"github.com/clausecker/nfc/v2"
	"log"
)

func main() {
	//TODO: initialize the docker container (postgres database). maybe do docker file. for now start it manually

	/*	control := controller.NewController()
		control.Start()*/
	// Initialize NFC library
	version := nfc.Version()
	fmt.Printf("Using NFC library version: %s\n", version)

	// List available NFC devices
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

	// TODO: Connect to an NFC device and perform read/write operations
	// This will require additional functions and logic, which are not provided in the initial code.
	// For demonstration purposes, I'll show connecting to the first device and attempting a read.

	// Connect to the first NFC device
	device := devices[0]
	fmt.Printf("Connecting to NFC device: %s\n", device)
	// TODO: Use appropriate functions to connect to the device and read data

	// Cleanup and exit
	// TODO: Any cleanup tasks, like closing connections to NFC devices
	fmt.Println("Exiting program.")
}
