package ip2location

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var DBpath = "IP-COUNTRY-SAMPLE.BIN"

func TestNewLocation(t *testing.T) {
	t.Run("Read from file", func(t *testing.T) {
		db, err := NewFileDB(DBpath, false)
		require.NoError(t, err)

		r, err := db.Query("8.8.8.8", All)
		require.NoError(t, err)

		if r.CountryLong != "UNITED STATES" {
			t.Error("County name is not equal. Expected UNITED STATES, got", r.CountryLong)
		}

		err = db.Close()
		require.NoError(t, err)
	})

	t.Run("Read from memory", func(t *testing.T) {
		db, err := NewFileDB(DBpath, true)
		require.NoError(t, err)

		r, err := db.Query("8.8.8.8", CountryLong)
		require.NoError(t, err)

		if r.CountryLong != "UNITED STATES" {
			t.Error("County name is not equal. Expected UNITED STATES, got", r.CountryLong)
		}

		err = db.Close()
		require.NoError(t, err)
	})

	t.Run("Wrong path", func(t *testing.T) {
		_, err := NewFileDB("", false)
		require.Error(t, err)
	})

	t.Run("Dir", func(t *testing.T) {
		_, err := NewFileDB(".", false)
		require.Error(t, err)
	})
}
