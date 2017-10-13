package main

import (
	"bufio"
	"flag"
	"io"
	"net"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	host     = flag.String("host", "127.0.0.1:23", "Host on which to connect")
	wasdMode = flag.Bool("wasd", false, "Map arrow keys to wasd values")

	oldState *terminal.State
	err      error
)

const (
	backspace = 127
	arrowLeft = 1000 + iota
	arrowRight
	arrowUp
	arrowDown
	delKey
	homeKey
	endKey
	pageUp
	pageDown
)

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *host)
	if err != nil {
		panic(err)
	}

	rawMode()
	defer endRawMode()

	go paint(conn)

	for {
		keyPress := readKey()
		if keyPress == ctrlKey('q') {
			return
		}

		keySend := make([]byte, 2)
		keySend[1] = '\n'

		keySend[0] = byte(keyPress)

		if *wasdMode {
			switch keyPress {
			case arrowUp:
				keySend[0] = 'w'
			case arrowDown:
				keySend[0] = 's'
			case arrowLeft:
				keySend[0] = 'a'
			case arrowRight:
				keySend[0] = 'd'
			}
		}

		conn.Write(keySend)
	}
}

func paint(b io.Reader) {
	print("\033[H\033[2J")

	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		print(scanner.Text())
		print("\r\n")
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
