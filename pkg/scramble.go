package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
)

func GenerateScrambledZip(infiles []zipFilePath, outFilePath string) error {
	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)
	defer zw.Close()

	if err := CompressFilesAndWrite(zw, infiles); err != nil {
		return fmt.Errorf("failed to write to zip writer: %s", err)
	}

	if err := zw.Close(); err != nil {
		return err
	}

	zipBytes := bb.Bytes()
	scambledBytes := PKToWQW(zipBytes)
	if err := ioutil.WriteFile(outFilePath, scambledBytes, 0644); err != nil {
		return err
	}
	fmt.Printf("file saved as %s\n", outFilePath)
	return nil
}
