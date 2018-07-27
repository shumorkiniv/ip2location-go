package ip2location

import (
	"testing"
	"io"
	"os"
	"syscall"
	"bytes"
)

func TestNewDb(t *testing.T) {
	t.Run("Read from memory", func(t *testing.T) {
		var r io.ReaderAt
		var err error

		s, err := os.Stat(DBpath)
		if err != nil {
			t.Error(err)
		}

		var fd int
		if fd, err = syscall.Open(DBpath, syscall.O_RDONLY, 0); err != nil {
			t.Error(err)
		}

		var data []byte
		if data, err = syscall.Mmap(fd, 0, int(s.Size()), syscall.PROT_READ, syscall.MAP_SHARED); err != nil {
			t.Error(err)
		}

		r = bytes.NewReader(data)
		db, err := NewDb(r)
		if err != nil {
			t.Error(err)
		}

		record, err := db.Query("8.8.8.8", All)
		if err != nil {
			t.Error(err)
		}

		if record.CountryLong != "UNITED STATES" {
			t.Error("County name is not equal. Expected UNITED STATES, got", record.CountryLong)
		}

		db.Close()
	})

	t.Run("Read from file", func(t *testing.T) {
		var f *os.File
		var err error
		var r io.ReaderAt
		if f, err = os.Open(DBpath); err != nil {
			t.Error(err)
		}
		r = f

		db, err := NewDb(r)
		if err != nil {
			t.Error(err)
		}

		record, err := db.Query("8.8.8.8", All)
		if err != nil {
			t.Error(err)
		}

		if record.CountryLong != "UNITED STATES" {
			t.Error("County name is not equal. Expected UNITED STATES, got", record.CountryLong)
		}

		db.Close()
	})
}
