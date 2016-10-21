package data

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"github.com/pborman/uuid"
	"os"
	"path"
	"io/ioutil"
	"encoding/json"
)

type SimpleData struct {
	Parts    []string `json:"parts"`
	Filename string `json:"filename"`
	Tracker  string `json:"tracker"`
	UUID     uuid.UUID `json:"uuid"`
}

func New(in *os.File) *SimpleData {
	d := SimpleData{
		Parts: []string{},
		Filename: path.Base(in.Name()),
		Tracker: "simple-host:9999",
		UUID:uuid.NewRandom(),
	}
	d.SetFileParts(in)

	return &d
}

func (meta *SimpleData) SetFileParts(f io.Reader) {
	scanner := bufio.NewScanner(f)
	scanner.Split(ScanChunk)

	hasher := sha1.New()
	for scanner.Scan() {
		data := scanner.Bytes()
		hasher.Write([]byte(data))
		sha := fmt.Sprintf("%x", hasher.Sum(nil))
		meta.Parts = append(meta.Parts, sha)

	}
}

const KB int = 1024
const ChunkSize  = KB * 5

func ScanChunk(data []byte, atEOF bool) (int, []byte, error) {
	dataLength := len(data)
	if atEOF && dataLength == 0 {
		return 0, nil, nil
	}
	var till int
	if till = dataLength; till > ChunkSize {
		till = ChunkSize
	}
	return till, data[0:till], nil
}

func LoadSimpleDataFile(path string) (*SimpleData, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	simpleData := SimpleData{}
	err = json.Unmarshal(data, &simpleData)

	if err != nil {
		return nil, err
	}

	return &simpleData, nil
}