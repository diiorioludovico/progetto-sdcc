package query

//server.go

func GetMeasure() string {
	return "INSERT INTO measures(sensor_id, park_id, temperature, humidity, brightness, air_quality, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?)"
}

func GetSensor() string {
	return "SELECT * FROM sensors WHERE serial_number = ?"
}

// OPERATIONS//
func UpdateSensorStatus() string {
	return "UPDATE sensors SET is_active = ? WHERE id = ?"
}
