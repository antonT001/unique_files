package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/inancgumus/screen"
)

type InfoFile struct {
	buf             string
	repeatedRequest bool
	byteMap         map[string]string
	hashMap         map[string]string
}

var unique map[int64]InfoFile

func hashFilePath(filePath string) string {

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func scanFiles(filesFrom string) int {
	i := 0
	err := filepath.Walk(filesFrom,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if !info.IsDir() {
				i++
				screen.MoveTopLeft()
				fmt.Println("scaning", i, "file(s)")
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return i
}

func readBits(path string, info os.FileInfo) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	buf := make([]byte, 50)
	_, err = f.ReadAt(buf, (info.Size() - 50))
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func biteComparison(path string, info os.FileInfo) (string, bool) {
	str := unique[info.Size()]
	elementCurrentFile := readBits(path, info)

	if !str.repeatedRequest {
		elementSavedFile := readBits(str.buf, info)
		if elementCurrentFile == elementSavedFile {
			str.repeatedRequest = true
			byteMap := map[string]string{}
			byteMap[elementCurrentFile] = str.buf
			str.byteMap = byteMap
			return str.buf, true
		}
		byteMap := map[string]string{}
		byteMap[elementCurrentFile] = path
		byteMap[elementSavedFile] = str.buf
		str.byteMap = byteMap
		return path, false

	}
	byteMap := str.byteMap
	if _, ok := byteMap[elementCurrentFile]; ok {
		return path, true
	}
	byteMap[elementCurrentFile] = path
	return path, false
	
}

func hashComparison(path string, sp string, info os.FileInfo) bool {
	str := unique[info.Size()]
	hash := hashFilePath(path)
	if path != sp {
		hashSp := hashFilePath(sp)
		if hash == hashSp {
			hashMap := make(map[string]string)
			hashMap[hash] = path
			hashMap[hashSp] = sp
			str.hashMap = hashMap
			return true
		}
		hashMap := make(map[string]string)
		hashMap[hashSp] = sp
		str.hashMap = hashMap
		return false

	}
	hashMap := str.hashMap
	if _, ok := hashMap[hash]; ok {
		return true
	}
	hashMap[hash] = path
	return false

}

func main() {

	if len(os.Args) != 2 {

		log.Fatal("Invalid path in arguments. Аrgument example: /Users/username/Documents/files")
	}

	screen.Clear()

	filesFrom := os.Args[1]

	itemsFile := scanFiles(filesFrom)

	unique = make(map[int64]InfoFile, itemsFile/4)

	dt := time.Now()

	newNameFolder := dt.Format("2006-01-02 15:04:05")

	err := os.Mkdir("doubles "+newNameFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}
	var item = 0
	var itemCopy = 0
	err = filepath.Walk(filesFrom,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if !info.IsDir() {
				item++
				screen.MoveTopLeft()
				fmt.Println("verified", item, "from", itemsFile, "file(s)")

				if _, ok := unique[info.Size()]; ok {

					if sp, ok := biteComparison(path, info); ok {
						if ok := hashComparison(path, sp, info); ok {
							itemCopy++
							err := os.Rename(path, "./doubles "+newNameFolder+"/"+info.Name())
							if err != nil {
								log.Fatal(err)
							}
						}
					}

				} else {
					str := InfoFile{}
					str.buf = path
					unique[info.Size()] = str
				}
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	fmt.Println("moved", itemCopy, "file(s)")
}
