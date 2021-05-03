package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
	"time"
)

type Hero struct {
	velocity pixel.Vec
	rect pixel.Rect
	hp int
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

	//type Hero struct {
	//	velocity pixel.Vec
	//	rect pixel.Rect
	//	hp int
	//}

	hero := Hero{
		velocity: pixel.V(200, 0),
		rect: pixel.R(0,0,200,100).Moved(pixel.V(win.Bounds().W() / 4, win.Bounds().H() / 2)),
		hp: 100,
	}

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Skyblue)

		if win.Pressed(pixelgl.KeySpace) {
			hero.velocity.Y = 500
		}
		hero.rect = hero.rect.Moved(pixel.V(0, hero.velocity.Y * dt))
		hero.velocity.Y -= 900 * dt

		mat := pixel.IM
		mat = mat.ScaledXY(pixel.ZV, pixel.V(hero.rect.W() / sprite.Frame().W(), hero.rect.H() / sprite.Frame().H()))
		mat = mat.Moved(pixel.V(hero.rect.Min.X + hero.rect.W() / 2, hero.rect.Min.Y + hero.rect.H() / 2))
		sprite.Draw(win, mat)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
