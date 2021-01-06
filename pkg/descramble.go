package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func DecodeFileAndSave(in, out string) error {
	zipfiles, err := Descramble(in)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
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

	return nil
}
