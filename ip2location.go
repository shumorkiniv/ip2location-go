package ip2location

import (
	"bytes"
	"errors"
	"io"
	"os"
	"syscall"
)

type Location struct {
	f    *os.File
	data *DB
}

//const api_version string = "8.0.3"

// New читает файл
func New(path string, mmap bool) (*Location, error) {
	var err error
	s, err := os.Stat(path)
	if err != nil {
		return nil, errors.New("wrong path")
	}

	if s.IsDir() {
		return nil, errors.New("can't open directory")
	}

	data := &Location{}
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
		data.data, err = NewDb(r)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (location *Location) Query(ip string, mode uint32) (*Record, error) {
	return location.data.Query(ip, mode)
}

// Close закрывает файл БД
func (location *Location) Close() error {
	if location.f == nil {
		return nil
	}
	err := location.f.Close()
	if err != nil {
		return err
	}
	return nil
}
