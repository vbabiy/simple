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

type Store struct {
	records map[string]*data.SimpleData
}

func New() (*Store, error) {
	s := &Store{
		records: make(map[string]*data.SimpleData),
	}
	if err := s.Reload(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) All() []*data.SimpleData {
	out := make([]*data.SimpleData, len(s.records))
	var idx int
	for _, v := range s.records {
		out[idx] = v
		idx++
	}
	return out
}

func (s *Store) Reload() error {
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
			s.addToStore(simpleData)
		}
	}
	return nil
}

func (s *Store) addToStore(d *data.SimpleData) {
	log.Println("Loading", d.UUID, "...")
	s.records[string(d.UUID)] = d
}

func (s *Store) Write(in io.Reader, data *data.SimpleData) error {
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

func (s *Store) Add(in *os.File) (*data.SimpleData, error) {
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

	s.addToStore(m)
	return m, nil
}
