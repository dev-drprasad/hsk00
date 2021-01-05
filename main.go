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
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/markbates/pkger"
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

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "add games to category",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println(args)
		// games, err := cmd.Flags().GetStringArray("nes")
		// if err != nil {
		// 	return err
		// }

		if len(args) == 0 {
			return errors.New("pass games path as argument")
		}

		categoryID, err := cmd.Flags().GetInt("category")
		if err != nil {
			return err
		}

		if categoryID < 0 {
			return errors.New("category number must be 0 or greater")
		}

		rootDir, err := cmd.Flags().GetString("root")
		if err != nil {
			return err
		}
		if _, err := os.Stat(rootDir); os.IsNotExist(err) {
			return fmt.Errorf("root directory %s not exists", rootDir)
		}

		return add(rootDir, categoryID, args)
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

	addCommand.Flags().Int("category", 0, "number of category starting from 0, left -> right")
	addCommand.MarkFlagRequired("category")
	// addCommand.Flags().StringArray("nes", nil, "location of nes game file")
	// addCommand.MarkFlagRequired("nes")
	addCommand.Flags().String("root", "", "root path of sd card")
	addCommand.MarkFlagRequired("root")
	rootCmd.AddCommand(addCommand)
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

func decompress(in string) ([]*zip.File, error) {
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

func decodeFileAndSave(in, out string) error {
	zipfiles, err := decompress(in)
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

	gamelistfiles, err := decompress(hsk00)
	if err != nil {
		return err
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
					list = append(list, t)
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

			w.Write([]byte{byte(len(list)), 0x00, 0x00, 0x00})
		}
	}

	zw.Close()
	bout := bb.Bytes()
	if debug {
		ioutil.WriteFile("hsk00.asd.debug.zip", bout, 0644)
	}

	return ioutil.WriteFile("hsk00.asd", PKToWQW(bout), 0644)
}

func getMenuList(hsk00Path string) ([]string, error) {
	hsk00files, err := decompress(hsk00Path)
	if err != nil {
		return nil, err
	}

	menuList := []string{}

	for _, fi := range hsk00files {
		if encodeFileName(fi.Name) == "Hsk00.lst" {
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

func fileNameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func makeHsk00(menuList []string) ([]byte, error) {
	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)

	w1, err := zw.Create(encodeFileName("Hsk00.lst"))
	if err != nil {
		return nil, err
	}
	w1.Write([]byte(strings.Join(menuList, "\n")))

	w2, err := zw.Create(encodeFileName("GameNumber.bin"))
	if err != nil {
		return nil, err
	}

	w2.Write([]byte{byte(len(menuList)), 0x00, 0x00, 0x00})
	zw.Close()
	return bb.Bytes(), nil
}

func compressReaderAndWrite(zw *zip.Writer, r io.ReadCloser, fname string) error {
	w, err := zw.Create(encodeFileName(fname))
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy to writer: %s", err)
	}
	return nil
}

func compressFilesAndWrite(zw *zip.Writer, filePaths []string) error {
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

		if err := compressReaderAndWrite(zw, f, fi.Name()); err != nil {
			return err
		}

		// header, err := zip.FileInfoHeader(fi)
		// if err != nil {
		// 	return fmt.Errorf("failed to read header: %s", err)
		// }
		// header.Name = encodeFileName(fi.Name())

		// fw, err := zw.CreateHeader(header)
		// if _, err = io.Copy(fw, f); err != nil {
		// 	return fmt.Errorf("failed to copy to writer: %s", err)
		// }
	}

	return nil
}

func update(filePath string, gamePaths []string) error {
	zipFileInfos, err := decompress(filePath)
	if err != nil {
		return fmt.Errorf("failed to decompress file %s: %s", filePath, err)
	}

	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)
	defer zw.Close()

	for _, zipfi := range zipFileInfos {
		header, err := zip.FileInfoHeader(zipfi.FileInfo())
		if err != nil {
			return fmt.Errorf("failed to read header: %s", err)
		}

		fr, _ := zipfi.Open()
		defer fr.Close()

		fw, err := zw.CreateHeader(header)
		if _, err = io.Copy(fw, fr); err != nil {
			return fmt.Errorf("failed to copy to writer: %s", err)
		}
	}

	if err := compressFilesAndWrite(zw, gamePaths); err != nil {
		return fmt.Errorf("failed to write to zip writer: %s", err)
	}

	if err := zw.Close(); err != nil {
		return err
	}

	zipBytes := bb.Bytes()
	scambledBytes := PKToWQW(zipBytes)
	return ioutil.WriteFile(filePath, scambledBytes, 0644)
}

