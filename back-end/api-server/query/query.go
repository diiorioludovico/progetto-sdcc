package query

import (
	"fmt"
)

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
    		measures AS m
		WHERE
    		m.timestamp >= DATE_SUB(CURRENT_DATE(), INTERVAL 7 DAY) AND m.timestamp < CURRENT_DATE() AND park_id = ?
		GROUP BY
    		m.park_id,
    		DATE(m.timestamp)
		ORDER BY
    		observation_date;`
}
