# AALBox

AALBox is a music management system for Raspberry pi that allows users to add, and remove to/form database, play songs
using NFC interactions. The project is structured into different packages, each responsible for a specific
functionality.

## Features

- **Play Songs**: Play a song associated with a specific word or NFC tag.
- **Manage Songs**: Add or remove songs from the database.
- **NFC Integration**: Listen for NFC tags and play the associated song.

## Getting Started

### Prerequisites

- Go (version 1.16 or higher)
- Docker
- PostgreSQL Docker Image
- mpg123 (for playing songs)
- NFC Reader compatible with the [`clausecker/nfc/v2`](https://github.com/clausecker/nfc) library

### Setting up a PN532 NFC module on a Raspberry Pi using I2C

Before starting the application, you need to set up the PN532 NFC module on your Raspberry Pi. Follow the steps from
this [article](https://blog.stigok.com/2017/10/12/setting-up-a-pn532-nfc-module-on-a-raspberry-pi-using-i2c.html) to get
it up and running.

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/adelbables/AALBox.git

2. **Navigate to the project directory**:
   ```bash
   cd AALBox

3. **Install mpg123**:

   ```bash
    sudo apt-get install mpg123

4. **Set up the PostgreSQL Docker container**:

   ```bash
   docker run --name postgres -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 -d postgres

5. **Install the required Go packages**:

   ```bash
    cd src && go build -o AALBox

## Usage

Use the CLI to manage and play songs:

- **To play a song**:
   ```bash
    AALBox play <tag_id>

- **To add a new song to the database**:

   ```bash
  AALBox add <tag_id> <path_to_song>

- **To remove a song from the database**:

   ```bash
    AALBox remove <tag_id>

- **To start the NFC listener**:

   ```bash
    AALBox start-nfc

## Packages

- controller: Manages songs and NFC interactions.
- repo: Handles database interactions for song management.
- view: Provides functionalities to play songs.
- chip: Manages NFC interactions.

## Troubleshooting

Low Sound Volume on Raspberry Pi Zero with USB Audio Card

If you are using a Raspberry Pi Zero and a USB audio card, you might experience issues with low sound volume. To resolve
this, follow the steps below:

- Install alsamixer:

   ```bash
    sudo apt-get install alsamixer

- Adjust the Sound Volume:
  Open alsamixer:

  ```bash
    alsamixer

Use the arrow keys to adjust the volume level to your preference.

- Save Your Settings:
  After adjusting the volume, save your settings to ensure they persist after a reboot:

    ```bash
    sudo alsactl store
  ```

### ps:
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY
