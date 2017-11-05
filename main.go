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
	nodes     []*node
	d         direction
	potential int
}

func (s *snake) head() *node {
	return s.nodes[len(s.nodes)-1]
}

func (s *snake) move() {
	head := s.head()
	head.d = s.d

	var newHead *node
	if s.potential > 0 {
		s.potential--
		newHead = &node{}
		s.nodes = append(s.nodes, newHead)
	} else {
		newHead = s.nodes[0]
		s.nodes = append(s.nodes[1:], newHead)
	}

	newHead.d = s.d
	newHead.x = head.x
	newHead.y = head.y
	switch newHead.d {
	case up:
		newHead.y--
	case down:
		newHead.y++
	case left:
		newHead.x--
	case right:
		newHead.x++
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

	snake := newSnake(1, right)
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

func newSnake(size, d direction) snake {
	const (
		initX         = 5
		initY         = 5
		initPotential = 5
	)
	var nodes = make([]*node, size)
	for i := range nodes {
		nodes[i] = &node{initX + i, initY, right}
	}
	return snake{nodes: nodes, potential: initPotential, d: d}
}
