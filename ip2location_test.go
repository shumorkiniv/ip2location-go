package ip2location

import (
	"fmt"
	"testing"
)

var DBpath = "IP-COUNTRY-SAMPLE.BIN"

func TestNew(t *testing.T) {
	t.Run("Read from file", func(t *testing.T) {
		db, err := New(DBpath, false)
		if err != nil {
			t.Error(err)
		}

		r, err := db.Query("8.8.8.8", All)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(r)

		err = db.Close()
		if err != nil {
			t.Error(err)
		}

	})

	t.Run("Read from memory", func(t *testing.T) {
		db, err := New(DBpath, true)
		if err != nil {
			t.Error(err)
		}

		_, err = db.Query("8.8.8.8", CountryLong)
		if err != nil {
			t.Error(err)
		}

		err = db.Close()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Wrong path", func(t *testing.T) {
		_, err := New("", false)
		if err == nil {
			t.Error(err)
		}
	})

	t.Run("Dir", func(t *testing.T) {
		_, err := New(".", false)
		if err == nil {
			t.Error(err)
		}
	})
}
