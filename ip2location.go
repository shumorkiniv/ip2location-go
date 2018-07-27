package ip2location

import (
	"bytes"
	"errors"
	"io"
	"os"
	"syscall"
)

type FileDB struct {
	f    *os.File
	data *DB
}

// New читает файл
func New(path string, mmap bool) (*FileDB, error) {
	var err error
	s, err := os.Stat(path)
	if err != nil {
		return nil, errors.New("wrong path")
	}

	if s.IsDir() {
		return nil, errors.New("can't open directory")
	}

	data := &FileDB{}
	var r io.ReaderAt
	if mmap {
		var fd int
		if fd, err = syscall.Open(path, syscall.O_RDONLY, 0); err != nil {
			return nil, err
		}

		var data []byte
		if data, err = syscall.Mmap(fd, 0, int(s.Size()), syscall.PROT_READ, syscall.MAP_SHARED); err != nil {
			return nil, err
		}

		r = bytes.NewReader(data)
	} else {
		if data.f, err = os.Open(path); err != nil {
			return nil, err
		}
		r = data.f
	}

	if r != nil {
		data.data, err = NewDB(r)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (fileDB *FileDB) Query(ip string, mode uint32) (*Record, error) {
	return fileDB.data.Query(ip, mode)
}

// Close закрывает файл БД
func (fileDB *FileDB) Close() error {
	if fileDB.f == nil {
		return nil
	}
	err := fileDB.f.Close()
	if err != nil {
		return err
	}
	return nil
}
