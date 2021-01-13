package main

import (
	"log"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/leaanthony/mewn"
	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails"
)

type Runtime struct {
	runtime *wails.Runtime
}

// WailsInit initialize wails
func (r *Runtime) WailsInit(wr *wails.Runtime) error {
	r.runtime = wr
	return nil
}

func (r *Runtime) SelectGames() []string {
	files, err := zenity.SelectFileMutiple(zenity.Filename(""), zenity.FileFilters{{"NES ROMs", []string{"*.nes"}}})
	if err != nil {
		log.Println("err ", err)
	}
	log.Println("files", files)
	return files
}

func (r *Runtime) SelectRootDir() string {
	file, _ := zenity.SelectFile(zenity.Filename(""), zenity.Directory())
	return file
}
func (r *Runtime) AddGames(rootDir string, categoryID int, newGamesIn []interface{}) error {
	var newGames []string
	for _, g := range newGamesIn {
		newGames = append(newGames, g.(string))
	}

	return pkg.Add(rootDir, categoryID, newGames, "")
}

func main() {

	js := mewn.String("./frontend/build/static/js/main.js")
	css := mewn.String("./frontend/build/static/css/main.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  420,
		Height: 520,
		Title:  "hsk00",
		JS:     js,
		CSS:    css,
		Colour: "#0d1117",
	})

	r := &Runtime{}

	app.Bind(r)
	app.Run()
}

// import (
// 	"archive/zip"
// 	"bufio"
// 	"bytes"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"regexp"
// 	"strconv"
// 	"strings"

// 	"github.com/dev-drprasad/hsk00/util"
// )

// func createZip(in []string) ([]byte, error) {
// 	bb := bytes.NewBuffer(nil)
// 	outz := bufio.NewWriter(bb)

// 	outzw := zip.NewWriter(outz)
// 	defer outzw.Close()

// 	for _, fname := range in {
// 		f, err := os.Open(fname)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer f.Close()

// 		fi, err := f.Stat()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if fi.IsDir() {
// 			return nil, fmt.Errorf("%s is directory", fname)
// 		}
// 		header, err := zip.FileInfoHeader(fi)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read reader: %s", err)
// 		}
// 		header.Name = util.EncodeFileName(fi.Name())
// 		header.Method = zip.Deflate

// 		fw, err := outzw.CreateHeader(header)
// 		// fw, _ := zw.Create(decodeFileName(file.Name))
// 		if _, err = io.Copy(fw, f); err != nil {
// 			return nil, fmt.Errorf("failed to copy to writer: %s", err)
// 		}

// 	}

// 	// if err := outzw.Flush(); err != nil {
// 	// 	log.Printf("failed to write to hsk06.asd: %s\n", err)
// 	// }

// 	// needs to be closed to flush to buffer
// 	outzw.Close()
// 	return bb.Bytes(), nil

// }

// func encodeFile(in []string, hsk00 string, out string) error {
// 	re := regexp.MustCompile(`^hsk(\d+).asd$`)
// 	matches := re.FindStringSubmatch(out)
// 	if len(matches) < 2 {
// 		return errors.New("invalid file name")
// 	}
// 	id, err := strconv.Atoi(matches[1])
// 	if err != nil {
// 		return fmt.Errorf("invalid file name: %s", err)
// 	}

// 	gameID := (id + 1) / 2

// 	outb, err := createZip(in)
// 	if err != nil {
// 		return err
// 	}

// 	if debug {
// 		if err := ioutil.WriteFile(out+".debug.zip", outb, 0644); err != nil {
// 			log.Printf("failed to write to hsk debug: %s\n", err)
// 		}
// 	}

// 	eReplaced := util.PKToWQW(outb)

// 	if err := ioutil.WriteFile(out, eReplaced, 0644); err != nil {
// 		return fmt.Errorf("failed to write to %s: %s", out, err)
// 	}

// 	gamelistfiles, err := decompress(hsk00)
// 	if err != nil {
// 		return err
// 	}

// 	bb := bytes.NewBuffer(nil)
// 	outz := bufio.NewWriter(bb)
// 	zw := zip.NewWriter(outz)

// 	list := []string{}
// 	for _, fi := range gamelistfiles {
// 		if util.EncodeFileName(fi.Name) == "Hsk00.lst" {
// 			f, err := fi.Open()
// 			if err != nil {
// 				return err
// 			}
// 			scanner := bufio.NewScanner(f)
// 			updated := false
// 			for scanner.Scan() {
// 				t := scanner.Text()
// 				if !strings.HasPrefix(t, strings.Title(out)[:len(out)-4]) {
// 					list = append(list, t)
// 				} else if !updated {
// 					for _, fn := range in {
// 						t := makeGameListLine(gameID, out, fn)
// 						log.Println("changing to", t)
// 						list = append(list, t)
// 					}
// 					updated = true
// 				}
// 			}
// 			// append
// 			if !updated {
// 				for _, fn := range in {
// 					t := makeGameListLine(gameID, out, fn)
// 					log.Println("appending", t)
// 					list = append(list, t)
// 				}
// 				updated = true
// 			}
// 			f.Close()

// 			if len(list) == 0 {
// 				return errors.New("failed to read gamelist")
// 			}

// 			w, err := zw.Create(util.EncodeFileName("Hsk00.lst"))
// 			if err != nil {
// 				return err
// 			}
// 			w.Write([]byte(strings.Join(list, "\n")))
// 		} else {
// 			w, err := zw.Create(fi.Name)
// 			if err != nil {
// 				return err
// 			}

// 			w.Write([]byte{byte(len(list)), 0x00, 0x00, 0x00})
// 		}
// 	}

// 	zw.Close()
// 	bout := bb.Bytes()
// 	if debug {
// 		ioutil.WriteFile("hsk00.asd.debug.zip", bout, 0644)
// 	}

// 	return ioutil.WriteFile("hsk00.asd", util.PKToWQW(bout), 0644)
// }

// // 0,1 =
// // 1,2 = works
// // 1,0 = game running with sound but image
// // 2,2 = works
// //
