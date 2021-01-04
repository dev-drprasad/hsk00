package main

import (
	"log"

	"github.com/fogleman/gg"
)

const canvasWidth = 720
const canvasHeight = 480

const listStartX = 200
const listStartY = 116
const listGap = 8
const fontSize = 24
const imageQuality = 100

var gameList = []string{"Super Mario", "Contra", "Chrono Trigger", "Super Internation Cricket", "Bomber Man 2",
	"Dr. Mario", "Adventure Island 2", "F1 Race", "Route 16 Turboy", "James_Bond_Jr.nes"}

func testImageGeneration() {
	im, err := gg.LoadImage("bg.jpg")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(canvasWidth, canvasHeight)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("/Library/Fonts/Arial Unicode.ttf", fontSize); err != nil {
		panic(err)
	}

	dc.DrawImage(im, 0, 0)
	for i, name := range gameList {
		dc.DrawStringAnchored(name, listStartX, float64(listStartY+(i*(listGap+fontSize))), 0, 1)
	}

	gg.SaveJPG("out.jpg", dc.Image(), imageQuality)
}
