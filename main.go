package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	width  = 80
	height = 24
)

type Scene struct {
	buffer       [][]rune
	treeOffset   int
	starOffset   int
	trees        []Tree
	stars        []Star
	frameCounter int
}

type Star struct {
	x, y int
	char rune
}

type Tree struct {
	x      int
	height int
	shape  []string
}

func main() {
	scene := NewScene()
	scene.Run()
}

func NewScene() *Scene {
	scene := &Scene{
		buffer: make([][]rune, height),
		trees:  generateStaticTrees(),
		stars:  generateStars(20),
	}

	for i := range scene.buffer {
		scene.buffer[i] = make([]rune, width)
	}

	return scene
}

func (s *Scene) Run() {
	fmt.Print("\033[?25l\033[2J")
	defer fmt.Print("\033[?25h")

	for {
		s.update()
		s.render()
		fmt.Print("\033[H")
		s.display()
		time.Sleep(200 * time.Millisecond)
	}
}

func (s *Scene) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func (s *Scene) update() {
	s.frameCounter++
	s.treeOffset++
	if s.treeOffset > width*4 {
		s.treeOffset = 0
	}

	if s.frameCounter%3 == 0 {
		s.starOffset++
		if s.starOffset > width*6 {
			s.starOffset = 0
		}
	}

	if s.frameCounter%300 == 0 {
		s.stars = generateStars(20)
	}
}

func (s *Scene) render() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			s.buffer[y][x] = ' '
		}
	}

	s.renderNightSky()
	s.renderTrees()
	s.renderCarWindow()
}

func (s *Scene) renderNightSky() {
	for _, star := range s.stars {
		starX := star.x - s.starOffset/3
		if starX < 0 {
			starX += width * 2
		}
		if star.y < height-8 && starX < width-10 && starX > 5 {
			s.buffer[star.y][starX] = star.char
		}
	}
}

func (s *Scene) renderTrees() {
	for _, tree := range s.trees {
		treeX := tree.x - s.treeOffset
		if treeX >= -20 && treeX < width+20 {
			adjustedTree := Tree{
				x:      treeX,
				height: tree.height,
				shape:  tree.shape,
			}
			s.drawTree(adjustedTree)
		}
	}
}

func (s *Scene) drawTree(tree Tree) {
	for i, line := range tree.shape {
		y := tree.height - len(tree.shape) + i + 1
		if y >= 0 && y < height-3 {
			for j, char := range line {
				x := tree.x + j
				if x >= 5 && x < width-5 && char != ' ' {
					s.buffer[y][x] = char
				}
			}
		}
	}
}

func (s *Scene) renderCarWindow() {
	windowTop := 2
	windowBottom := height - 4
	windowLeft := 3
	windowRight := width - 4

	for y := windowTop; y <= windowBottom; y++ {
		if y == windowTop || y == windowBottom {
			for x := windowLeft; x <= windowRight; x++ {
				s.buffer[y][x] = '═'
			}
		} else {
			s.buffer[y][windowLeft] = '║'
			s.buffer[y][windowRight] = '║'
		}
	}

	s.buffer[windowTop][windowLeft] = '╔'
	s.buffer[windowTop][windowRight] = '╗'
	s.buffer[windowBottom][windowLeft] = '╚'
	s.buffer[windowBottom][windowRight] = '╝'

	dashboardY := height - 3
	for x := 0; x < width; x++ {
		s.buffer[dashboardY][x] = '▓'
	}

	for x := 0; x < windowLeft; x++ {
		for y := 0; y < height-3; y++ {
			s.buffer[y][x] = '█'
		}
	}
	for x := windowRight + 1; x < width; x++ {
		for y := 0; y < height-3; y++ {
			s.buffer[y][x] = '█'
		}
	}
}

func (s *Scene) display() {
	var sb strings.Builder
	sb.Grow(width * height * 4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sb.WriteRune(s.buffer[y][x])
		}
		if y < height-1 {
			sb.WriteRune('\n')
		}
	}
	fmt.Print(sb.String())
}

