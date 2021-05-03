package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

type Hero struct {
	velocity pixel.Vec
	rect pixel.Rect
	hp int
}

type Wall struct {
	X float64
	Y float64
	W float64
	H float64
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game Screen",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("flying.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	hero := Hero{
		velocity: pixel.V(200, 0),
		rect: pixel.R(0,0,170,100).Moved(pixel.V(win.Bounds().W() / 4, win.Bounds().H() / 2)),
		hp: 100,
	}

	var walls []Wall

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(win.Bounds().Center(), basicAtlas)

	distance := 0.0

	last := time.Now()

	status := "playing"

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Skyblue)

		switch status {
		case "playing":
			if len(walls) <= 0 || distance - (walls[len(walls) - 1].X - win.Bounds().W()) >= 400 {
				newWall := Wall{}
				newWall.X = distance + win.Bounds().W()
				newWall.W = win.Bounds().W() / 10
				newWall.H = rand.Float64() * win.Bounds().H() * 0.7
				if rand.Intn(2) >= 1 {
					newWall.Y = newWall.H
				} else {
					newWall.Y = win.Bounds().H()
				}
				walls = append(walls, newWall)
			}

			drawing := imdraw.New(nil)

			for _, wall := range walls {
				drawing.Color = colornames.Beige
				drawing.Push(pixel.V(wall.X - distance, wall.Y))
				drawing.Push(pixel.V(wall.X - distance + wall.W, wall.Y - wall.H))
				drawing.Rectangle(0)

				if wall.X - distance <= hero.rect.Max.X && wall.X - distance + wall.W >= hero.rect.Min.X && wall.Y >= hero.rect.Min.Y && wall.Y - wall.H <= hero.rect.Max.Y {
					drawing.Color = colornames.Red
					drawing.Push(hero.rect.Min)
					drawing.Push(hero.rect.Max)
					drawing.Rectangle(0)
					hero.hp -= 1
					if hero.hp <= 0 {
						status = "gameover"
					}
				}

				if wall.X - distance < - wall.W {
					walls = walls[1:]
				}
			}

			drawing.Draw(win)
			distance += hero.velocity.X * dt

			if hero.rect.Max.Y < 0 || hero.rect.Min.Y > win.Bounds().H() {
				status = "gameover"
			}

			if win.Pressed(pixelgl.KeySpace) {
				hero.velocity.Y = 500
			}
			hero.rect = hero.rect.Moved(pixel.V(0, hero.velocity.Y * dt))
			hero.velocity.Y -= 900 * dt

			mat := pixel.IM
			mat = mat.ScaledXY(pixel.ZV, pixel.V(hero.rect.W() / sprite.Frame().W(), hero.rect.H() / sprite.Frame().H()))
			// add
			mat = mat.Rotated(pixel.ZV, math.Atan(hero.velocity.Y / (hero.velocity.X * 3)))
			mat = mat.Moved(pixel.V(hero.rect.Min.X + hero.rect.W() / 2, hero.rect.Min.Y + hero.rect.H() / 2))
			sprite.Draw(win, mat)
		case "gameover":
			//os.Exit(3)
			basicTxt.Clear()
			basicTxt.Color = colornames.Green
			line := fmt.Sprintf("Game Over! Score: %d\n", int(distance))
			basicTxt.Dot.X -= basicTxt.BoundsOf(line).W() / 2
			fmt.Fprintf(basicTxt, line)
			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
			if win.Pressed(pixelgl.KeySpace) {
				hero.hp = 100
				hero.rect = pixel.R(0,0,100,100).Moved(pixel.V(win.Bounds().W() / 4, win.Bounds().H() / 2))
				status = "playing"
				distance = 0.0
				last = time.Now()
				walls = walls[:0]
			}
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
