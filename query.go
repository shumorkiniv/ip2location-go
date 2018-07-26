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
		low = readuint32(db.r, ipindex)
		high = readuint32(db.r, ipindex+4)
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(ipno, big.NewInt(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowoffset = baseaddr + (mid * colsize)
		rowoffset2 = rowoffset + colsize

		if iptype == 4 {
			ipfrom = big.NewInt(int64(readuint32(db.r, rowoffset)))
			ipto = big.NewInt(int64(readuint32(db.r, rowoffset2)))
		} else {
			ipfrom = readuint128(db.r, rowoffset)
			ipto = readuint128(db.r, rowoffset2)
		}

		if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
			if iptype == 6 {
				rowoffset = rowoffset + 12 // coz below is assuming All columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			if mode&CountryShort == 1 && countryEnabled {
				record.CountryShort = readstr(db.r, readuint32(db.r, rowoffset+countryPositionOffset))
			}

			if mode&CountryLong != 0 && countryEnabled {
				record.CountryLong = readstr(db.r, readuint32(db.r, rowoffset+countryPositionOffset)+3)
			}

			if mode&Region != 0 && regionEnabled {
				record.Region = readstr(db.r, readuint32(db.r, rowoffset+regionPositionOffset))
			}

			if mode&City != 0 && cityEnabled {
				record.City = readstr(db.r, readuint32(db.r, rowoffset+cityPositionOffset))
			}

			if mode&ISP != 0 && ispEnabled {
				record.ISP = readstr(db.f, readuint32(db.r, rowoffset+ispPositionOffset))
			}

			if mode&Latitude != 0 && latitudeEnabled {
				record.Latitude = readfloat(db.f, rowoffset+latitudePositionOffset)
			}

			if mode&Longitude != 0 && longitudeEnabled {
				record.Longitude = readfloat(db.f, rowoffset+longitudePositionOffset)
			}

			if mode&Domain != 0 && domainEnabled {
				record.Domain = readstr(db.f, readuint32(db.f, rowoffset+domainPositionOffset))
			}

			if mode&ZipCode != 0 && zipCodeEnabled {
				record.ZipCode = readstr(db.f, readuint32(db.f, rowoffset+zipCodePositionOffset))
			}

			if mode&TimeZone != 0 && timeZoneEnabled {
				record.TimeZone = readstr(db.f, readuint32(db.f, rowoffset+timeZonePositionOffset))
			}

			if mode&NetSpeed != 0 && netSpeedEnabled {
				record.NetSpeed = readstr(db.f, readuint32(db.f, rowoffset+netSpeedPositionOffset))
			}

			if mode&IddCode != 0 && iddCodeEnabled {
				record.IddCode = readstr(db.f, readuint32(db.f, rowoffset+iddCodePositionOffset))
			}

			if mode&AreaCode != 0 && areaCodeEnabled {
				record.AreaCode = readstr(db.f, readuint32(db.f, rowoffset+areaCodePositionOffset))
			}

			if mode&WeatherStationCode != 0 && weatherStationCodeEnabled {
				record.WeatherStationCode = readstr(db.f, readuint32(db.f, rowoffset+weatherStationCodePositionOffset))
			}

			if mode&WeatherStationName != 0 && weatherStationNameEnabled {
				record.WeatherStationName = readstr(db.f, readuint32(db.f, rowoffset+weatherStationNamePositionOffset))
			}

			if mode&Mcc != 0 && mccEnabled {
				record.Mcc = readstr(db.f, readuint32(db.f, rowoffset+mccPositionOffset))
			}

			if mode&Mnc != 0 && mncEnabled {
				record.Mnc = readstr(db.f, readuint32(db.f, rowoffset+mncPositionOffset))
			}

			if mode&MobileBrand != 0 && mobileBrandEnabled {
				record.MobileBrand = readstr(db.f, readuint32(db.f, rowoffset+mobileBrandPositionOffset))
			}

			if mode&Elevation != 0 && elevationEnabled {
				f, _ := strconv.ParseFloat(readstr(db.f, readuint32(db.f, rowoffset+elevationPositionOffset)), 32)
				record.Elevation = float32(f)
			}

			if mode&UsageType != 0 && usageTypeEnabled {
				record.UsageType = readstr(db.f, readuint32(db.f, rowoffset+usageTypePositionOffset))
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
func checkip(ip string, meta ip2LocationMeta) (iptype uint32, ipnum *big.Int, ipindex uint32) {
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
