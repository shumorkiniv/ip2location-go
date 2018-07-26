package ip2location

import (
	"errors"
	"math/big"
	"net"
	"strconv"
)

// Query поиск местоположения по IP
func (db *DB) Query(ip string, mode uint32) (*Record, error) {
	record := &Record{}

	// read metadata
	if !metaOk {
		return nil, errors.New("file not found")
	}

	// check IP type and return IP number & index (if exists)
	iptype, ipno, ipindex := checkip(ip, db.meta)

	if iptype == 0 {
		return nil, errors.New("invalid IP address")
	}

	var colsize uint32
	var baseaddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowoffset uint32
	var rowoffset2 uint32
	var err error

	ipfrom := big.NewInt(0)
	ipto := big.NewInt(0)
	maxip := big.NewInt(0)

	if iptype == 4 {
		baseaddr = db.meta.ipv4DataBaseAddr
		high = db.meta.ipv4DataBaseCount
		maxip = maxIPv4Range
		colsize = db.meta.ipv4ColumnSize
	} else {
		baseaddr = db.meta.ipv6DataBaseAddr
		high = db.meta.ipv6DataBaseCount
		maxip = maxIPv6Range
		colsize = db.meta.ipv6ColumnSize
	}

	// reading index
	if ipindex > 0 {
		low, err = db.readuint32(ipindex)
		high, err = db.readuint32(ipindex + 4)
	}

	if err != nil {
		return nil, err
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(ipno, big.NewInt(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowoffset = baseaddr + (mid * colsize)
		rowoffset2 = rowoffset + colsize

		if iptype == 4 {
			fromOffset, err := db.readuint32(rowoffset)
			if err != nil {
				return nil, err
			}
			ipfrom = big.NewInt(int64(fromOffset))

			toOffset, err := db.readuint32(rowoffset2)
			if err != nil {
				return nil, err
			}
			ipto = big.NewInt(int64(toOffset))
		} else {
			ipfrom, err = db.readuint128(rowoffset)
			ipto, err = db.readuint128(rowoffset2)

			if err != nil {
				return nil, err
			}
		}

		if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
			if iptype == 6 {
				rowoffset = rowoffset + 12 // coz below is assuming All columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			if mode&CountryShort == 1 && countryEnabled {
				offset, err := db.readuint32(rowoffset + countryPositionOffset)
				record.CountryShort, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&CountryLong != 0 && countryEnabled {
				offset, err := db.readuint32(rowoffset + countryPositionOffset)
				record.CountryLong, err = db.readstr(offset + 3)

				if err != nil {
					return nil, err
				}
			}

			if mode&Region != 0 && regionEnabled {
				offset, err := db.readuint32(rowoffset + regionPositionOffset)
				record.Region, err = db.readstr(offset)

				if err != nil {
					return nil, err
				}

			}

			if mode&City != 0 && cityEnabled {
				offset, err := db.readuint32(rowoffset + cityPositionOffset)
				record.City, err = db.readstr(offset)

				if err != nil {
					return nil, err
				}
			}

			if mode&ISP != 0 && ispEnabled {
				offset, err := db.readuint32(rowoffset + ispPositionOffset)
				record.ISP, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Latitude != 0 && latitudeEnabled {
				record.Latitude, err = db.readfloat(rowoffset + latitudePositionOffset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Longitude != 0 && longitudeEnabled {
				record.Longitude, err = db.readfloat(rowoffset + longitudePositionOffset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Domain != 0 && domainEnabled {
				offset, err := db.readuint32(rowoffset + domainPositionOffset)
				record.Domain, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&ZipCode != 0 && zipCodeEnabled {
				offset, err := db.readuint32(rowoffset + zipCodePositionOffset)
				record.ZipCode, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&TimeZone != 0 && timeZoneEnabled {
				offset, err := db.readuint32(rowoffset + timeZonePositionOffset)
				record.TimeZone, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&NetSpeed != 0 && netSpeedEnabled {
				offset, err := db.readuint32(rowoffset + netSpeedPositionOffset)
				record.NetSpeed, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&IddCode != 0 && iddCodeEnabled {
				offset, err := db.readuint32(rowoffset + iddCodePositionOffset)
				record.IddCode, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&AreaCode != 0 && areaCodeEnabled {
				offset, err := db.readuint32(rowoffset + areaCodePositionOffset)
				record.AreaCode, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&WeatherStationCode != 0 && weatherStationCodeEnabled {
				offset, err := db.readuint32(rowoffset + weatherStationCodePositionOffset)
				record.WeatherStationCode, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&WeatherStationName != 0 && weatherStationNameEnabled {
				offset, err := db.readuint32(rowoffset + weatherStationNamePositionOffset)
				record.WeatherStationName, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Mcc != 0 && mccEnabled {
				offset, err := db.readuint32(rowoffset + mccPositionOffset)
				record.Mcc, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Mnc != 0 && mncEnabled {
				offset, err := db.readuint32(rowoffset + mncPositionOffset)
				record.Mnc, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&MobileBrand != 0 && mobileBrandEnabled {
				offset, err := db.readuint32(rowoffset + mobileBrandPositionOffset)
				record.MobileBrand, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			if mode&Elevation != 0 && elevationEnabled {
				offset, err := db.readuint32(rowoffset + elevationPositionOffset)
				float, err := db.readstr(offset)
				if err != nil {
					return nil, err
				}
				f, _ := strconv.ParseFloat(float, 32)
				record.Elevation = float32(f)
			}

			if mode&UsageType != 0 && usageTypeEnabled {
				offset, err := db.readuint32(rowoffset + usageTypePositionOffset)
				record.UsageType, err = db.readstr(offset)
				if err != nil {
					return nil, err
				}
			}

			return record, nil
		}
		if ipno.Cmp(ipfrom) < 0 {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return nil, errors.New("no match found")
}

// get IP type and calculate IP number; calculates index too if exists
func checkip(ip string, meta meta) (iptype uint32, ipnum *big.Int, ipindex uint32) {
	iptype = 0
	ipnum = big.NewInt(0)
	ipnumtmp := big.NewInt(0)
	ipindex = 0
	ipaddress := net.ParseIP(ip)

	if ipaddress != nil {
		v4 := ipaddress.To4()

		if v4 != nil {
			iptype = 4
			ipnum.SetBytes(v4)
		} else {
			v6 := ipaddress.To16()

			if v6 != nil {
				iptype = 6
				ipnum.SetBytes(v6)
			}
		}
	}
	if iptype == 4 {
		if meta.ipv4IndexBaseAddr > 0 {
			ipnumtmp.Rsh(ipnum, 16)
			ipnumtmp.Lsh(ipnumtmp, 3)
			ipindex = uint32(ipnumtmp.Add(ipnumtmp, big.NewInt(int64(meta.ipv4IndexBaseAddr))).Uint64())
		}
	} else if iptype == 6 {
		if meta.ipv6IndexBaseAddr > 0 {
			ipnumtmp.Rsh(ipnum, 112)
			ipnumtmp.Lsh(ipnumtmp, 3)
			ipindex = uint32(ipnumtmp.Add(ipnumtmp, big.NewInt(int64(meta.ipv6IndexBaseAddr))).Uint64())
		}
	}
	return
}
