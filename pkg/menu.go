package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"io/ioutil"

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

const defaultFont = "NotoSans"
const defaultBg = "default"

var fonts = map[string]string{
	defaultFont:      "NotoSans-SemiBold.ttf",
	"04B_30":         "04B_30.ttf",
	"VideoPhreak":    "Video-Phreak.ttf",
	"PressStart2P":   "PressStart2P.ttf",
	"SuperMarioBros": "Super Mario Bros 2.ttf",
}

var backgrounds = map[string]string{
	defaultBg:        "default.jpg",
	"SuperMarioBros": "SuperMarioBros.jpg",
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

func generateMenuImage(gameNames []string, fontName string, bgName string) ([]byte, error) {
	fontFilename := fonts[defaultFont]
	if fn, ok := fonts[fontName]; ok {
		fontFilename = fn
	}

	bgFilename := backgrounds[defaultBg]
	if fn, ok := backgrounds[bgName]; ok {
		bgFilename = fn
	}
	bgf, err := pkger.Open(fmt.Sprintf("/assets/%s", bgFilename))
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
	n := 2.5
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
	b, err := getHsk00lstContent(hsk00Path)
	if err != nil {
		return nil, err
	}

	menuList := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		menuList = append(menuList, scanner.Text())
	}

	return menuList, nil
}

func makeHsk00(menuList GameItemList) ([]byte, error) {
	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)
	defer zw.Close()

	w1, err := zw.Create(EncodeFileName("Hsk00.lst"))
	if err != nil {
		return nil, err
	}

	for _, g := range menuList {
		if g.Hsk == "" || g.Filename == "" || g.BGFilename == "" {
			return nil, fmt.Errorf("one of them is empty -> hsk: '%s', fn: '%s', bg fn: '%s'", g.Hsk, g.Filename, g.BGFilename)
		}
		w1.Write([]byte(fmt.Sprintf("%s,%s,0,1,%s\n", g.Hsk, g.Filename, g.BGFilename)))
	}

	w2, err := zw.Create(EncodeFileName("GameNumber.bin"))
	if err != nil {
		return nil, err
	}

	h := fmt.Sprintf("%08x", len(menuList))
	var hexFlipped string
	for i := len(h); i >= 2; i -= 2 {
		hexFlipped += h[i-2 : i]
	}
	b, _ := hex.DecodeString(hexFlipped)
	w2.Write(b)
	zw.Close()
	return bb.Bytes(), nil
}
