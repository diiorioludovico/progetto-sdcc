package menu

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"progetto/server/go/logger"
	qr "progetto/server/go/query"
	"runtime"
	"strings"
	"text/tabwriter"
)

type Sensor struct {
	Id            int
	Is_active     bool
	Park_id       sql.NullInt64
	Serial_number string
}

type Park struct {
	id          int
	location    string
	name        string
	is_observed bool
}

func ShowMenu(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printMenu()

		fmt.Print("Select an option: ")
		// Legge l’input da tastiera
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Gestisce le scelte
		switch input {
		case "1":
			showSensors(db)
		case "2":
			showParks(db)
		case "3":
			addSensor(db, reader)
		case "4":
			addPark(db, reader)
		case "5":
			removeSensor(db, reader)
		case "6":
			removePark(db, reader)
		case "7":
			associateSensor(db, reader)
		case "8":
			deassociateSensor(db, reader)
		default:
			fmt.Println("\nInvalid option. Try Again.")
		}
		waitForEnter(reader)

		fmt.Println()
	}
}

// Pulisce lo schermo in modo compatibile con Windows/Linux/macOS
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default: // Unix-like systems
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printMenu() {
	fmt.Println("====== Menu ======")
	fmt.Println("1) Show sensors")
	fmt.Println("2) Show parks")
	fmt.Println("3) Add sensor")
	fmt.Println("4) Add park")
	fmt.Println("5) Remove sensor")
	fmt.Println("6) Remove park")
	fmt.Println("7) Associate sensor to park")
	fmt.Println("8) Disassociate sensor from park")

}

func showSensors(db *sql.DB) {
	//fmt.Println("\nshowSensors")
	var sensors []Sensor

	rows, err := db.Query(qr.ShowSensors())
	if err != nil {
		logger.Error.Println("Query error: ", err)

	}

	defer rows.Close()

	var count int

	for rows.Next() {
		var sen Sensor
		if err := rows.Scan(&sen.Id, &sen.Is_active, &sen.Park_id, &sen.Serial_number); err != nil {
			logger.Error.Println("Scan error: ", err)
		}

		sensors = append(sensors, sen)
		count += 1
	}

	if err := rows.Err(); err != nil {
		logger.Error.Println("Rows error: ", err)
	} else if count > 0 {
		sensorTable(sensors)
	} else {
		logger.Info.Println("There are not sensor records")
	}
}

func showParks(db *sql.DB) {
	//fmt.Println("\nshowParks")
	var parks []Park

	rows, err := db.Query(qr.ShowParks())
	if err != nil {
		logger.Error.Println("Query error: ", err)
	}

	defer rows.Close()

	var count int

	for rows.Next() {
		var park Park
		if err := rows.Scan(&park.id, &park.location, &park.name, &park.is_observed); err != nil {
			logger.Error.Println("Scan error: ", err)
		}

		parks = append(parks, park)
		count += 1
	}

	if err := rows.Err(); err != nil {
		logger.Error.Println("Rows error: ", err)
	} else if count > 0 {
		parkTable(parks)
	} else {
		logger.Info.Println("There are not park records")
	}
}

func addSensor(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\naddSensor")
	fmt.Print("\nInsert sensor serial number: ")
	// Legge l’input da tastiera
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	_, err := db.Exec(qr.InsertSensor(), input)
	if err != nil {
		logger.Error.Println("error in inserting new sensor: ", err)
	}
}

