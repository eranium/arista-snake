# Arista Snake
Play the old school Snake game on your Arista switch. Proof of concept written in Go and tested on the 96 port Arista DCS-7280CR3-96. This is a switch that the Dutch internet exchange [ERA-IX](https://www.era-ix.com/) deploys in major PoPs to facilitate 100G interconnection. A video was posted on [Eranium's LinkedIn](https://www.linkedin.com/posts/eranium-bv_what-do-we-do-on-our-spare-fridays-answering-activity-7351915453620838401-4zce) and after many requests we have decided to open-source this great piece of code. Is the game experience great like the original? No, but it was certainly worth it. Note that this game isn't efficient as this switch uses around 700W idle.

![](snake.gif)

## Introduction
This concept is inspired by the nostalgic snake game "Blockade" made in 1976:
> The original Blockade from 1976 and its many clones are two-player games. Viewed from a top-down perspective, each player controls a "snake" with a fixed starting position. The "head" of the snake continually moves forward, unable to stop, growing ever longer. It must be steered left, right, up, and down to avoid hitting walls and the body of either snake. The player who survives the longest wins. Single-player versions are less prevalent and have one or more snakes controlled by the computer, as in the light cycles segment of the 1982 Tron arcade game.

If you're unfamiliar with this game read more about the history [here](https://en.wikipedia.org/wiki/Snake_(video_game_genre)).

## Quick Start
```shell
# Build
GOOS=linux GOARCH=amd64 make build
# Copy the files
scp bin/arista-snake-linux-amd64 <user>@<switch-ip>:/home/<user>/arista-snake
scp grid.json <user>@<switch-ip>:/home/<user>/grid.json
scp eapiconfig-local.ini <user>@<switch-ip>:/home/<user>/eapiconfig-local.ini
# SSH into switch and run
ssh <user>@<switch-ip>
enable
bash
./arista-snake --config-file eapiconfig-local.ini --connection localhost --grid-definition-file grid.json --difficulty hard
```

## Building
In order to build this project, you need to have Go and `make` installed.

To build, you can use the included Makefile:
```shell
make build
```

This will produce a binary under `bin/` for your current OS/architecture.

If you want to run this program directly on an Arista switch, you might have to cross-compile the binary for a different
architecture. For example, to compile this program for an x86-64 machine:
```shell
GOOS=linux GOARCH=amd64 make build
```

And then just copy the binary to the switch.

## Running
In order to run this program, you need:
- A configuration file containing the host and credentials to connect to, see [this link](https://github.com/aristanetworks/goeapi?tab=readme-ov-file#example-eapiconf-file).
  - Example for a remote connection: [eapiconfig.ini](eapiconfig.ini)
  - Example for a local connection: [eapiconfig-local.ini](eapiconfig-local.ini)
- A JSON file containing the game grid, see [grid.json](grid.json).

Program usage:
```shell
Usage of ./arista-snake:
  -config-file string
    	eAPI configuration file
  -connection string
    	eAPI connection to use
  -difficulty string
    	Game difficulty, available options are 'easy', 'medium' and 'hard' (default "medium")
  -grid-definition-file string
    	JSON file containing the port grid definitions for the switch
```

So for example:
```shell
./arista-snake --config-file eapiconfig.ini --connection arista1 --grid-definition-file grid.json --difficulty hard
```

## License
Mozilla Public License Version 2.0
