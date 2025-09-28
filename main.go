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
	bushes       []Bush
	stars        []Star
	ninja        Ninja
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

type Bush struct {
	x      int
	height int
	shape  []string
}

type Ninja struct {
	x          int
	y          int
	velocityX  float64
	velocityY  float64
	onGround   bool
	state      NinjaState
	animFrame  int
	targetTree int
}

type NinjaState int

const (
	Running NinjaState = iota
	Jumping
	Landing
	Crouching
)

func main() {
	scene := NewScene()
	scene.Run()
}

func NewScene() *Scene {
	scene := &Scene{
		buffer: make([][]rune, height),
		trees:  generateStaticTrees(),
		bushes: generateStaticBushes(),
		stars:  generateStars(20),
		ninja: Ninja{
			x:          30,
			y:          height - 8,
			velocityX:  0.5,
			velocityY:  0,
			onGround:   true,
			state:      Running,
			animFrame:  0,
			targetTree: 0,
		},
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
		time.Sleep(120 * time.Millisecond)
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

	s.updateNinja()
}

func (s *Scene) render() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			s.buffer[y][x] = ' '
		}
	}

	s.renderNightSky()
	s.renderTrees()
	s.renderBushes()
	s.renderCarWindow()
	s.renderNinja()
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
	groundLevel := height - 6

	x := 50 // Start position
	for i := 0; i < 20; i++ {
		tree := generateTree(x, groundLevel, i)
		trees = append(trees, tree)

		// Variable spacing between trees
		spacing := 25 + rand.Intn(20) // 25-45 units apart
		x += spacing
	}
	return trees
}

func generateStaticBushes() []Bush {
	bushes := make([]Bush, 0)
	groundLevel := height - 6

	for i := 0; i < 30; i++ {
		x := rand.Intn(700) + 10
		if rand.Intn(100) < 25 {
			bush := generateBush(x, groundLevel, i)
			bushes = append(bushes, bush)
		}
	}
	return bushes
}

func (s *Scene) renderBushes() {
	for _, bush := range s.bushes {
		bushX := bush.x - s.treeOffset
		if bushX >= -10 && bushX < width+10 {
			adjustedBush := Bush{
				x:      bushX,
				height: bush.height,
				shape:  bush.shape,
			}
			s.drawBush(adjustedBush)
		}
	}
}

func (s *Scene) drawBush(bush Bush) {
	for i, line := range bush.shape {
		y := bush.height - len(bush.shape) + i + 1
		if y >= 0 && y < height-3 {
			for j, char := range line {
				x := bush.x + j
				if x >= 5 && x < width-5 && char != ' ' {
					s.buffer[y][x] = char
				}
			}
		}
	}
}

