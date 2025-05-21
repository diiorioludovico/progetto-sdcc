package menu

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func ShowMenu() {
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
			showSensors()
		case "2":
			showParks()
		case "3":
			addSensor()
		case "4":
			addPark()
		case "5":
			removeSensor()
		case "6":
			removePark()
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

func showSensors() {
	fmt.Println("\nshowSensors")
}

func showParks() {
	fmt.Println("\nshowParks")
}

func addSensor() {
	fmt.Println("\naddSensor")
}

func addPark() {
	fmt.Println("\naddPark")
}

func removeSensor() {
	fmt.Println("\nremoveSensor")
}

func removePark() {
	fmt.Println("\nremovePark")
}

func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPress ENTER to return to rhe menu.")
	reader.ReadString('\n')
}
