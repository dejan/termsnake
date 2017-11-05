package main

import (
	"time"

	"github.com/nsf/termbox-go"
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
	termbox.Clear(termbox.ColorDefault, termbox.Attribute(1))
	for _, n := range s.nodes {
		termbox.SetCell(n.x, n.y, ' ', termbox.ColorDefault, termbox.Attribute(2))
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	snake := newSnake(5, right)
	ticker := time.NewTicker(70 * time.Millisecond)
	events := make(chan termbox.Event)

	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case <-ticker.C:
			snake.move()
			snake.draw()
		case ch := <-events:
			switch ch.Key {
			case termbox.KeyEsc:
				fallthrough
			case termbox.KeyCtrlC:
				return
			case termbox.KeyArrowUp:
				snake.redirect(up)
			case termbox.KeyArrowDown:
				snake.redirect(down)
			case termbox.KeyArrowRight:
				snake.redirect(right)
			case termbox.KeyArrowLeft:
				snake.redirect(left)
			}
		default:
		}
	}
}

func newSnake(size int, d direction) snake {
	const (
		startX = 2
		startY = 5
	)
	var nodes = make([]*node, size)
	for i := range nodes {
		nodes[i] = &node{startX + i, startY, right}
	}
	return snake{nodes, d}
}