func generateBush(x, groundY, seed int) Bush {
	rand.Seed(int64(seed * 23 + x*7))

	bushType := rand.Intn(3)
	var shape []string

	switch bushType {
	case 0:
		shape = []string{
			" *** ",
			"*****",
			" *** ",
		}
	case 1:
		shape = []string{
			"  ~~~  ",
			" ~~~~~ ",
			"~~~~~~~",
			" ~~~~~ ",
		}
	default:
		shape = []string{
			" ooo ",
			"ooooo",
		}
	}

	return Bush{
		x:      x,
		height: groundY,
		shape:  shape,
	}
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

func (s *Scene) updateNinja() {
	gravity := 0.3
	groundLevel := float64(height - 6)

	s.ninja.animFrame++

	switch s.ninja.state {
	case Running:
		// Ninja stays in fixed position - only check for trees when they're very close
		nearbyTree := s.findNearbyTree()
		if nearbyTree != nil {
			s.ninja.state = Jumping
			s.ninja.velocityY = -3.2  // Moderate jump height
			s.ninja.onGround = false
		}

	case Jumping:
		s.ninja.velocityY += gravity
		s.ninja.y += int(s.ninja.velocityY)

		// Land when reaching ground level
		if float64(s.ninja.y) >= groundLevel-2 {
			s.ninja.y = int(groundLevel) - 2
			s.ninja.state = Landing
			s.ninja.onGround = true
			s.ninja.velocityY = 0
		}

	case Landing:
		if s.ninja.animFrame%4 == 0 {
			s.ninja.state = Running
		}
	}

	// Keep ninja completely stationary horizontally
	s.ninja.x = 25
}

func (s *Scene) findNextTree() *Tree {
	for _, tree := range s.trees {
		treeX := tree.x - s.treeOffset
		if treeX > s.ninja.x {
			return &tree
		}
	}
	return nil
}

func (s *Scene) findNearbyTree() *Tree {
	for _, tree := range s.trees {
		treeX := tree.x - s.treeOffset
		// Only jump when tree is very close - right before collision
		if treeX > s.ninja.x+2 && treeX < s.ninja.x+8 && treeX > 10 && treeX < width-10 {
			return &tree
		}
	}
	return nil
}

func (s *Scene) renderNinja() {
	// Only render ninja within the car window bounds
	windowLeft := 6
	windowRight := width - 7
	windowTop := 4
	windowBottom := height - 5

	if s.ninja.x < windowLeft || s.ninja.x > windowRight || s.ninja.y < windowTop || s.ninja.y > windowBottom {
		return
	}

	// Only render on empty space - never overwrite anything
	if s.buffer[s.ninja.y][s.ninja.x] != ' ' {
		return
	}

	var sprite rune
	switch s.ninja.state {
	case Running:
		if (s.ninja.animFrame/8)%2 == 0 {
			sprite = 'R'
		} else {
			sprite = 'r'
		}
	case Jumping:
		sprite = 'J'
	case Landing:
		sprite = 'L'
	default:
		sprite = 'R'
	}

	// Simply replace the empty space with ninja sprite
	s.buffer[s.ninja.y][s.ninja.x] = sprite
}

func generateTree(x, groundY, seed int) Tree {
	rand.Seed(int64(seed * 17 + x*3))

	treeType := rand.Intn(4)
	treeHeight := rand.Intn(3) + 5

	var shape []string

	switch treeType {
	case 0: // Pine tree
		if treeHeight <= 5 {
			shape = []string{
				"   ^   ",
				"  ^^^  ",
				" ^^^^^ ",
				"^^^^^^^",
				"   |   ",
				"   |   ",
			}
		} else if treeHeight <= 6 {
			shape = []string{
				"    ^    ",
				"   ^^^   ",
				"  ^^^^^  ",
				" ^^^^^^^ ",
				"^^^^^^^^^",
				"    |    ",
				"    |    ",
			}
		} else {
			shape = []string{
				"     ^     ",
				"    ^^^    ",
				"   ^^^^^   ",
				"  ^^^^^^^  ",
				" ^^^^^^^^^ ",
				"^^^^^^^^^^^",
				"     |     ",
				"     |     ",
			}
		}
	case 1: // Oak tree
		if treeHeight <= 5 {
			shape = []string{
				"  @@@  ",
				" @@@@@ ",
				"@@@@@@@",
				" @@@@@ ",
				"   |   ",
				"   |   ",
			}
		} else if treeHeight <= 6 {
			shape = []string{
				"   @@@   ",
				"  @@@@@  ",
				" @@@@@@@ ",
				"@@@@@@@@@",
				" @@@@@@@ ",
				"  @@@@@  ",
				"    |    ",
				"    |    ",
			}
		} else {
			shape = []string{
				"    @@@    ",
				"   @@@@@   ",
				"  @@@@@@@  ",
				" @@@@@@@@@ ",
				"@@@@@@@@@@@",
				" @@@@@@@@@ ",
				"  @@@@@@@  ",
				"     |     ",
				"     |     ",
			}
		}
	case 2: // Birch tree
		if treeHeight <= 5 {
			shape = []string{
				"  ###  ",
				" ##### ",
				"#######",
				"   |   ",
				"   |   ",
				"   |   ",
			}
		} else {
			shape = []string{
				"   ###   ",
				"  #####  ",
				" ####### ",
				"#########",
				"    |    ",
				"    |    ",
				"    |    ",
			}
		}
	default: // Maple tree
		if treeHeight <= 5 {
			shape = []string{
				"  &&&  ",
				" &&&&& ",
				"&&&&&&&",
				" &&&&& ",
				"  &&&  ",
				"   |   ",
				"   |   ",
			}
		} else {
			shape = []string{
				"   &&&   ",
				"  &&&&&  ",
				" &&&&&&& ",
				"&&&&&&&&&",
				" &&&&&&& ",
				"  &&&&&  ",
				"   &&&   ",
				"    |    ",
				"    |    ",
			}
		}
	}

	return Tree{
		x:      x,
		height: groundY,
		shape:  shape,
	}
}
