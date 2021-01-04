package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

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

var debug bool

var rootCmd = &cobra.Command{
	Use:   "hsk00",
	Short: "A tool to add/replace games to datafrog handheld console",
	Long:  `🚧 WIP 🚧`,
}

var makeCommand = &cobra.Command{
	Use:   "make",
	Short: "Creates hskXX.asd and hsk00.asd files",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		infiles, err := cmd.Flags().GetStringArray("in")
		if err != nil {
			return err
		}
		gamelistfile, err := cmd.Flags().GetString("hsk00")
		if err != nil {
			return err
		}
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		return encodeFile(infiles, gamelistfile, out)
	},
}
var descrambleCommand = &cobra.Command{
	Use:   "descramble",
	Short: "converts hskXX.asd files to usable zip files",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		in, err := cmd.Flags().GetString("in")
		if err != nil {
			return err
		}

		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		return decodeFileAndSave(in, out)
	},
}

var xorCommand = &cobra.Command{
	Use:   "xor",
	Short: "Just some debug helper",
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(in)
		if err != nil {
			return err
		}
		for _, c := range b {
			r := c ^ 0xe5
			os.Stdout.Write([]byte{r})
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")
	makeCommand.Flags().StringArray("in", nil, "nes files")
	makeCommand.Flags().String("out", "", "output hskXX file")
	makeCommand.Flags().String("hsk00", "hsk00.asd", "hsk00.asd  files")
	rootCmd.AddCommand(makeCommand)

	xorCommand.Flags().String("file", "", "file to xor")
	rootCmd.AddCommand(xorCommand)

	descrambleCommand.Flags().String("in", "", "scrambled file")
	descrambleCommand.Flags().String("out", "", "proper zip file")
	rootCmd.AddCommand(descrambleCommand)
}

func encodeFileName(in string) string {
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

func decodeFile(in string) ([]*zip.File, error) {
	b, err := ioutil.ReadFile(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	endReplaced := WQWToPK(b)

	// if err := ioutil.WriteFile("sample.zip", endReplaced, 0644); err != nil {
	// 	log.Printf("failed to write to zip file: %s\n", err)
	// }

	// zf, err := zip.OpenReader("sample.zip")
	// if err != nil {
	// 	log.Printf("faled to read zip file: %s\n", err)
	// }

	// for _, file := range zf.File {
	// 	log.Printf("=%s\n", decodeFileName(file.Name))
	// 	log.Printf("=%X\n", file.Name)
	// }

	r := bytes.NewReader(endReplaced)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %s", err)
	}

	return zr.File, nil
}

func decodeFileAndSave(in, out string) error {
	zipfiles, err := decodeFile(in)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	conglomerateZip, err := os.Create(out)
	if err != nil {
		log.Printf("failed to create zip writer: %s\n", err)
	}
	defer conglomerateZip.Close()
	zw := zip.NewWriter(conglomerateZip)
	defer zw.Close()

	for _, file := range zipfiles {

		header, err := zip.FileInfoHeader(file.FileInfo())
		if err != nil {
			log.Printf("failed to read reader: %s\n", err)
		}

		header.Name = encodeFileName(file.Name)
		header.Method = zip.Deflate

		fr, _ := file.Open()
		defer fr.Close()

		fw, err := zw.CreateHeader(header)
		// fw, _ := zw.Create(decodeFileName(file.Name))
		if _, err = io.Copy(fw, fr); err != nil {
			log.Printf("failed to copy to writer: %s\n", err)
		}
	}

	return nil
}

func makeGameListLine(gameID int, out, fn string) string {
	return fmt.Sprintf("%s,%s,0,2,Game%s.bin", strings.ToUpper(out[:1])+out[1:], filepath.Base(fn), fmt.Sprintf("%02d", gameID))
}

func createZip(in []string) ([]byte, error) {
	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)

	// outz := bufio.NewWriterSize(
	// 	bb,
	// 	4096*2,
	// )

	// outz, err := os.Create("hsk06.zip")
	// if err != nil {
	// 	log.Printf("failed to create zip writer: %s\n", err)
	// }
	// defer outz.Close()

	outzw := zip.NewWriter(outz)
	defer outzw.Close()

	for _, fname := range in {
		f, err := os.Open(fname)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			return nil, fmt.Errorf("%s is directory", fname)
		}
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return nil, fmt.Errorf("failed to read reader: %s", err)
		}
		header.Name = encodeFileName(fi.Name())
		header.Method = zip.Deflate

		fw, err := outzw.CreateHeader(header)
		// fw, _ := zw.Create(decodeFileName(file.Name))
		if _, err = io.Copy(fw, f); err != nil {
			return nil, fmt.Errorf("failed to copy to writer: %s", err)
		}

	}

	// if err := outzw.Flush(); err != nil {
	// 	log.Printf("failed to write to hsk06.asd: %s\n", err)
	// }

	// needs to be closed to flush to buffer
	outzw.Close()
	return bb.Bytes(), nil

}

