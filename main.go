package main

import (
	"embed"
	"encoding/hex"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

//go:embed pokemon/main-sprites/crystal/*.png
var sprites embed.FS
var parser = regexp.MustCompile(`Party\:\s*(.+)`)
var window fyne.Window

func loadSprite(id int) (*canvas.Image, error) {
	sprite := fmt.Sprintf("pokemon/main-sprites/crystal/%d.png", id)
	if id == 201 {
		sprite = "pokemon/main-sprites/crystal/201-a.png"
	}
	data, err := sprites.Open(sprite)
	if err != nil {
		return nil, err
	}
	return canvas.NewImageFromReader(data, sprite), nil
}

func refresh(party []int) error {
	grid := container.NewGridWithColumns(6)
	for i := range party {
		sprite, err := loadSprite(party[i])
		if err != nil {
			return err
		}
		sprite.FillMode = canvas.ImageFillOriginal
		sprite.ScaleMode = canvas.ImageScalePixels
		grid.Add(container.NewStack(sprite))
	}
	window.SetContent(container.NewStack(canvas.NewRectangle(color.White), grid))
	return nil
}

func main() {
	a := app.New()
	window = a.NewWindow("Crystal Kaizo+ Party")
	err := refresh([]int{166, 166, 166, 166, 166, 166})
	if err != nil {
		panic(err)
	}
	window.Resize(fyne.NewSize(400, 200))
	window.SetFixedSize(false)
	quit := make(chan struct{})
	go watch(quit)
	window.ShowAndRun()
	close(quit)
}

func watch(quit chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := update()
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
		case <-quit:
			return
		}
	}
}

func update() error {
	resp, err := http.Get("http://localhost:31123/update")
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	hexData := parser.FindSubmatch(body)[1]
	data := make([]byte, len(hexData)/2)
	_, err = hex.Decode(data, hexData)
	if err != nil {
		return err
	}
	species := readPokemonList(data, 0, 6)
	return refresh(species)
}

func readPokemonList(bytes []byte, start int, capacity int) []int {
	count := int(bytes[start])
	p := start + 1
	var species []int
	for i := 0; i < count; i++ {
		species = append(species, int(bytes[p+i]))
	}
	p += capacity + 1
	return species
}
