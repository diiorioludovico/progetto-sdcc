package query

func ShowSensors() string {
	return "SELECT * FROM sensors"
}

func ShowParks() string {
	return "SELECT * FROM parks"
}

func GetSensorParkid() string {
	return "SELECT park_id FROM sensors WHERE id = ?"
}

func GetParkStatus() string {
	return "SELECT is_observed, sensors.id FROM parks LEFT JOIN sensors ON sensors.park_id = parks.id WHERE parks.id = ?"
}

func InsertSensor() string {
	return "INSERT INTO sensors(serial_number) VALUES(?)"
}

func InsertPark() string {
	return "INSERT INTO parks(location, name) VALUES(?, ?)"
}

func UpdateParkStatus() string {
	return "UPDATE parks SET is_observed = ? WHERE id = ?"
}

func UpdateSensorPark() string {
	return "UPDATE sensors SET park_id = ? WHERE id = ?"
}

func UpdateSensorParkAndStatus() string {
	return "UPDATE sensors SET park_id = ?, is_active = ? WHERE id = ?"
}

func DeleteSensor() string {
	return "DELETE FROM sensors WHERE id = ?"
}

func DeletePark() string {
	return "DELETE FROM parks WHERE id = ?"
}
