package main

import "github.com/veandco/go-sdl2/sdl"

type food struct {
	pos sdl.Rect
}

func newFood(snek *snek) *food {
	for {
		x := randInt(20, winWidth-snekSize)
		y := randInt(20, winHeight-snekSize)
		x = x - x%snekSize
		y = y - y%snekSize
		pos := sdl.Rect{x, y, snekSize, snekSize}
		if snek.collidesWithNewFood(&pos) {
			continue
		}

		return &food{
			pos: pos,
		}
	}
	return nil
}
