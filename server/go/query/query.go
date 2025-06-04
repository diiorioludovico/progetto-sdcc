package query

import (
	"fmt"
)

//server.go

// QUERIES//
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

//menu.go

// QUERIES//
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

// OPERATIONS//
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

//fejs.go

// QUERIES
func GetMetricMeasures(metric string) string {
	return fmt.Sprintf(`
        SELECT
            DATE_FORMAT(timestamp, '%%H:00') AS ora, 
            AVG(%s) AS valore
        FROM 
            measures
        WHERE 
            park_id = ? AND DATE(timestamp) = CURDATE()
        GROUP BY 
            HOUR(timestamp)
        ORDER BY 
            HOUR(timestamp)`, metric)
}

func GetMostRecentParkMeasure() string {
	return `
	WITH latest_measures AS ( 
    	SELECT m.*
    	FROM measures m
    	INNER JOIN (
        	SELECT park_id, MAX(timestamp) AS max_ts
        	FROM measures
        	GROUP BY park_id
    	) latest ON m.park_id = latest.park_id AND m.timestamp = latest.max_ts
	)	
	SELECT 
    	p.id AS park_id,
    	p.name,
    	p.location,
    	lm.temperature,
    	lm.humidity,
    	lm.brightness,
    	lm.air_quality,
    	lm.timestamp
	FROM parks p
	JOIN latest_measures lm ON p.id = lm.park_id
	WHERE p.is_observed = true;`
}

func GetOldData() string {
	return `
		SELECT 
			DATE(m.timestamp) AS observation_date,
    		MAX(m.temperature) AS max_temperature,
    		MIN(m.temperature) AS min_temperature
		FROM
    		MEASURES AS m
		WHERE
    		m.timestamp >= DATE_SUB(CURRENT_DATE(), INTERVAL 7 DAY) AND m.timestamp < CURRENT_DATE() AND park_id = ?
		GROUP BY
    		m.park_id,
    		DATE(m.timestamp)
		ORDER BY
    		observation_date;`
}
