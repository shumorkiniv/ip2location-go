package ip2location

import (
	"testing"
	"os"
)

func TestNewDb(t *testing.T) {
	t.Run("Read from file", func(t *testing.T) {
		var f *os.File
		var err error
		if f, err = os.Open(DBpath); err != nil {
			t.Error(err)
		}

		db, err := NewDB(f)
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