func addPark(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\naddPark")
	fmt.Print("\nInsert park location: ")
	location, _ := reader.ReadString('\n')
	location = strings.TrimSpace(location)

	fmt.Print("\nInsert park name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	_, err := db.Exec(qr.InsertPark(), location, name)
	if err != nil {
		logger.Error.Println("error in inserting new parl: ", err)
	}
}

func removeSensor(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("\nInsert sensor id to delete: ")
	sensor_id, _ := reader.ReadString('\n')
	sensor_id = strings.TrimSpace(sensor_id)

	// Avvia una transazione
	tx, err := db.Begin()
	if err != nil {
		logger.Error.Println("Transaction begin failed:", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			logger.Error.Println("Recovered from panic, transaction rolled back:", r)
		}
	}()

	// 1. Controlla se il sensore è associato a un parco
	var park_id sql.NullInt64
	err = tx.QueryRow(qr.GetSensorParkid(), sensor_id).Scan(&park_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Query park_id failed, transaction rolled back:", err)
		return
	}

	// 2. Se associato, aggiorna is_observed del parco
	if park_id.Valid {
		_, err = tx.Exec(qr.UpdateParkStatus(), false, park_id.Int64)
		if err != nil {
			_ = tx.Rollback()
			logger.Error.Println("Park update error, transaction rolled back:", err)
			return
		}
	}

	// 3. Elimina il sensore
	_, err = tx.Exec(qr.DeleteSensor(), sensor_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Sensor delete error, transaction rolled back:", err)
		return
	}

	// 4. Commit
	err = tx.Commit()
	if err != nil {
		logger.Error.Println("Transaction commit failed:", err)
		return
	}

	logger.Info.Println("Sensor removed successfully")
}

func removePark(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("\nInsert park id to delete: ")
	park_id, _ := reader.ReadString('\n')
	park_id = strings.TrimSpace(park_id)

	// Avvia una transazione
	tx, err := db.Begin()
	if err != nil {
		logger.Error.Println("Transaction begin failed:", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			logger.Error.Println("Recovered from panic, transaction rolled back:", r)
		}
	}()

	// 1. Verifica se il parco è osservato
	var is_observed bool
	var sensor_id sql.NullInt64
	err = tx.QueryRow(qr.GetParkStatus(), park_id).Scan(&is_observed, &sensor_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Query error, transaction rolled back:", err)
		return
	}

	if is_observed {
		_ = tx.Rollback()
		logger.Error.Printf("\nYou need to deassociate it from the sensor %d to remove it\n", sensor_id.Int64)
		return
	}

	// 2. Cancella il parco
	_, err = tx.Exec(qr.DeletePark(), park_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Error deleting park, transaction rolled back:", err)
		return
	}

	// 3. Commit finale
	err = tx.Commit()
	if err != nil {
		logger.Error.Println("Transaction commit failed:", err)
		return
	}

	logger.Info.Println("Park removed successfully")
}

func associateSensor(db *sql.DB, reader *bufio.Reader) {
	//esecuzione di una transazione
	fmt.Print("\nInsert park id: ")
	park_id, _ := reader.ReadString('\n')
	park_id = strings.TrimSpace(park_id)

	fmt.Print("\nInsert sensor id: ")
	sensor_id, _ := reader.ReadString('\n')
	sensor_id = strings.TrimSpace(sensor_id)

	tx, err := db.Begin()
	if err != nil {
		logger.Error.Println("Transaction error: ", err)
	}

	// Se qualcosa va storto, annulla tutto
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			logger.Error.Println("Recovered from panic, transaction rolled back:", r)
		}
	}()

	// 1. Update sul parco
	_, err = tx.Exec(qr.UpdateParkStatus(), true, park_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Park update error, transaction rolled back:", err)
		return
	}

	// 2. Update sul sensore
	_, err = tx.Exec(qr.UpdateSensorPark(), park_id, sensor_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Sensor update error, transaction rolled back:", err)
		return
	}

	// Tutto ok, conferma la transazione
	err = tx.Commit()
	if err != nil {
		logger.Error.Println("Transaction commit failed:", err)
		return
	}

	logger.Info.Println("Sensor associated to park successfully")
}

func deassociateSensor(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("\nInsert sensor id: ")
	sensor_id, _ := reader.ReadString('\n')
	sensor_id = strings.TrimSpace(sensor_id)

	// Avvia una transazione
	tx, err := db.Begin()
	if err != nil {
		logger.Error.Println("Transaction begin failed:", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			logger.Error.Println("Recovered from panic, transaction rolled back:", r)
		}
	}()

	// 1. Ricava il park_id associato al sensore
	var park_id int
	err = tx.QueryRow(qr.GetSensorParkid(), sensor_id).Scan(&park_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Query park_id failed, transaction rolled back:", err)
		return
	}

	// 2. Setta is_observed a false per il parco
	_, err = tx.Exec(qr.UpdateParkStatus(), false, park_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Park update error, transaction rolled back:", err)
		return
	}

	// 3. Disassocia il sensore e aggiorna il suo stato
	_, err = tx.Exec(qr.UpdateSensorParkAndStatus(), nil, false, sensor_id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error.Println("Sensor update error, transaction rolled back:", err)
		return
	}

	// Commit finale se tutto è andato bene
	err = tx.Commit()
	if err != nil {
		logger.Error.Println("Transaction commit failed:", err)
		return
	}

	logger.Info.Println("Sensor deassociated from park successfully")
}

func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPress ENTER to return to the menu.")
	reader.ReadString('\n')
}

func sensorTable(sensors []Sensor) {
	// Tabwriter per stampa tabellare allineata
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "\nid\tis active\tpark id\tserial number")
	fmt.Fprintln(writer, "--\t----\t-----\t-----")

	for _, s := range sensors {
		fmt.Fprintf(writer, "%d\t%t\t%d\t%s\n", s.Id, s.Is_active, s.Park_id.Int64, s.Serial_number)
	}

	writer.Flush()
}

func parkTable(parks []Park) {
	// Tabwriter per stampa tabellare allineata
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "\nid\tlocation\tname\tis observed")
	fmt.Fprintln(writer, "--\t----\t-----\t-----")

	for _, p := range parks {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%t\n", p.id, p.location, p.name, p.is_observed)
	}

	writer.Flush()
}
