package menu

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
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

	rows, err := db.Query("SELECT * FROM sensors")
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	var count int

	for rows.Next() {
		var sen Sensor
		if err := rows.Scan(&sen.Id, &sen.Is_active, &sen.Park_id, &sen.Serial_number); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}

		sensors = append(sensors, sen)
		count += 1
	}

	if err := rows.Err(); err != nil {
		fmt.Println("ERROR: rows error: ", err)
	} else if count > 0 {
		sensorTable(sensors)
	} else {
		fmt.Println("There are not sensor records")
	}
}

func showParks(db *sql.DB) {
	//fmt.Println("\nshowParks")
	var parks []Park

	rows, err := db.Query("SELECT * FROM parks")
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	var count int

	for rows.Next() {
		var park Park
		if err := rows.Scan(&park.id, &park.location, &park.name, &park.is_observed); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}

		parks = append(parks, park)
		count += 1
	}

	if err := rows.Err(); err != nil {
		fmt.Println("ERROR: rows error: ", err)
	} else if count > 0 {
		parkTable(parks)
	} else {
		fmt.Println("There are not park records")
	}
}

func addSensor(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\naddSensor")
	fmt.Print("\nInsert sensor serial number: ")
	// Legge l’input da tastiera
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	_, err := db.Exec("INSERT INTO sensors(serial_number) VALUES(?)", input)
	if err != nil {
		fmt.Println("ERROR: error in inserting new sensor: ", err)
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

	_, err := db.Exec("INSERT INTO parks(location, name) VALUES(?, ?)", location, name)
	if err != nil {
		fmt.Println("ERROR: error in inserting new park: ", err)
	}
}

func removeSensor(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\nremoveSensor")
	fmt.Print("\nInsert sensor id to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	_, err := db.Exec("DELETE FROM sensors WHERE id = ?", id)
	if err != nil {
		fmt.Println("ERROR: error in inserting new park: ", err)
	}
}

func removePark(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\nremovePark")
	fmt.Print("\nInsert park id to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	_, err := db.Exec("DELETE FROM parks WHERE id = ?", id)
	if err != nil {
		fmt.Println("ERROR: error in inserting new park: ", err)
	}
}

func associateSensor(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\nremovePark")
	fmt.Print("\nInsert park id: ")
	park_id, _ := reader.ReadString('\n')
	park_id = strings.TrimSpace(park_id)

	fmt.Print("\nInsert sensor id: ")
	sensor_id, _ := reader.ReadString('\n')
	sensor_id = strings.TrimSpace(sensor_id)

	//modifica del valore is_observed del parco per indicare che è stato posto un sensore e che tra poco sarà attivato
	_, err := db.Exec("UPDATE parks SET is_observed = ? WHERE id = ?", true, park_id)
	if err != nil {
		fmt.Println("ERROR: park update error: ", err)
	}

	//modifica del valore park_id del sensore per indicare il parco a cui è stato assegnato
	_, err = db.Exec("UPDATE sensors SET park_id = ? WHERE id = ?", park_id, sensor_id)
	if err != nil {
		fmt.Println("ERROR: sensor update error: ", err)
	}
}

func deassociateSensor(db *sql.DB, reader *bufio.Reader) {
	//fmt.Println("\nremovePark")
	fmt.Print("\nInsert sensor id: ")
	sensor_id, _ := reader.ReadString('\n')
	sensor_id = strings.TrimSpace(sensor_id)

	var park_id int
	err := db.QueryRow("SELECT park_id FROM sensors WHERE id = ?", sensor_id).Scan(&park_id)
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	//modifica del valore is_observed del parco per indicare che non è più osservato
	_, err = db.Exec("UPDATE parks SET is_observed = ? WHERE id = ?", false, park_id)
	if err != nil {
		fmt.Println("ERROR: park update error: ", err)
	}

	//modifica del valore park_id del sensore per indicare il parco a cui è stato assegnato
	_, err = db.Exec("UPDATE sensors SET park_id = ?, is_active = ? WHERE id = ?", nil, false, sensor_id)
	if err != nil {
		fmt.Println("ERROR: sensor update error: ", err)
	}

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