func encodeFile(in []string, hsk00 string, out string) error {
	re := regexp.MustCompile(`^hsk(\d+).asd$`)
	matches := re.FindStringSubmatch(out)
	if len(matches) < 2 {
		return errors.New("invalid file name")
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("invalid file name: %s", err)
	}

	gameID := (id + 1) / 2
	log.Println("gameID", gameID)
	// fns := []string{"Super_Sprint.Nes",
	// 	"Ufo_Race.Nes",
	// 	"Vindicators.Nes",
	// 	"Zippy_Race.Nes"}

	// fis, err := ioutil.ReadDir("in")
	// if err != nil {
	// 	log.Printf("failed to read dir in: %s\n", err)
	// }

	// outz := bufio.NewWriterSize(
	// 	bb,
	// 	4096*2,
	// )

	// outz, err := os.Create("hsk06.zip")
	// if err != nil {
	// 	log.Printf("failed to create zip writer: %s\n", err)
	// }
	// defer outz.Close()

	outb, err := createZip(in)
	if err != nil {
		return err
	}

	if debug {
		if err := ioutil.WriteFile(out+".debug.zip", outb, 0644); err != nil {
			log.Printf("failed to write to hsk debug: %s\n", err)
		}
	}

	eReplaced := PKToWQW(outb)

	if err := ioutil.WriteFile(out, eReplaced, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %s", out, err)
	}

	gamelistfiles, err := decodeFile(hsk00)
	if err != nil {
		return nil
	}

	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)

	list := []string{}
	for _, fi := range gamelistfiles {
		log.Println("finame", encodeFileName(fi.Name))
		if encodeFileName(fi.Name) == "Hsk00.lst" {
			f, err := fi.Open()
			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(f)
			updated := false
			for scanner.Scan() {
				t := scanner.Text()
				if !strings.HasPrefix(t, strings.Title(out)[:len(out)-4]) {
					if strings.HasPrefix(t, "Hsk02.asd") {
						tt := strings.Replace(t, ",0,2", ",0,1", 1)
						list = append(list, tt)
					} else {
						list = append(list, t)
					}
				} else if !updated {
					for _, fn := range in {
						t := makeGameListLine(gameID, out, fn)
						log.Println("changing to", t)
						list = append(list, t)
					}
					updated = true
				}
			}
			// append
			if !updated {
				for _, fn := range in {
					t := makeGameListLine(gameID, out, fn)
					log.Println("appending", t)
					list = append(list, t)
				}
				updated = true
			}
			f.Close()

			if len(list) == 0 {
				return errors.New("failed to read gamelist")
			}

			w, err := zw.Create(encodeFileName("Hsk00.lst"))
			if err != nil {
				return err
			}
			w.Write([]byte(strings.Join(list, "\n")))
		} else {
			w, err := zw.Create(fi.Name)
			if err != nil {
				return err
			}

			// r, err := fi.Open()
			// if err != nil {
			// 	return err
			// }
			log.Println("games in this ", len(list))
			// w.Write([]byte(fmt.Sprintf("%X000000", len(list))))
			w.Write([]byte{byte(len(list)), 0x00, 0x00, 0x00})
			// w.Write([]byte{0x1b, 0x00, 0x00, 0x00})
			// if _, err := io.Copy(w, r); err != nil {
			// 	return err
			// }
		}
	}

	zw.Close()
	bout := bb.Bytes()
	if debug {
		ioutil.WriteFile("hsk00.asd.debug.zip", bout, 0644)
	}

	return ioutil.WriteFile("hsk00.asd", PKToWQW(bout), 0644)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

// 0,1 =
// 1,2 = works
// 1,0 = game running with sound but image
// 2,2 = works
//
