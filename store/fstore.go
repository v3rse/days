package store

import (
	"encoding/json"
	"os"

	"github.com/v3rse/days/utils"
)

type FileStore struct {
	file    *os.File
	decoder *json.Decoder
	encoder *json.Encoder
}

func NewFileStore(fileName string) FileStore {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	utils.Check(err)

	decoder := json.NewDecoder(file)
	encoder := json.NewEncoder(file)

	store := FileStore{
		file,
		decoder,
		encoder,
	}

	store.init([]byte("{\"start\":null, \"habits\":[], \"end\": null}"))

	return store
}

func (f *FileStore) init(initContent []byte) {
	info, err := f.file.Stat()
	utils.Check(err)

	if info.Size() == 0 {
		f.file.Write(initContent)
		f.file.Seek(0, 0)
	}
}

func (f *FileStore) Load(dest interface{}) {
	err := f.decoder.Decode(&dest)
	utils.Check(err)
}

func (f *FileStore) Save(data interface{}) {
	f.file.Seek(0, 0)
	err := f.encoder.Encode(data)
	utils.Check(err)
}

func (f *FileStore) Close() {
	f.file.Close()
}
