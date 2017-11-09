package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type direction int8

type gameState int8

const (
	up    = direction(-1)
	down  = direction(1)
	right = direction(2)
	left  = direction(-2)

	welcome  = gameState(0)
	playing  = gameState(1)
	gameOver = gameState(2)
	exit     = gameState(3)
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

func (s *snake) tail() []*node {
	return s.nodes[0 : len(s.nodes)-1]
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

type food struct {
	x int
	y int
}

type game struct {
	snake  *snake
	ticker *time.Ticker
	events chan termbox.Event
	state  gameState
	food   *food
}

func (g *game) tick() {
	if g.state == playing {
		g.snake.move()
		g.consolidate()
	}
}

func (g *game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.Attribute(1))
	for _, n := range g.snake.nodes {
		termbox.SetCell(n.x, n.y, ' ', termbox.ColorDefault, termbox.Attribute(2))
	}

	termbox.SetCell(g.food.x, g.food.y, ' ', termbox.ColorDefault, termbox.Attribute(4))
	termbox.Flush()
}

func freeSpot(nodes []*node) (int, int) {
	sx, sy := termbox.Size()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := r.Intn(sx)
	y := r.Intn(sy)
	for _, n := range nodes {
		if (n.x == x) && (n.y == y) {
			return freeSpot(nodes)
		}
	}
	return x, y
}

func (g *game) consolidate() {
	head := g.snake.head()
	for _, n := range g.snake.tail() {
		if (n.x == head.x) && (n.y == head.y) {
			g.state = gameOver
		}
	}

	if (head.x == g.food.x) && (head.y == g.food.y) {
		g.snake.potential = 10
		x, y := freeSpot(g.snake.nodes)
		g.food.x = x
		g.food.y = y
	}
}

func (g *game) loop() {
	for {

		switch g.state {

		case welcome:
			fmt.Println("Welcome. Press space to continue...")
			select {
			case ch := <-g.events:
				switch ch.Key {
				case termbox.KeyEsc:
					fallthrough
				case termbox.KeyCtrlC:
					return
				case termbox.KeySpace:
					g.state = playing
				}
			}

		case playing:
			select {
			case <-g.ticker.C:
				g.tick()
				g.draw()
			case ch := <-g.events:
				switch ch.Key {
				case termbox.KeyEsc:
					fallthrough
				case termbox.KeyCtrlC:
					g.state = exit
				case termbox.KeyArrowUp:
					g.snake.redirect(up)
				case termbox.KeyArrowDown:
					g.snake.redirect(down)
				case termbox.KeyArrowRight:
					g.snake.redirect(right)
				case termbox.KeyArrowLeft:
					g.snake.redirect(left)
				}
			default:
			}

		case gameOver:
			fmt.Println("Game Over. Press space to start again...")
			select {
			case ch := <-g.events:
				switch ch.Key {
				case termbox.KeyEsc:
					fallthrough
				case termbox.KeyCtrlC:
					return
				case termbox.KeySpace:
					g.snake = newSnake(1, right)
					g.state = playing
				}
			}
		case exit:
			return
		}
	}
}

func newSnake(size, d direction) *snake {
	const (
		initX         = 5
		initY         = 5
		initPotential = 20
	)
	var nodes = make([]*node, size)
	for i := range nodes {
		nodes[i] = &node{initX + i, initY, right}
	}
	return &snake{nodes: nodes, potential: initPotential, d: d}
}

func newGame(events chan termbox.Event) game {
	return game{
		snake:  newSnake(1, right),
		state:  welcome,
		ticker: time.NewTicker(70 * time.Millisecond),
		events: events,
		food:   &food{20, 20},
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	game := newGame(events)
	game.loop()
}
