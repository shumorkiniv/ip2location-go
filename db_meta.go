package ip2location

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

func (db *DB) readMeta() error {
	var err error

	if db.meta.dataBaseType, err = db.readUint8(1); err != nil {
		return err
	}
	if db.meta.dataBaseColumn, err = db.readUint8(2); err != nil {
		return err
	}
	if db.meta.dataBaseYear, err = db.readUint8(3); err != nil {
		return err
	}
	if db.meta.dataBaseMonth, err = db.readUint8(4); err != nil {
		return err
	}
	if db.meta.dataBaseDay, err = db.readUint8(5); err != nil {
		return err
	}
	if db.meta.ipv4DataBaseCount, err = db.readUint32(6); err != nil {
		return err
	}
	if db.meta.ipv4DataBaseAddr, err = db.readUint32(10); err != nil {
		return err
	}
	if db.meta.ipv6DataBaseCount, err = db.readUint32(14); err != nil {
		return err
	}
	if db.meta.ipv6DataBaseAddr, err = db.readUint32(18); err != nil {
		return err
	}
	if db.meta.ipv4IndexBaseAddr, err = db.readUint32(22); err != nil {
		return err
	}
	if db.meta.ipv6IndexBaseAddr, err = db.readUint32(26); err != nil {
		return err
	}
	db.meta.ipv4ColumnSize = uint32(db.meta.dataBaseColumn << 2)              // 4 bytes each column
	db.meta.ipv6ColumnSize = uint32(16 + ((db.meta.dataBaseColumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	dbt := db.meta.dataBaseType
	// since both IPv4 and IPv6 use 4 bytes for the below columns, can just do it once here
	if countryPosition[dbt] != 0 {
		db.offsets.countryPositionOffset = uint32(countryPosition[dbt]-1) << 2
		db.fieldEnabled.countryEnabled = true
	}
	if regionPosition[dbt] != 0 {
		db.offsets.regionPositionOffset = uint32(regionPosition[dbt]-1) << 2
		db.fieldEnabled.regionEnabled = true
	}
	if cityPosition[dbt] != 0 {
		db.offsets.cityPositionOffset = uint32(cityPosition[dbt]-1) << 2
		db.fieldEnabled.cityEnabled = true
	}
	if ispPosition[dbt] != 0 {
		db.offsets.ispPositionOffset = uint32(ispPosition[dbt]-1) << 2
		db.fieldEnabled.ispEnabled = true
	}
	if domainPosition[dbt] != 0 {
		db.offsets.domainPositionOffset = uint32(domainPosition[dbt]-1) << 2
		db.fieldEnabled.domainEnabled = true
	}
	if zipCodePosition[dbt] != 0 {
		db.offsets.zipCodePositionOffset = uint32(zipCodePosition[dbt]-1) << 2
		db.fieldEnabled.zipCodeEnabled = true
	}
	if latitudePosition[dbt] != 0 {
		db.offsets.latitudePositionOffset = uint32(latitudePosition[dbt]-1) << 2
		db.fieldEnabled.latitudeEnabled = true
	}
	if longitudePosition[dbt] != 0 {
		db.offsets.longitudePositionOffset = uint32(longitudePosition[dbt]-1) << 2
		db.fieldEnabled.longitudeEnabled = true
	}
	if timeZonePosition[dbt] != 0 {
		db.offsets.timeZonePositionOffset = uint32(timeZonePosition[dbt]-1) << 2
		db.fieldEnabled.timeZoneEnabled = true
	}
	if netSpeedPosition[dbt] != 0 {
		db.offsets.netSpeedPositionOffset = uint32(netSpeedPosition[dbt]-1) << 2
		db.fieldEnabled.netSpeedEnabled = true
	}
	if iddCodePosition[dbt] != 0 {
		db.offsets.iddCodePositionOffset = uint32(iddCodePosition[dbt]-1) << 2
		db.fieldEnabled.iddCodeEnabled = true
	}
	if areaCodePosition[dbt] != 0 {
		db.offsets.areaCodePositionOffset = uint32(areaCodePosition[dbt]-1) << 2
		db.fieldEnabled.areaCodeEnabled = true
	}
	if weatherStationCodePosition[dbt] != 0 {
		db.offsets.weatherStationCodePositionOffset = uint32(weatherStationCodePosition[dbt]-1) << 2
		db.fieldEnabled.weatherStationCodeEnabled = true
	}
	if weatherStationNamePosition[dbt] != 0 {
		db.offsets.weatherStationNamePositionOffset = uint32(weatherStationNamePosition[dbt]-1) << 2
		db.fieldEnabled.weatherStationNameEnabled = true
	}
	if mccPosition[dbt] != 0 {
		db.offsets.mccPositionOffset = uint32(mccPosition[dbt]-1) << 2
		db.fieldEnabled.mccEnabled = true
	}
	if mncPosition[dbt] != 0 {
		db.offsets.mncPositionOffset = uint32(mncPosition[dbt]-1) << 2
		db.fieldEnabled.mncEnabled = true
	}
	if mobileBrandPosition[dbt] != 0 {
		db.offsets.mobileBrandPositionOffset = uint32(mobileBrandPosition[dbt]-1) << 2
		db.fieldEnabled.mobileBrandEnabled = true
	}
	if elevationPosition[dbt] != 0 {
		db.offsets.elevationPositionOffset = uint32(elevationPosition[dbt]-1) << 2
		db.fieldEnabled.elevationEnabled = true
	}
	if usageTypePosition[dbt] != 0 {
		db.offsets.usageTypePositionOffset = uint32(usageTypePosition[dbt]-1) << 2
		db.fieldEnabled.usageTypeEnabled = true
	}

	return nil
}
