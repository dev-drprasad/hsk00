package pkg

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func DecodeFileAndSave(in, out string) error {
	scrambledBytes, err := ioutil.ReadFile(in)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if bytes.HasPrefix(scrambledBytes, scrambledZipHeader) {
		if out == "" {
			out = in + ".zip"
		}
		zipfiles, err := DescrambleZipBytes(scrambledBytes)
		if err != nil {
			return fmt.Errorf("failed to descramble file '%s' <- %s", in, err)
		}

		conglomerateZip, err := os.Create(out)
		if err != nil {
			return fmt.Errorf("failed to create zip writer: %s", err)
		}
		defer conglomerateZip.Close()
		zw := zip.NewWriter(conglomerateZip)
		defer zw.Close()

		for _, file := range zipfiles {

			header, err := zip.FileInfoHeader(file.FileInfo())
			if err != nil {
				return fmt.Errorf("failed to read reader: %s", err)
			}

			header.Name = EncodeFileName(file.Name)
			header.Method = zip.Deflate

			fr, _ := file.Open()
			defer fr.Close()

			fw, err := zw.CreateHeader(header)
			// fw, _ := zw.Create(decodeFileName(file.Name))
			if _, err = io.Copy(fw, fr); err != nil {
				return fmt.Errorf("failed to copy to writer: %s", err)
			}
		}

		fmt.Printf("zip file saved as %s\n", out)
	} else if bytes.HasPrefix(scrambledBytes, []byte{0xFF, 0xDB, 0x0F, 0xD2}) {
		i := bytes.Index(scrambledBytes, []byte{0xFF, 0xD9, 0xFF})
		if i == -1 {
			return errors.New("file started with FFDB0FD2, but no image found")
		}
		if out == "" {
			out = in + ".jpg"
		}
		imageBytes := append([]byte{0xFF, 0xD8}, scrambledBytes[i+2:]...)
		if err := ioutil.WriteFile(out, imageBytes, 0644); err != nil {
			return fmt.Errorf("failed to save image file '%s' <- %s", out, err)
		}
		fmt.Printf("image saved as: %s\n", out)
		return nil
	} else {
		return errors.New("file is unknow type")
	}

	return nil
}
