package sfile

import (
	"log"
	"os"
	"encoding/json"
	"github.com/vbabiy/simple/simple/data"
	"path/filepath"
	"strings"
)

func WriteSimpleFile(outputName string, meta *data.SimpleData) {
	log.Println("Saving simple file...", outputName)
	out, err := os.Create(outputName)
	defer out.Close()
	if err != nil {
		log.Fatal(err)
	}
	out.Write(MarshalSimpleData(meta))
}

func MarshalSimpleData(meta *data.SimpleData) []byte {

	if b, err := json.Marshal(meta); err != nil {
		log.Fatal(err)
		return nil
	} else {
		return b
	}
}

func SwapExt(inputName string) string {
	baseFilename := filepath.Base(inputName)
	name := strings.Replace(baseFilename, filepath.Ext(baseFilename), "", 1)
	return strings.Join([]string{name, ".simple"}, "")
}
