<p align="center"><img width="330px" src="/_docs/img/mcc.png" alt="mcc"/></p>

#

[![GitHub release](https://img.shields.io/github/release/qmu/mcc.svg)](https://github.com/qmu/mcc/releases)

mcc is a terminal dashboard easily configured by yaml. 

## Install

Availabe on Mac and Linux. Fetch the [latest release](https://github.com/qmu/mcc/releases) for your platform.

#### macOS

```bash
brew tap qmu/mcc
brew install mcc
```

#### Linux

```bash
sudo wget https://github.com/qmu/mcc/releases/download/v0.9.1/linux_amd64_mcc -O /usr/local/bin/mcc
sudo chmod +x /usr/local/bin/mcc
```

## Usage

```bash
# if ./mcc.yml exists, just
mcc

# or give its path
mcc -c path/to/mcc.yml
```

## Config

See the [_examples](https://github.com/qmu/mcc/tree/master/_examples)

## Key Bindings

KeyBinding          | Description
--------------------|---------------------------------------------------------
<kbd>Ctrl + j,k,h,l</kbd> | Switch widgets
<kbd>j, k, ↑, ↓</kbd> | Move cursor in the active widget
<kbd>Enter</kbd> in the Menu widget | Execute a command
<kbd>Ctrl-c, q</kbd> | quit

## License 

**MIT License**

Copyright (c) qmu Co., Inc.
