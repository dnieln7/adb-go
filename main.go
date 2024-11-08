package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	addresses := scanDeviceAdresses()

	scanner := bufio.NewScanner(os.Stdin)
	input := ""

	for input != "e" {
		displayMenu(addresses)

		scanner.Scan()
		input = scanner.Text()

		if input == "s" {
			addresses = scanDeviceAdresses()
			continue
		}

		option, err := strconv.Atoi(input)

		if err != nil {
			continue
		}

		realOption := option - 1

		if realOption < 0 || realOption >= len(addresses) {
			log.Println("Select one of the list")
			continue
		}

		address := addresses[realOption]

		if len(address) > 0 {
			connectToDevice(address)
			break
		} else {
			log.Println("Select one of the list")
		}
	}

	log.Print("adb go exit\n\n")
}

func scanDeviceAdresses() []string {
	log.Print("Scanning for devices...\n\n")

	outputBytes, err := exec.Command("adb", "mdns", "services").Output()

	if err != nil {
		log.Fatal("Failed to scan services", err)
	}

	outputString := string(outputBytes)
	outputSlice := strings.Split(outputString, "\n")
	addresses := getAddressesFromOutputSlice(outputSlice)

	return addresses
}

func connectToDevice(address string) {
	log.Printf("Connecting to %s...", address)

	outputBytes, err := exec.Command("adb", "connect", address).Output()

	if err != nil {
		log.Fatal("Failed to connect", err)
	}

	outputString := string(outputBytes)

	log.Print(outputString)

	if strings.HasPrefix(outputString, "failed to connect") {
		ip := strings.Split(address, ":")[0]
		pairDevice(ip)
	}
}

func pairDevice(ip string) {
	log.Printf("Starting pairing request to: %s...", ip)
	log.Println("Enter pairing port and code in the format: port-code")
	log.Println("Enter 'c' to cancel pairing")

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	input := scanner.Text()

	if input == "c" {
		return
	}

	inputArray := strings.Split(input, "-")

	if len(inputArray) != 2 {
		log.Printf("%s is not in the format: port-code", input)
	}

	port := inputArray[0]
	code := inputArray[1]

	pairingAddress := ip + ":" + port

	log.Printf("pairing to: %s... with code: %s", pairingAddress, code)

	outputBytes, err := exec.Command("adb", "pair", pairingAddress, code).Output()

	if err != nil {
		log.Println("Failed to pair", err)
	}

	outputString := string(outputBytes)

	log.Print(outputString)
}

func displayMenu(addresses []string) {
	fmt.Print("Select a device\n\n")
	fmt.Println("e - Exit")
	fmt.Println("s - Scan again")

	for index, address := range addresses {
		fmt.Printf("%d - %s\n", index+1, address)
	}

	fmt.Println("")
}

func substring(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func getAddressesFromOutputSlice(output []string) []string {
	regex, err := regexp.Compile(`\d{3}\..+:\d+`)

	if err != nil {
		log.Fatal("Could not create regex", err)
	}

	addresses := []string{}

	for _, line := range output {
		if len(line) > 0 {
			ipIndexes := regex.FindStringIndex(line)

			if len(ipIndexes) >= 2 {
				ip := substring(line, ipIndexes[0], ipIndexes[1]-ipIndexes[0])
				addresses = append(addresses, ip)
			}
		}
	}

	return addresses
}
