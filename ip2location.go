package ip2location

import (
	"bytes"
	"io"
	"math/big"
	"os"
	"syscall"

	"github.com/pkg/errors"
)

// DB содержит в себе данные из файла
type DB struct {
	f    *os.File
	r    io.ReaderAt
	meta ip2LocationMeta
}

type ip2LocationMeta struct {
	dataBaseType      uint8
	dataBaseColumn    uint8
	dataBaseDay       uint8
	dataBaseMonth     uint8
	dataBaseYear      uint8
	ipv4DataBaseCount uint32
	ipv4DataBaseAddr  uint32
	ipv6DataBaseCount uint32
	ipv6DataBaseAddr  uint32
	ipv4IndexBaseAddr uint32
	ipv6IndexBaseAddr uint32
	ipv4ColumnSize    uint32
	ipv6ColumnSize    uint32
}

var countryPosition = [25]uint8{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var regionPosition = [25]uint8{0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
var cityPosition = [25]uint8{0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
var ispPosition = [25]uint8{0, 0, 3, 0, 5, 0, 7, 5, 7, 0, 8, 0, 9, 0, 9, 0, 9, 0, 9, 7, 9, 0, 9, 7, 9}
var latitudePosition = [25]uint8{0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
var longitudePosition = [25]uint8{0, 0, 0, 0, 0, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6}
var domainPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 6, 8, 0, 9, 0, 10, 0, 10, 0, 10, 0, 10, 8, 10, 0, 10, 8, 10}
var zipCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 7, 7, 7, 0, 7, 0, 7, 7, 7, 0, 7}
var timeZonePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 7, 8, 8, 8, 7, 8, 0, 8, 8, 8, 0, 8}
var netSpeedPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 11, 0, 11, 8, 11, 0, 11, 0, 11, 0, 11}
var iddCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 12, 0, 12, 0, 12, 9, 12, 0, 12}
var areaCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 13, 0, 13, 0, 13, 10, 13, 0, 13}
var weatherStationCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 14, 0, 14, 0, 14, 0, 14}
var weatherStationNamePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 15, 0, 15, 0, 15, 0, 15}
var mccPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 16, 0, 16, 9, 16}
var mncPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 17, 0, 17, 10, 17}
var mobileBrandPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 18, 0, 18, 11, 18}
var elevationPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 19, 0, 19}
var usageTypePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 20}

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
		db.r = r

		db.meta.dataBaseType = readuint8(db.r, 1)
		db.meta.dataBaseColumn = readuint8(db.r, 2)
		db.meta.dataBaseYear = readuint8(db.r, 3)
		db.meta.dataBaseMonth = readuint8(db.r, 4)
		db.meta.dataBaseDay = readuint8(db.r, 5)
		db.meta.ipv4DataBaseCount = readuint32(db.r, 6)
		db.meta.ipv4DataBaseAddr = readuint32(db.r, 10)
		db.meta.ipv6DataBaseCount = readuint32(db.r, 14)
		db.meta.ipv6DataBaseAddr = readuint32(db.r, 18)
		db.meta.ipv4IndexBaseAddr = readuint32(db.r, 22)
		db.meta.ipv6IndexBaseAddr = readuint32(db.r, 26)
		db.meta.ipv4ColumnSize = uint32(db.meta.dataBaseColumn << 2)              // 4 bytes each column
		db.meta.ipv6ColumnSize = uint32(16 + ((db.meta.dataBaseColumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

		dbt := db.meta.dataBaseType

		// since both IPv4 and IPv6 use 4 bytes for the below columns, can just do it once here
		if countryPosition[dbt] != 0 {
			countryPositionOffset = uint32(countryPosition[dbt]-1) << 2
			countryEnabled = true
		}
		if regionPosition[dbt] != 0 {
			regionPositionOffset = uint32(regionPosition[dbt]-1) << 2
			regionEnabled = true
		}
		if cityPosition[dbt] != 0 {
			cityPositionOffset = uint32(cityPosition[dbt]-1) << 2
			cityEnabled = true
		}
		if ispPosition[dbt] != 0 {
			ispPositionOffset = uint32(ispPosition[dbt]-1) << 2
			ispEnabled = true
		}
		if domainPosition[dbt] != 0 {
			domainPositionOffset = uint32(domainPosition[dbt]-1) << 2
			domainEnabled = true
		}
		if zipCodePosition[dbt] != 0 {
			zipCodePositionOffset = uint32(zipCodePosition[dbt]-1) << 2
			zipCodeEnabled = true
		}
		if latitudePosition[dbt] != 0 {
			latitudePositionOffset = uint32(latitudePosition[dbt]-1) << 2
			latitudeEnabled = true
		}
		if longitudePosition[dbt] != 0 {
			longitudePositionOffset = uint32(longitudePosition[dbt]-1) << 2
			longitudeEnabled = true
		}
		if timeZonePosition[dbt] != 0 {
			timeZonePositionOffset = uint32(timeZonePosition[dbt]-1) << 2
			timeZoneEnabled = true
		}
		if netSpeedPosition[dbt] != 0 {
			netSpeedPositionOffset = uint32(netSpeedPosition[dbt]-1) << 2
			netSpeedEnabled = true
		}
		if iddCodePosition[dbt] != 0 {
			iddCodePositionOffset = uint32(iddCodePosition[dbt]-1) << 2
			iddCodeEnabled = true
		}
		if areaCodePosition[dbt] != 0 {
			areaCodePositionOffset = uint32(areaCodePosition[dbt]-1) << 2
			areaCodeEnabled = true
		}
		if weatherStationCodePosition[dbt] != 0 {
			weatherStationCodePositionOffset = uint32(weatherStationCodePosition[dbt]-1) << 2
			weatherStationCodeEnabled = true
		}
		if weatherStationNamePosition[dbt] != 0 {
			weatherStationNamePositionOffset = uint32(weatherStationNamePosition[dbt]-1) << 2
			weatherStationNameEnabled = true
		}
		if mccPosition[dbt] != 0 {
			mccPositionOffset = uint32(mccPosition[dbt]-1) << 2
			mccEnabled = true
		}
		if mncPosition[dbt] != 0 {
			mncPositionOffset = uint32(mncPosition[dbt]-1) << 2
			mncEnabled = true
		}
		if mobileBrandPosition[dbt] != 0 {
			mobileBrandPositionOffset = uint32(mobileBrandPosition[dbt]-1) << 2
			mobileBrandEnabled = true
		}
		if elevationPosition[dbt] != 0 {
			elevationPositionOffset = uint32(elevationPosition[dbt]-1) << 2
			elevationEnabled = true
		}
		if usageTypePosition[dbt] != 0 {
			usageTypePositionOffset = uint32(usageTypePosition[dbt]-1) << 2
			usageTypeEnabled = true
		}

		metaOk = true
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
