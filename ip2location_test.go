package ip2location

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var DBpath = "IP-COUNTRY-SAMPLE.BIN"

func TestNew(t *testing.T) {
	t.Run("Read from file", func(t *testing.T) {
		db, err := New(DBpath, false)
		assert.NoError(t, err)

		r, err := db.Query("8.8.8.8", All)
		assert.NoError(t, err)

		fmt.Println(r)

		err = db.Close()
		assert.NoError(t, err)
	})

	t.Run("Read from memory", func(t *testing.T) {
		db, err := New(DBpath, true)
		assert.NoError(t, err)

		r, err := db.Query("8.8.8.8", All)
		assert.NoError(t, err)

		fmt.Println(r)

		err = db.Close()
		assert.NoError(t, err)
	})

	t.Run("Wrong path", func(t *testing.T) {
		_, err := New("", false)
		assert.Error(t, err)
	})

	t.Run("Dir", func(t *testing.T) {
		_, err := New(".", false)
		assert.Error(t, err)
	})
}
