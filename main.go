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
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
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

	snake := newSnake()
	ticker := time.NewTicker(80 * time.Millisecond)
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

func newSnake() snake {
	return snake{
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
}
