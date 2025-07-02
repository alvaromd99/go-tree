# Go Tree Clone

A simple command-line tool written in Go that lists files and directories in a tree-like format, similar to the Unix `tree` command.

## Purpose

This project was created as a learning exercise to practice and improve Go programming skills, including working with filesystems, recursion, and terminal output formatting.

## Features

- Displays directory structure with branches and indentation.
- Colors directories and executable files for easy identification.
- Shows counts of total directories and files at the end.

## Requirements

- Go 1.13 or higher installed on your system.

### Build the project

```bash
go build -o gotree main.go
```

### Running

```bash
./gotree [directory]
```
