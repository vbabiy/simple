package store

import (
	"github.com/vbabiy/simple/simple/data"
	"io/ioutil"
	"fmt"
	"log"
	"os"
	"io"
	"path"
	"github.com/vbabiy/simple/simple/sfile"
	"strings"
	"path/filepath"
)

type store struct {
	records map[string]*data.SimpleData
}

const StorePath = ".store"

func (s *store) All() ([]*data.SimpleData) {
	out := make([]*data.SimpleData, len(s.records))
	var idx int
	for _, v := range s.records {
		out[idx] = v
		idx++
	}

	return out
}

func (s *store) Reload() (error) {
	s.records = make(map[string]*data.SimpleData) // reset the map
	if files, err := ioutil.ReadDir(StorePath); err != nil {
		return err
	} else {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".json") {
				simpleData, err := data.LoadSimpleDataFile(filepath.Join(StorePath, file.Name()))
				if err != nil {
					log.Fatal(err)
					continue
				}
				addToStore(simpleData)
			}
		}
	}

	return nil
}
func addToStore(s *data.SimpleData) {
	log.Println("Loading", s.UUID, "...")
	MetaStore.records[string(s.UUID)] = s
}

func (s *store) Write(in io.Reader, data *data.SimpleData) (error) {
	err := writeRawFile(in, data)
	if err != nil {
		log.Fatal(err)
	}
	err = writeSimpleData(data)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func writeSimpleData(data *data.SimpleData) error {
	writer, err := os.Create(path.Join(StorePath, fmt.Sprintf("%s.json", data.UUID)))
	defer writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	size, err := writer.Write(sfile.MarshalSimpleData(data))
	log.Println("Saved", size, "size data file")
	return err
}
func writeRawFile(in io.Reader, data *data.SimpleData) error {
	writer, err := os.Create(path.Join(StorePath, fmt.Sprintf("%s.data", data.UUID)))
	defer writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	size, err := io.Copy(writer, in)
	log.Println("Saved", size, "size data file")
	return err
}

func (s *store) Add(in *os.File) (*data.SimpleData, error) {
	m := data.New(in)

	// Reset the file
	_, err := in.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Write(in, m)
	if err != nil {
		log.Fatal(err)
	}

	addToStore(m)
	return m, nil
}

var MetaStore store

func init() {
	MetaStore = store{
		records:make(map[string]*data.SimpleData),
	}
	if err := MetaStore.Reload(); err != nil {
		log.Fatal(err)
	}
}
