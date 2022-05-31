package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

var snekSize int32 = 20

type snek struct {
	x, y     int32
	body     []sdl.Rect
	movement [2]int32
	tailDir  [2]int32
}

func newSnek() *snek {
	s := snek{
		x:        winWidth / 2,
		y:        winHeight / 2,
		body:     make([]sdl.Rect, 1),
		movement: [2]int32{-1, 0},
		tailDir:  [2]int32{-1, 0},
	}
	s.body[0] = sdl.Rect{s.x, s.y, snekSize, snekSize}
	return &s
}

func (s *snek) len() int32 {
	return int32(len(s.body))
}

func (s *snek) updateSnekBody() {
	prevX, prevY := s.x, s.y
	if s.isGoingLeft() {
		s.x = s.x - snekSize
	}
	if s.isGoingRight() {
		s.x = s.x + snekSize
	}
	if s.isGoingUp() {
		s.y = s.y - snekSize
	}
	if s.isGoingDown() {
		s.y = s.y + snekSize
	}
	s.body[0] = sdl.Rect{s.x, s.y, snekSize, snekSize}
	for i := 1; i < len(s.body); i++ {
		nextX, nextY := s.body[i].X, s.body[i].Y
		if i == len(s.body)-1 {
			// we know we reached the tail
			if nextX > prevX {
				// is moving left
				s.tailDir[0] = -1
				s.tailDir[1] = 0
			} else if nextX < prevX {
				// is moving right
				s.tailDir[0] = 1
				s.tailDir[1] = 0
			} else if nextY > prevY {
				// is moving up
				s.tailDir[0] = 0
				s.tailDir[1] = -1
			} else if nextY < prevY {
				// is moving down
				s.tailDir[0] = 0
				s.tailDir[1] = 1
			}
		}
		s.body[i] = sdl.Rect{prevX, prevY, snekSize, snekSize}
		prevX, prevY = nextX, nextY
	}
}

func (s *snek) goLeft() {
	if s.isGoingRight() {
		return
	}
	s.movement[0] = -1
	s.movement[1] = 0
}

func (s *snek) goRight() {
	if s.isGoingLeft() {
		return
	}
	s.movement[0] = 1
	s.movement[1] = 0
}

func (s *snek) goDown() {
	if s.isGoingUp() {
		return
	}
	s.movement[0] = 0
	s.movement[1] = 1
}

func (s *snek) goUp() {
	if s.isGoingDown() {
		return
	}
	s.movement[0] = 0
	s.movement[1] = -1
}

func (s *snek) isGoingLeft() bool {
	return s.movement[0] == -1 && s.movement[1] == 0
}

func (s *snek) isGoingRight() bool {
	return s.movement[0] == 1 && s.movement[1] == 0
}

func (s *snek) isGoingDown() bool {
	return s.movement[0] == 0 && s.movement[1] == 1
}

func (s *snek) isGoingUp() bool {
	return s.movement[0] == 0 && s.movement[1] == -1
}

func (s *snek) hasCollided() bool {
	return s.isOffBoundaries() || s.collidesWithItself()
}

func (s *snek) isOffBoundaries() bool {
	return s.x >= winWidth || s.x <= 0 || s.y <= 0 || s.y >= winHeight
}

func (s *snek) collidesWithItself() bool {
	if len(s.body) == 1 {
		return false
	}
	for i := 1; i < len(s.body); i++ {
		if s.body[0].HasIntersection(&s.body[i]) {
			return true
		}
	}
	return false
}

func (s *snek) ateFood(f *food) bool {
	return s.body[0].HasIntersection(&f.pos)
}

func (s *snek) grow() {
	x, y := s.body[len(s.body)-1].X, s.body[len(s.body)-1].Y
	if s.tailIsGoingLeft() {
		s.body = append(s.body, sdl.Rect{x + snekSize, y, snekSize, snekSize})
		return
	}
	if s.tailIsGoingRight() {
		s.body = append(s.body, sdl.Rect{x - snekSize, y, snekSize, snekSize})
		return
	}
	if s.tailIsGoingUp() {
		s.body = append(s.body, sdl.Rect{x, y + snekSize, snekSize, snekSize})
		return
	}
	if s.tailIsGoingLeft() {
		s.body = append(s.body, sdl.Rect{x, y - snekSize, snekSize, snekSize})
		return
	}
}

// check tail direction
func (s *snek) tailIsGoingLeft() bool {
	return s.tailDir[0] == -1 && s.tailDir[1] == 0
}

func (s *snek) tailIsGoingRight() bool {
	return s.tailDir[0] == 1 && s.tailDir[1] == 0
}

func (s *snek) tailIsGoingDown() bool {
	return s.tailDir[0] == 0 && s.tailDir[1] == 1
}

func (s *snek) tailIsGoingUp() bool {
	return s.tailDir[0] == 0 && s.tailDir[1] == -1
}

// check if new food collides with snake's body
func (s *snek) collidesWithNewFood(food *sdl.Rect) bool {
	for i := 0; i < len(s.body); i++ {
		if food.HasIntersection(&s.body[i]) {
			return true
		}
	}
	return false
}
