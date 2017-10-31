package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/term"
)

type direction uint8

const (
	up    = direction(5)
	down  = direction(6)
	right = direction(7)
	left  = direction(8)
)

type node struct {
	x int
	y int
	d direction
}

func (n *node) move() {
	switch n.d {
	case up:
		n.y--
	case down:
		n.y++
	case left:
		n.x--
	case right:
		n.x++
	}
}

type snake struct {
	nodes []*node
}

func (s *snake) move() {
	for _, n := range s.nodes {
		n.move()
	}
}

func (s *snake) redirect(d direction) {
	s.nodes[0].d = d
}

func (s *snake) draw() {
	for _, n := range s.nodes {
		fmt.Printf("\033[2J\033[%d;%dHX", n.y, n.x)
	}
}

func getch() []byte {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	t.Restore()
	t.Close()
	if err != nil {
		return nil
	}
	return bytes[0:numRead]
}

func main() {
	snake := snake{
		nodes: []*node{&node{5, 5, right}},
	}

	ticker := time.NewTicker(80 * time.Millisecond)
	redirections := make(chan direction)
	exit := make(chan bool)

	go keybordListener(redirections, exit)

	for {
		select {
		case <-exit:
			fmt.Println()
			return
		case d := <-redirections:
			snake.redirect(d)
		case <-ticker.C:
			snake.move()
			snake.draw()
		default:

		}
	}
}

func keybordListener(redirections chan direction, exit chan bool) {
	for {
		c := getch()
		switch {
		case bytes.Equal(c, []byte{3}): // ctrl + c
			exit <- true
		case bytes.Equal(c, []byte{27, 91, 65}): // up
			redirections <- up
		case bytes.Equal(c, []byte{27, 91, 66}): // down
			redirections <- down
		case bytes.Equal(c, []byte{27, 91, 67}): // right
			redirections <- right
		case bytes.Equal(c, []byte{27, 91, 68}): // left
			redirections <- left
		default:
			fmt.Println("\nUnknown pressed", c)
		}
	}
}
