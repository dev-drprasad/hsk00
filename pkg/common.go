package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dev-drprasad/hsk00/util"
)

var hskidrx = regexp.MustCompile(`(?i)^hsk(\d+).asd`)

var Debug = false

/*
 * 57 51 57 03 -> 50 4B 03 04 file header
 * 57 51 57 02 -> 50 4B 01 02 central directory file header
 * 57 51 57 01 -> 50 4B 05 06 End of central directory record
 */

var fileHeaderEnc = []byte{0x57, 0x51, 0x57, 0x03}
var fileHeaderDec = []byte{0x50, 0x4b, 0x03, 0x04}
var centralFileHeaderEnc = []byte{0x57, 0x51, 0x57, 0x02}
var centralFileHeaderDec = []byte{0x50, 0x4b, 0x01, 0x02}
var endEnc = []byte{0x57, 0x51, 0x57, 0x01}
var endDec = []byte{0x50, 0x4b, 0x05, 0x06}

func EncodeFileName(in string) string {
	bs := []byte(in)
	var out []byte
	for _, c := range bs {
		out = append(out, byte(0xE5^c))
	}
	return string(out)
}

func WQWToPK(in []byte) []byte {
	headerReplaced := bytes.ReplaceAll(in, fileHeaderEnc, fileHeaderDec)
	centralHeaderReplaced := bytes.ReplaceAll(headerReplaced, centralFileHeaderEnc, centralFileHeaderDec)
	endReplaced := bytes.ReplaceAll(centralHeaderReplaced, endEnc, endDec)
	return endReplaced
}

func PKToWQW(in []byte) []byte {
	hreplaced := bytes.ReplaceAll(in, fileHeaderDec, fileHeaderEnc)
	crReplaced := bytes.ReplaceAll(hreplaced, centralFileHeaderDec, centralFileHeaderEnc)
	eReplaced := bytes.ReplaceAll(crReplaced, endDec, endEnc)
	return eReplaced
}

func CompressReaderAndWrite(zw *zip.Writer, r io.ReadCloser, fname string) error {
	w, err := zw.Create(EncodeFileName(fname))
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy to writer: %s", err)
	}
	return nil
}

func CompressFilesAndWrite(zw *zip.Writer, filePaths []string) error {
	for _, fname := range filePaths {
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("%s is directory", fname)
		}

		if err := CompressReaderAndWrite(zw, f, fi.Name()); err != nil {
			return err
		}
	}

	return nil
}

func Descramble(in string) ([]*zip.File, error) {
	scrambledBytes, err := ioutil.ReadFile(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}
	zipBytes := WQWToPK(scrambledBytes)

	r := bytes.NewReader(zipBytes)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %s", err)
	}

	return zr.File, nil
}

type GameItem struct {
	ID         int    `json:"id"`       // game number
	Hsk        string `json:"hsk"`      // hsk filename where game present
	Filename   string `json:"filename"` // filename in hsk00.lst
	SourcePath string `json:"srcPath"`  // source file path. unsaved file path
	BGFilename string `json:"-"`        // ex: Game03.bin
	Name       string `json:"name"`     // name of game in menu list
}

type GameItemList []*GameItem

func (l GameItemList) Len() int {
	return len(l)
}
func (l GameItemList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}
func (l GameItemList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l GameItemList) NextHskID() (int, error) {
	max := 0
	for _, g := range l {
		fmt.Println("g", g)
		matches := hskidrx.FindStringSubmatch(g.Hsk)
		if len(matches) != 2 {
			return 0, fmt.Errorf("for hsk '%s', expected match len is 2, but got %d", g.Hsk, len(matches))
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse hsk '%s': %s", matches[1], err)
		}
		if max < id {
			max = id
		}
	}
	return max + 1, nil
}

func ParseHsk00lstContent(content []byte) (GameItemList, error) {
	if content == nil {
		return nil, errors.New("content is empty")
	}

	var list GameItemList
	scanner := bufio.NewScanner(bytes.NewReader(content))
	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, ",")

		if len(lineSlice) != 5 {
			return nil, fmt.Errorf("expected items per line are 5, but got %d", len(line))
		}

		i := GameItem{
			ID:         index + 1,
			Hsk:        lineSlice[0],
			Filename:   lineSlice[1],
			Name:       GameNameFromFilename(lineSlice[1]),
			BGFilename: lineSlice[4],
		}
		list = append(list, &i)
		index++
	}

	return list, nil
}

func GameNameFromFilename(fn string) string {
	return strings.ReplaceAll(util.FileNameWithoutExt(fn), "_", " ")
}

func getHsk00lstContent(hsk00Path string) ([]byte, error) {
	hsk00files, err := Descramble(hsk00Path)
	if err != nil {
		return nil, err
	}

	var contentB []byte
	for _, fi := range hsk00files {
		if EncodeFileName(fi.Name) == "Hsk00.lst" {
			f, err := fi.Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			contentB = b
			break
		}
	}
	if contentB == nil {
		return nil, errors.New("no hsk00.lst found")
	}
	return contentB, nil
}
