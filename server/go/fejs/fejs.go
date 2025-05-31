package fejs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// API Response hello
type JSMessage struct {
	Parks []Park `json:"parks"`
}

type Park struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Location    string         `json:"location"`
	Temperature string         `json:"temperature"`
	Humidity    string         `json:"humidity"`
	Brightness  string         `json:"brightness"`
	AirQuality  string         `json:"airquality"`
	Timestamp   string         `json:"timestamp"`
	OldData     []ParksOldData `json:"olddata"`
}

type ParksOldData struct {
	Date string `json:"date"`
	Icon string `json:"icon"`
	Min  string `json:"min"`
	Max  string `json:"max"`
}

// API Response getData
type MetricResp struct {
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Hour  string `json:"hour"`
	Value string `json:"value"`
}

var dbsql *sql.DB

func getMetricData(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	w.Header().Set("Content-Type", "application/json")
	mx := getData(r)
	json.NewEncoder(w).Encode(mx)
}

func getData(r *http.Request) MetricResp {
	park_id := r.URL.Query().Get("id")
	metric := r.URL.Query().Get("metric")

	//lavorazione stringa metrica
	metric = strings.ToLower(metric)
	metric = strings.ReplaceAll(metric, " ", "_")

	fmt.Println(park_id)
	fmt.Println(metric)

	query := fmt.Sprintf(`
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

	rows, err := dbsql.Query(query, park_id)
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	defer rows.Close()

	var data MetricResp

	for rows.Next() {
		var metric Metric
		if err := rows.Scan(&metric.Hour, &metric.Value); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}

		data.Metrics = append(data.Metrics, metric)
	}
	return data
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	w.Header().Set("Content-Type", "application/json")
	mx := makeResponse()
	json.NewEncoder(w).Encode(mx)
}

func makeResponse() JSMessage {
	//query per recuperare la misuraziona più recente per ogni parco
	query := `
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

	rows, err := dbsql.Query(query)
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	defer rows.Close()

	var count int
	var data JSMessage

	for rows.Next() {
		var park Park
		if err := rows.Scan(&park.Id, &park.Name, &park.Location, &park.Temperature, &park.Humidity, &park.Brightness, &park.AirQuality, &park.Timestamp); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}

		data.Parks = append(data.Parks, park)
		data.Parks[count].OldData = getOldData(data.Parks[count].Id)
		count += 1
	}

	if count == 0 {
		fmt.Println("INFO: no observed parks")
		return data
	}

	return data
}

func getOldData(id string) []ParksOldData {
	query := `
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

	rows, err := dbsql.Query(query, id)
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	defer rows.Close()

	var old_data []ParksOldData

	for rows.Next() {
		var data ParksOldData
		if err := rows.Scan(&data.Date, &data.Max, &data.Min); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}
		data.Icon = "1"
		old_data = append(old_data, data)
	}

	return old_data
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func StartFrontendSetup(db *sql.DB) {
	dbsql = db
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/getData", getMetricData)
	http.ListenAndServe(":8080", nil)
}
