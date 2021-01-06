package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"io/ioutil"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/markbates/pkger"
)

const canvasWidth = 720
const canvasHeight = 480

const listStartX = 200
const listStartY = 114
const listGap = 8
const fontSize = 24
const imageQuality = 100

var fonts = map[string]string{
	"Gotham-Medium": "Gotham-SemiBold.ttf",
	"Video-Phreak":  "Video-Phreak.ttf",
}

func loadFontFacePkger(dc *gg.Context, fpath string, points float64) error {
	f, err := pkger.Open(fpath)
	if err != nil {
		return fmt.Errorf("failed to open font file %s: %s", fpath, err)
	}
	fontBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	p, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}

	face := truetype.NewFace(p, &truetype.Options{
		Size: points,
	})

	dc.SetFontFace(face)
	return nil
}

func generateMenuImage(gameNames []string, fontName string) ([]byte, error) {
	fontFilename := fonts["Gotham-Medium"]
	if fn, ok := fonts[fontName]; ok {
		fontFilename = fn
	}

	bgf, err := pkger.Open("/assets/bg.jpg")
	im, _, err := image.Decode(bgf)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContext(canvasWidth, canvasHeight)
	dc.Clear()
	// if err := dc.LoadFontFace("./OpenSans-Bold.ttf", fontSize); err != nil {
	// 	panic(err)
	// }
	if err := loadFontFacePkger(dc, fmt.Sprintf("/assets/%s", fontFilename), fontSize); err != nil {
		return nil, err
	}

	dc.DrawImage(im, 0, 0)
	n := 3
	for i, name := range gameNames {
		listItemStartY := float64(listStartY + (i * (listGap + fontSize)))
		// dc.SetRGB(float64(255)/255, float64(174)/255, float64(182)/255)
		dc.SetRGB(float64(210)/255, float64(220)/255, float64(225)/255)
		for dy := -n; dy <= n; dy++ {
			for dx := -n; dx <= n; dx++ {
				if dx*dx+dy*dy >= n*n {
					// give it rounded corners
					continue
				}

				x := listStartX + float64(dx)
				y := listItemStartY + float64(dy)
				dc.DrawStringAnchored(name, x, y, 0, 1)
			}
		}
		// dc.SetRGB(float64(57)/255, float64(49)/255, float64(75)/255)
		dc.SetRGB(float64(40)/255, float64(40)/255, float64(40)/255)
		dc.DrawStringAnchored(name, listStartX, listItemStartY, 0, 1)
	}

	var opt jpeg.Options
	opt.Quality = imageQuality

	bb := bytes.NewBuffer(nil)
	w := bufio.NewWriter(bb)

	if err := jpeg.Encode(w, dc.Image(), &opt); err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

func getMenuList(hsk00Path string) ([]string, error) {
	hsk00files, err := Descramble(hsk00Path)
	if err != nil {
		return nil, err
	}

	menuList := []string{}

	for _, fi := range hsk00files {
		if EncodeFileName(fi.Name) == "Hsk00.lst" {
			f, err := fi.Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				menuList = append(menuList, scanner.Text())
			}
			break
		}
	}
	return menuList, nil
}

func makeHsk00(menuList []string) ([]byte, error) {
	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)

	w1, err := zw.Create(EncodeFileName("Hsk00.lst"))
	if err != nil {
		return nil, err
	}
	w1.Write([]byte(strings.Join(menuList, "\n")))

	w2, err := zw.Create(EncodeFileName("GameNumber.bin"))
	if err != nil {
		return nil, err
	}

	w2.Write([]byte{byte(len(menuList)), 0x00, 0x00, 0x00})
	zw.Close()
	return bb.Bytes(), nil
}
