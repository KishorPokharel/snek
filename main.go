package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var winWidth, winHeight int32 = 700, 700
var running = true

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var score int32
	window, err := sdl.CreateWindow(
		"Snek",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("could not create renderer: %v", err)
	}
	defer renderer.Destroy()

	if err = ttf.Init(); err != nil {
		return fmt.Errorf("could not init ttf: %v", err)
	}
	defer ttf.Quit()

	font, err := ttf.OpenFont("res/fonts/SplineSansMonoBold.ttf", 50)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}
	defer font.Close()

	snek := newSnek()
	f := newFood(snek)
	if f == nil {
		return fmt.Errorf("could not create food")
	}

	for running {
		pollEvent(snek)

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// draw food
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.FillRect(&f.pos)

		snek.updateSnekBody()

		//check boundary
		if snek.hasCollided() {
			running = false
		}
		if snek.ateFood(f) {
			score++
			snek.grow()
			f = newFood(snek)
			if f == nil {
				return fmt.Errorf("could not create food")
			}
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.FillRect(&f.pos)
		}

		// draw snake
		renderer.SetDrawColor(0, 255, 0, 255)
		renderer.FillRects(snek.body)

		renderer.Present()
		sdl.Delay(100)
	}
	text := fmt.Sprintf("Score: %d", score)
	fontSurface, err := font.RenderUTF8Solid(text, sdl.Color{255, 255, 255, 255})
	if err != nil {
		return fmt.Errorf("could not create font surface: %v", err)
	}
	defer fontSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		return fmt.Errorf("could not create texture from font surface: %v", err)
	}
	defer texture.Destroy()

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	var textSize int32 = 400
	r := &sdl.Rect{
		(winWidth / 2) - (textSize / 2),
		(winHeight / 2) - (textSize / 2),
		textSize,
		textSize,
	}

	renderer.Copy(texture, nil, r)
	renderer.Present()
	time.Sleep(2 * time.Second)

	return nil
}

func randInt(min, max int32) int32 {
	return min + int32(rand.Intn(int(max-min)))
}

func pollEvent(snek *snek) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_UP, sdl.K_k:
					snek.goUp()
				case sdl.K_DOWN, sdl.K_j:
					snek.goDown()
				case sdl.K_RIGHT, sdl.K_l:
					snek.goRight()
				case sdl.K_LEFT, sdl.K_h:
					snek.goLeft()
				}
			}
		}
	}
}
