package ip2location

import (
	"bytes"
	"errors"
	"io"
	"math/big"
	"os"
	"syscall"
)

//const api_version string = "8.0.3"

var maxIPv4Range = big.NewInt(4294967295)
var maxIPv6Range = big.NewInt(0)

var metaOk bool

var countryPositionOffset uint32
var regionPositionOffset uint32
var cityPositionOffset uint32
var ispPositionOffset uint32
var domainPositionOffset uint32
var zipCodePositionOffset uint32
var latitudePositionOffset uint32
var longitudePositionOffset uint32
var timeZonePositionOffset uint32
var netSpeedPositionOffset uint32
var iddCodePositionOffset uint32
var areaCodePositionOffset uint32
var weatherStationCodePositionOffset uint32
var weatherStationNamePositionOffset uint32
var mccPositionOffset uint32
var mncPositionOffset uint32
var mobileBrandPositionOffset uint32
var elevationPositionOffset uint32
var usageTypePositionOffset uint32

var countryEnabled bool
var regionEnabled bool
var cityEnabled bool
var ispEnabled bool
var domainEnabled bool
var zipCodeEnabled bool
var latitudeEnabled bool
var longitudeEnabled bool
var timeZoneEnabled bool
var netSpeedEnabled bool
var iddCodeEnabled bool
var areaCodeEnabled bool
var weatherStationCodeEnabled bool
var weatherStationNameEnabled bool
var mccEnabled bool
var mncEnabled bool
var mobileBrandEnabled bool
var elevationEnabled bool
var usageTypeEnabled bool

// New читает файл
func New(path string, mmap bool) (*DB, error) {
	var err error
	s, err := os.Stat(path)
	if err != nil {
		return nil, errors.New("Wrong path")
	}

	if s.IsDir() {
		return nil, errors.New("Can't open directory")
	}

	db := &DB{}
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
		if db.f, err = os.Open(path); err != nil {
			return nil, err
		}
		r = db.f
	}

	if r != nil {
		err := db.readMeta(r)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Close закрывает файл БД
func (db *DB) Close() error {
	// если читали из памяти, то ничего не закрываем
	if db.f == nil {
		return nil
	}
	err := db.f.Close()
	if err != nil {
		return err
	}
	return nil
}
