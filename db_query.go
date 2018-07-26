package ip2location

import (
	"errors"
	"math/big"
	"net"
	"strconv"
)

// Query поиск местоположения по IP
func (db *DB) query(ip string, mode uint32) (*Record, error) {
	record := &Record{}

	// check IP type and return IP number & index (if exists)
	ipType, ipNo, ipIndex := db.checkIP(ip)

	if ipType == 0 {
		return nil, errors.New("invalid IP address")
	}

	var colSize uint32
	var baseAddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowOffset uint32
	var rowOffset2 uint32
	var err error

	ipfrom := big.NewInt(0)
	ipto := big.NewInt(0)
	maxip := big.NewInt(0)

	if ipType == 4 {
		baseAddr = db.meta.ipv4DataBaseAddr
		high = db.meta.ipv4DataBaseCount
		maxip = maxIPv4Range
		colSize = db.meta.ipv4ColumnSize
	} else {
		baseAddr = db.meta.ipv6DataBaseAddr
		high = db.meta.ipv6DataBaseCount
		maxip = maxIPv6Range
		colSize = db.meta.ipv6ColumnSize
	}

	// reading index
	if ipIndex > 0 {
		low, err = db.readUint32(ipIndex)
		high, err = db.readUint32(ipIndex + 4)
	}

	if err != nil {
		return nil, err
	}

	if ipNo.Cmp(maxip) >= 0 {
		ipNo = ipNo.Sub(ipNo, big.NewInt(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowOffset = baseAddr + (mid * colSize)
		rowOffset2 = rowOffset + colSize

		if ipType == 4 {
			fromOffset, err := db.readUint32(rowOffset)
			if err != nil {
				return nil, err
			}
			ipfrom = big.NewInt(int64(fromOffset))

			toOffset, err := db.readUint32(rowOffset2)
			if err != nil {
				return nil, err
			}
			ipto = big.NewInt(int64(toOffset))
		} else {
			ipfrom, err = db.readUint128(rowOffset)
			ipto, err = db.readUint128(rowOffset2)

			if err != nil {
				return nil, err
			}
		}

		if ipNo.Cmp(ipfrom) >= 0 && ipNo.Cmp(ipto) < 0 {
			if ipType == 6 {
				rowOffset = rowOffset + 12 // coz below is assuming All columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			if mode&CountryShort == 1 && db.fieldEnabled.countryEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.countryPositionOffset)
				record.CountryShort, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&CountryLong != 0 && db.fieldEnabled.countryEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.countryPositionOffset)
				record.CountryLong, err = db.readStr(offset + 3)

				if err != nil {
					return nil, err
				}
			}

			if mode&Region != 0 && db.fieldEnabled.regionEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.regionPositionOffset)
				record.Region, err = db.readStr(offset)

				if err != nil {
					return nil, err
				}

			}

			if mode&City != 0 && db.fieldEnabled.cityEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.cityPositionOffset)
				record.City, err = db.readStr(offset)

				if err != nil {
					return nil, err
				}
			}

			if mode&ISP != 0 && db.fieldEnabled.ispEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.ispPositionOffset)
				record.ISP, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Latitude != 0 && db.fieldEnabled.latitudeEnabled {
				record.Latitude, err = db.readFloat(rowOffset + db.offsets.latitudePositionOffset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Longitude != 0 && db.fieldEnabled.longitudeEnabled {
				record.Longitude, err = db.readFloat(rowOffset + db.offsets.longitudePositionOffset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Domain != 0 && db.fieldEnabled.domainEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.domainPositionOffset)
				record.Domain, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&ZipCode != 0 && db.fieldEnabled.zipCodeEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.zipCodePositionOffset)
				record.ZipCode, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&TimeZone != 0 && db.fieldEnabled.timeZoneEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.timeZonePositionOffset)
				record.TimeZone, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&NetSpeed != 0 && db.fieldEnabled.netSpeedEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.netSpeedPositionOffset)
				record.NetSpeed, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&IddCode != 0 && db.fieldEnabled.iddCodeEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.iddCodePositionOffset)
				record.IddCode, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&AreaCode != 0 && db.fieldEnabled.areaCodeEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.areaCodePositionOffset)
				record.AreaCode, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&WeatherStationCode != 0 && db.fieldEnabled.weatherStationCodeEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.weatherStationCodePositionOffset)
				record.WeatherStationCode, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&WeatherStationName != 0 && db.fieldEnabled.weatherStationNameEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.weatherStationNamePositionOffset)
				record.WeatherStationName, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Mcc != 0 && db.fieldEnabled.mccEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.mccPositionOffset)
				record.Mcc, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Mnc != 0 && db.fieldEnabled.mncEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.mncPositionOffset)
				record.Mnc, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&MobileBrand != 0 && db.fieldEnabled.mobileBrandEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.mobileBrandPositionOffset)
				record.MobileBrand, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Elevation != 0 && db.fieldEnabled.elevationEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.elevationPositionOffset)
				float, err := db.readStr(offset)
				if err != nil {
					return nil, err
				}
				f, _ := strconv.ParseFloat(float, 32)
				record.Elevation = float32(f)
			}

			if mode&UsageType != 0 && db.fieldEnabled.usageTypeEnabled {
				offset, err := db.readUint32(rowOffset + db.offsets.usageTypePositionOffset)
				record.UsageType, err = db.readStr(offset)
				if err != nil {
					return nil, err
				}
			}

			return record, nil
		}
		if ipNo.Cmp(ipfrom) < 0 {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return nil, errors.New("no match found")
}

// get IP type and calculate IP number; calculates index too if exists
func (db *DB) checkIP(ip string) (ipType uint32, ipNum *big.Int, ipIndex uint32) {
	ipType = 0
	ipNum = big.NewInt(0)
	ipNumTmp := big.NewInt(0)
	ipIndex = 0
	ipAddress := net.ParseIP(ip)

	if ipAddress != nil {
		v4 := ipAddress.To4()

		if v4 != nil {
			ipType = 4
			ipNum.SetBytes(v4)
		} else {
			v6 := ipAddress.To16()

			if v6 != nil {
				ipType = 6
				ipNum.SetBytes(v6)
			}
		}
	}
	if ipType == 4 {
		if db.meta.ipv4IndexBaseAddr > 0 {
			ipNumTmp.Rsh(ipNum, 16)
			ipNumTmp.Lsh(ipNumTmp, 3)
			ipIndex = uint32(ipNumTmp.Add(ipNumTmp, big.NewInt(int64(db.meta.ipv4IndexBaseAddr))).Uint64())
		}
	} else if ipType == 6 {
		if db.meta.ipv6IndexBaseAddr > 0 {
			ipNumTmp.Rsh(ipNum, 112)
			ipNumTmp.Lsh(ipNumTmp, 3)
			ipIndex = uint32(ipNumTmp.Add(ipNumTmp, big.NewInt(int64(db.meta.ipv6IndexBaseAddr))).Uint64())
		}
	}
	return
}
