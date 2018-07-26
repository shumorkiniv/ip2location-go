package ip2location

import "io"

func (db *DB) readMeta(r io.ReaderAt) error {
	var err error

	db.r = r
	db.meta.dataBaseType, err = db.readuint8(1)
	db.meta.dataBaseColumn, err = db.readuint8(2)
	db.meta.dataBaseYear, err = db.readuint8(3)
	db.meta.dataBaseMonth, err = db.readuint8(4)
	db.meta.dataBaseDay, err = db.readuint8(5)
	db.meta.ipv4DataBaseCount, err = db.readuint32(6)
	db.meta.ipv4DataBaseAddr, err = db.readuint32(10)
	db.meta.ipv6DataBaseCount, err = db.readuint32(14)
	db.meta.ipv6DataBaseAddr, err = db.readuint32(18)
	db.meta.ipv4IndexBaseAddr, err = db.readuint32(22)
	db.meta.ipv6IndexBaseAddr, err = db.readuint32(26)
	db.meta.ipv4ColumnSize = uint32(db.meta.dataBaseColumn << 2)              // 4 bytes each column
	db.meta.ipv6ColumnSize = uint32(16 + ((db.meta.dataBaseColumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	if err != nil {
		return err
	}

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

	return nil
}
