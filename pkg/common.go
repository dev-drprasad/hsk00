package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

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
