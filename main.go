package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/term"
)

type direction int8

const (
	up    = direction(-1)
	down  = direction(1)
	right = direction(2)
	left  = direction(-2)
)

type node struct {
	x int
	y int
	d direction
}

type snake struct {
	nodes []*node
	d     direction
}

func (s *snake) head() *node {
	return s.nodes[len(s.nodes)-1]
}

func (s *snake) move() {
	s.head().d = s.d
	for i, n := range s.nodes {
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
		if i < (len(s.nodes) - 1) {
			n.d = s.nodes[i+1].d
		}
	}
}

func (s *snake) redirect(d direction) {
	if (s.head().d + d) != 0 {
		s.d = d
	}
}

func (s *snake) draw() {
	fmt.Printf("\033[2J")
	for _, n := range s.nodes {
		fmt.Printf("\033[%d;%dHX", n.y, n.x)
	}
}

func main() {
	snake := snake{
		nodes: []*node{
			&node{2, 5, right},
			&node{3, 5, right},
			&node{4, 5, right},
			&node{5, 5, right},
			&node{6, 5, right},
			&node{7, 5, right},
			&node{8, 5, right},
			&node{9, 5, right},
		},
		d: right,
	}

	ticker := time.NewTicker(80 * time.Millisecond)
	redirections := make(chan direction)
	exit := make(chan bool)

	go keybordListener(redirections, exit)

	fmt.Print("\033[?25l") // hide cursor
	for {
		select {
		case <-exit:
			fmt.Println("\033[?25h") // show cursor
			fmt.Println("\033[0m")   // reset
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
			fallthrough
		case bytes.Equal(c, []byte{27}): // esc
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