func create(outFilePath string, gamePaths []string) error {
	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)
	defer zw.Close()

	if err := compressFilesAndWrite(zw, gamePaths); err != nil {
		return fmt.Errorf("failed to write to zip writer: %s", err)
	}

	if err := zw.Close(); err != nil {
		return err
	}

	zipBytes := bb.Bytes()
	scambledBytes := PKToWQW(zipBytes)
	return ioutil.WriteFile(outFilePath, scambledBytes, 0644)
}

func add(rootDir string, categoryID int, newGames []string) error {
	gamesDirectoryName := fmt.Sprintf("Game%02d", categoryID)
	gamesPath := path.Join(rootDir, gamesDirectoryName)
	if _, err := os.Stat(gamesPath); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' doesn't exist", gamesPath)
	}
	hsk00Path := path.Join(gamesPath, "hsk00.asd")
	menuList, err := getMenuList(hsk00Path)
	if err != nil {
		return fmt.Errorf("failed to get list from hsk00.asd: %s", err)
	}
	if debug {
		log.Println("current game list:")
		for i, name := range menuList {
			log.Printf("%03d. %s\n", i+1, name)
		}
	}

	initialEndIndex := 5 - (len(menuList) % 5)
	for endIndex := initialEndIndex; endIndex <= len(newGames)+5; endIndex += 5 {
		y := endIndex
		// dont mutate endIndex, will cause ∞ loop
		if endIndex > len(newGames) {
			y = len(newGames)
		}
		startIndex := endIndex - 5
		if startIndex < 0 {
			startIndex = 0
		}
		batch := newGames[startIndex:y]

		hskID := (len(menuList) / 5) + 1
		hskFileName := fmt.Sprintf("Hsk%02d.asd", hskID)
		hskFilePath := path.Join(gamesPath, hskFileName)
		// update only last partial hsk files. create all other files from scratch
		if _, err := os.Stat(hskFilePath); err == nil && endIndex < 5 {
			if err := update(hskFilePath, batch); err != nil {
				return fmt.Errorf("failed to update games to file: %s: %s", hskFilePath, err)
			}
		} else {
			if err := create(hskFilePath, batch); err != nil {
				return fmt.Errorf("failed to create file: %s: %s", hskFilePath, err)
			}
		}

		newMenuImageFileNames := make(map[string]int) // simple set
		for _, gamePath := range batch {
			pageNo := len(menuList)/10 + 1
			menuImageFileName := fmt.Sprintf("Game%02d.bin", pageNo)
			menuListItem := fmt.Sprintf("%s,%s,0,2,%s", hskFileName, filepath.Base(gamePath), menuImageFileName)
			menuList = append(menuList, menuListItem)
			newMenuImageFileNames[menuImageFileName] = pageNo
		}

		binPrefixF, err := pkger.Open("/assets/binprefix")
		if err != nil {
			return fmt.Errorf("failed to open binprefix file: %s", err)
		}
		binPrefix, err := ioutil.ReadAll(binPrefixF)
		if err != nil {
			return fmt.Errorf("failed to read binprefix: %s", err)
		}
		for fname, pageNo := range newMenuImageFileNames {
			start := (pageNo - 1) * 10
			end := pageNo * 10
			if end > len(menuList) {
				end = len(menuList)
			}
			menuListInPage := menuList[start:end]
			gameNames := []string{}
			for i, menuItem := range menuListInPage {
				gameFilename := strings.SplitN(menuItem, ",", 3)[1]
				gameFileName := fileNameWithoutExtension(gameFilename)
				gameName := strings.ReplaceAll(gameFileName, "_", " ")
				gameNames = append(gameNames, fmt.Sprintf("%02d.%s", start+i+1, gameName))
			}

			imageBytes, err := generateMenuImage(gameNames)
			if err != nil {
				return fmt.Errorf("menu image generation failed: %s", err)
			}

			filePath := path.Join(gamesPath, fname)
			if debug {
				ioutil.WriteFile(filePath+".debug.jpg", imageBytes, 0644)
			}

			scrambledBytes := append(binPrefix, imageBytes[2:]...)
			if err := ioutil.WriteFile(filePath, scrambledBytes, 0644); err != nil {
				return nil
			}
		}
		hskID++
	}

	if debug {
		log.Println("new game list:")
		log.Printf("%s\n", strings.Join(menuList, "\n"))
	}

	listZipBytes, err := makeHsk00(menuList)
	if err != nil {
		return err
	}

	hsk00FilePath := path.Join(gamesPath, "hsk00.asd")
	if err := ioutil.WriteFile(hsk00FilePath, PKToWQW(listZipBytes), 0644); err != nil {
		return err
	}

	return nil
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
