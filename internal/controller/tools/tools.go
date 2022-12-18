package tools

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func SortByAlgorithmABC(ar []string) { //"ABC" sort
	sort.Slice(ar, func(i, j int) bool {
		for k := 0; true; k++ {
			if k == len(ar[i]) {
				return false
			}
			if k == len(ar[j]) {
				return false
			}
			if ar[i][k] == ar[j][k] {
				continue
			}
			return ar[i][k] < ar[j][k]
		}
		return false
	})
}

func Swapper(first []string) (ContentTextOut string) {

	var k = true
	ContentTextOut = first[len(first)-1]
	for l := 0; k; l++ {
		if l+1 == len(first)-1 {
			k = false
		}
		ContentTextOut += " " + first[l]
	}
	ContentTextOut += "\n"
	return ContentTextOut
}

func ReaderCSV(file io.Reader) (recordAll [][]string, err error) {

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	recordAll, err = reader.ReadAll()

	return recordAll, nil
}

func OpenFile(fileName string) (file *os.File, recordAll [][]string, err error) {
	fileName = fileName
	file, err = os.Open(filepath.Join("src/simulator/" + fileName))
	if err != nil {
		fmt.Println("err")
	}
	if _, err = file.Seek(0, 0); err != nil {
		log.Printf("err seek")
	}

	recordAll, err = ReaderCSV(file)

	if err != nil {
		err = fmt.Errorf("err record")
	}

	return file, recordAll, nil
}
