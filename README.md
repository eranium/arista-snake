# arista-snake
Play the old school Snake game on your Arista switch.

## TLDR
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

To build, you can use the Makefile:
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