func generateStaticTrees() []Tree {
	trees := make([]Tree, 0)
	treeSpacing := 35
	groundLevel := height - 6

	for i := 0; i < 20; i++ {
		x := i * treeSpacing
		tree := generateTree(x, groundLevel, i)
		trees = append(trees, tree)
	}
	return trees
}

func generateStars(count int) []Star {
	stars := make([]Star, count)
	starChars := []rune{'*', '·', '✦', '◦'}

	for i := 0; i < count; i++ {
		stars[i] = Star{
			x:    rand.Intn(width*2) + 5,
			y:    rand.Intn(height-12) + 2,
			char: starChars[rand.Intn(len(starChars))],
		}
	}
	return stars
}

func generateTree(x, groundY, seed int) Tree {
	rand.Seed(int64(seed * 17 + x*3))

	treeType := rand.Intn(5)
	treeHeight := rand.Intn(4) + 4

	var shape []string

	switch treeType {
	case 0:
		if treeHeight <= 5 {
			shape = []string{
				"  ▲  ",
				" ▲▲▲ ",
				"▲▲▲▲▲",
				"  ║  ",
			}
		} else if treeHeight <= 6 {
			shape = []string{
				"   ▲   ",
				"  ▲▲▲  ",
				" ▲▲▲▲▲ ",
				"▲▲▲▲▲▲▲",
				"   ║   ",
				"   ║   ",
			}
		} else {
			shape = []string{
				"    ▲    ",
				"   ▲▲▲   ",
				"  ▲▲▲▲▲  ",
				" ▲▲▲▲▲▲▲ ",
				"▲▲▲▲▲▲▲▲▲",
				"    ║    ",
				"    ║    ",
			}
		}
	case 1:
		if treeHeight <= 5 {
			shape = []string{
				" ♠♠♠ ",
				"♠♠♠♠♠",
				" ♠♠♠ ",
				"  ║  ",
			}
		} else {
			shape = []string{
				"   ♠♠♠   ",
				"  ♠♠♠♠♠  ",
				" ♠♠♠♠♠♠♠ ",
				"♠♠♠♠♠♠♠♠♠",
				"   ♠♠♠   ",
				"    ║    ",
				"    ║    ",
			}
		}
	case 2:
		if treeHeight <= 5 {
			shape = []string{
				" ◆◆◆ ",
				"◆◆◆◆◆",
				" ◆◆◆ ",
				"  ║  ",
			}
		} else {
			shape = []string{
				"  ◆◆◆◆◆  ",
				" ◆◆◆◆◆◆◆ ",
				"◆◆◆◆◆◆◆◆◆",
				" ◆◆◆◆◆◆◆ ",
				"  ◆◆◆◆◆  ",
				"    ║    ",
				"    ║    ",
			}
		}
	case 3:
		if treeHeight <= 5 {
			shape = []string{
				" ▼▼▼ ",
				"▼▼▼▼▼",
				"▼▼▼▼▼",
				"  ║  ",
			}
		} else {
			shape = []string{
				"  ▼▼▼▼▼  ",
				" ▼▼▼▼▼▼▼ ",
				"▼▼▼▼▼▼▼▼▼",
				"▼▼▼▼▼▼▼▼▼",
				"    ║    ",
				"    ║    ",
			}
		}
	default:
		if treeHeight <= 5 {
			shape = []string{
				" ◊◊◊ ",
				"◊◊◊◊◊",
				"◊◊◊◊◊",
				" ◊◊◊ ",
				"  ║  ",
			}
		} else {
			shape = []string{
				"   ◊◊◊   ",
				"  ◊◊◊◊◊  ",
				" ◊◊◊◊◊◊◊ ",
				"◊◊◊◊◊◊◊◊◊",
				" ◊◊◊◊◊◊◊ ",
				"   ◊◊◊   ",
				"    ║    ",
			}
		}
	}

	return Tree{
		x:      x,
		height: groundY,
		shape:  shape,
	}
}
