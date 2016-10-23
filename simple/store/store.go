package store

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/vbabiy/simple/simple/data"
	"github.com/vbabiy/simple/simple/sfile"
)

const StorePath = ".store"

type store struct {
	records map[string]*data.SimpleData
}

func (s *store) All() []*data.SimpleData {
	out := make([]*data.SimpleData, len(s.records))
	var idx int
	for _, v := range s.records {
		out[idx] = v
		idx++
	}
	return out
}

func (s *store) Reload() error {
	s.records = make(map[string]*data.SimpleData) // reset the map[type]type
	if _, err := os.Stat(StorePath); os.IsNotExist(err) {
		// StorePath does not exist. We try to create it.
		err = os.Mkdir(StorePath, 0700)
		if err != nil {
			return err
		}
	}
	files, err := ioutil.ReadDir(StorePath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			simpleData, err := data.LoadSimpleDataFile(filepath.Join(StorePath, file.Name()))
			if err != nil {
				return err
			}
			// TODO: why not a method
			addToStore(simpleData)
		}
	}
	return nil
}

func addToStore(s *data.SimpleData) {
	log.Println("Loading", s.UUID, "...")
	MetaStore.records[string(s.UUID)] = s
}

func (s *store) Write(in io.Reader, data *data.SimpleData) error {
	err := writeRawFile(in, data)
	if err != nil {
		return err
	}

	err = writeSimpleData(data)
	if err != nil {
		return err
	}
	return nil
}

func writeSimpleData(data *data.SimpleData) error {
	writer, err := os.Create(path.Join(StorePath, fmt.Sprintf("%s.json", data.UUID)))
	defer writer.Close()
	if err != nil {
		return nil
	}

	size, err := writer.Write(sfile.MarshalSimpleData(data))
	if err != nil {
		return err
	}

	fmt.Println("Saved", size, "size data file")
	return nil
}
func writeRawFile(in io.Reader, data *data.SimpleData) error {
	writer, err := os.Create(path.Join(StorePath, fmt.Sprintf("%s.data", data.UUID)))
	defer writer.Close()
	if err != nil {
		return err
	}

	size, err := io.Copy(writer, in)
	if err != nil {
		return err
	}

	fmt.Println("Saved", size, "size data file")
	return nil
}

func (s *store) Add(in *os.File) (*data.SimpleData, error) {
	m := data.New(in)

	// Reset the file
	_, err := in.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	err = s.Write(in, m)
	if err != nil {
		return nil, err
	}

	addToStore(m)
	return m, nil
}

var MetaStore store

func init() {
	MetaStore = store{
		records: make(map[string]*data.SimpleData),
	}
	if err := MetaStore.Reload(); err != nil {
		log.Fatal(err)
	}
}
