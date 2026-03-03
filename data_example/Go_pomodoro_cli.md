---
title: Go_pomodoro_cli
description: A Golang pomodoro cli tool built to get a better understanding of go routines
label: building
url: https://github.com/zidariu-sabin/Go_pomodoro_cli
updatedat: "2025-08-15T13:35:08Z"
languages:
- name: Go
  color: "#00ADD8"
  size: 98.6
- name: Makefile
  color: "#427819"
  size: 1.4
docfileoid: 1392b31
---
# Go pomodoro cli tool

Description: 
- A simple GO project that gives the option to start a sequence of pomodoro timers using the terminal and see the current progress
Benefits gained building the project:
- better understanding of golang routines and text parsing

### How to start

- Clone the github repository
git clone https://github.com/zidariu-sabin/Go_pomodoro_cli.git

- Build the project 
git mod init
git build 
- Start a pomodoro sequence
```shell
pomodoro
```
The default is set is set to 4 25 min work- 5 min break rounds

### Usage
flags available:
  -b int
        duration of the break time (default 5)
  -r int
        number of rounds the user wants the timer to run for (default 4)
  -w int
        duration of the working time (default 5)
```shell
pomodoro -w 1 -b 1 -r 2
```
    will start a sequence of 2 rounds with 1 minute of work and 1 minute of break time