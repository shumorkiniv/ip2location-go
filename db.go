package ip2location

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/big"
)

// DB содержит в себе данные из файла
type DB struct {
	r            io.ReaderAt
	meta         meta
	offsets      dbOffset
	fieldEnabled fieldEnabled
}

type meta struct {
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

type dbOffset struct {
	countryPositionOffset            uint32
	regionPositionOffset             uint32
	cityPositionOffset               uint32
	ispPositionOffset                uint32
	domainPositionOffset             uint32
	zipCodePositionOffset            uint32
	latitudePositionOffset           uint32
	longitudePositionOffset          uint32
	timeZonePositionOffset           uint32
	netSpeedPositionOffset           uint32
	iddCodePositionOffset            uint32
	areaCodePositionOffset           uint32
	weatherStationCodePositionOffset uint32
	weatherStationNamePositionOffset uint32
	mccPositionOffset                uint32
	mncPositionOffset                uint32
	mobileBrandPositionOffset        uint32
	elevationPositionOffset          uint32
	usageTypePositionOffset          uint32
}

type fieldEnabled struct {
	countryEnabled            bool
	regionEnabled             bool
	cityEnabled               bool
	ispEnabled                bool
	domainEnabled             bool
	zipCodeEnabled            bool
	latitudeEnabled           bool
	longitudeEnabled          bool
	timeZoneEnabled           bool
	netSpeedEnabled           bool
	iddCodeEnabled            bool
	areaCodeEnabled           bool
	weatherStationCodeEnabled bool
	weatherStationNameEnabled bool
	mccEnabled                bool
	mncEnabled                bool
	mobileBrandEnabled        bool
	elevationEnabled          bool
	usageTypeEnabled          bool
}

func newDb(r io.ReaderAt) (*DB, error) {
	db := &DB{}
	err := db.readMeta(r)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// read byte
func (db *DB) readUint8(pos int64) (uint8, error) {
	var retval uint8
	data := make([]byte, 1)
	_, err := db.r.ReadAt(data, pos-1)
	if err != nil {
		return 0, err
	}
	retval = data[0]
	return retval, nil
}

// read unsigned 32-bit integer
func (db *DB) readUint32(pos uint32) (uint32, error) {
	pos2 := int64(pos)
	var retval uint32
	data := make([]byte, 4)
	_, err := db.r.ReadAt(data, pos2-1)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		return 0, err
	}
	return retval, nil
}

// read unsigned 128-bit integer
func (db *DB) readUint128(pos uint32) (*big.Int, error) {
	pos2 := int64(pos)
	retval := big.NewInt(0)
	data := make([]byte, 16)
	_, err := db.r.ReadAt(data, pos2-1)
	if err != nil {
		return retval, err
	}

	// little endian to big endian
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	retval.SetBytes(data)
	return retval, nil
}

// read string
func (db *DB) readStr(pos uint32) (string, error) {
	pos2 := int64(pos)
	var retval string
	lenbyte := make([]byte, 1)
	_, err := db.r.ReadAt(lenbyte, pos2)
	if err != nil {
		return "", err
	}
	strlen := lenbyte[0]
	data := make([]byte, strlen)
	_, err = db.r.ReadAt(data, pos2+1)
	if err != nil {
		return "", err
	}
	retval = string(data[:strlen])
	return retval, nil
}

// read float
func (db *DB) readFloat(pos uint32) (float32, error) {
	pos2 := int64(pos)
	var retval float32
	data := make([]byte, 4)
	_, err := db.r.ReadAt(data, pos2-1)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		return 0, nil
	}
	return retval, nil
}
