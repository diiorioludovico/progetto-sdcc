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
	id            int
	is_active     bool
	park_id       sql.NullInt64
	serial_number string
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
			addSensor(db)
		case "4":
			addPark(db)
		case "5":
			removeSensor(db)
		case "6":
			removePark(db)
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
}

func showSensors(db *sql.DB) {
	fmt.Println("\nshowSensors")
	var sensors []Sensor

	rows, err := db.Query("SELECT * FROM sensors")
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	for rows.Next() {
		var sen Sensor
		if err := rows.Scan(&sen.id, &sen.is_active, &sen.park_id, &sen.serial_number); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}

		sensors = append(sensors, sen)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("ERROR: rows error: ", err)
	} else {
		sensorTable(sensors)
	}
}

func showParks(db *sql.DB) {
	fmt.Println("\nshowParks")
	var parks []Park

	rows, err := db.Query("SELECT * FROM parks")
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	for rows.Next() {
		var park Park
		if err := rows.Scan(&park.id, &park.location, &park.name, &park.is_observed); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}

		parks = append(parks, park)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("ERROR: rows error: ", err)
	} else {
		parkTable(parks)
	}
}

func addSensor(db *sql.DB) {
	fmt.Println("\naddSensor")
}

func addPark(db *sql.DB) {
	fmt.Println("\naddPark")
}

func removeSensor(db *sql.DB) {
	fmt.Println("\nremoveSensor")
}

func removePark(db *sql.DB) {
	fmt.Println("\nremovePark")
}

func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPress ENTER to return to the menu.")
	reader.ReadString('\n')
}

func sensorTable(sensors []Sensor) {
	// Tabwriter per stampa tabellare allineata
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "\nid\tis Active\tpark id\tserial number")
	fmt.Fprintln(writer, "--\t----\t-----\t-----")

	for _, s := range sensors {
		fmt.Fprintf(writer, "%d\t%t\t%d\t%s\n", s.id, s.is_active, s.park_id.Int64, s.serial_number)
	}

	writer.Flush()
}

func parkTable(parks []Park) {
	// Tabwriter per stampa tabellare allineata
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "\nid\tilocation\tname\tis observed")
	fmt.Fprintln(writer, "--\t----\t-----\t-----")

	for _, p := range parks {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%t\n", p.id, p.location, p.name, p.is_observed)
	}

	writer.Flush()
}
