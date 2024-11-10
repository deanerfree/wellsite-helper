package main

import (
	"bufio"
	// "errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deanerfree/wellsite-helper/utils"
	// "strings"
	// "text/scanner"
)

func scanLasFile(lasFile string) {
	// Check the file extension
	if filepath.Ext(lasFile) != ".las" {
		fmt.Println("Invalid file extension")
		return
	}

	file, err := os.Open(lasFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	target := ""

	for scanner.Scan() {
		// fmt.Printf("Scanning target: %s\n", target)
		line := scanner.Text()

		if strings.Contains(line, "~V") {
			// fmt.Println("Populate version information")
			// populate the version information struct skip to the next line
			target = "VersionInformation"
			continue
		}

		if strings.Contains(line, "~W") {
			// fmt.Println("Populate version information")
			target = "WellInformation"
			continue
		}

		if strings.Contains(line, "~C") {
			// fmt.Println("Populate curve information")
			target = "CurveInformation"
			continue
		}

		if strings.Contains(line, "~P") {
			// fmt.Println("Populate parameter information")
			target = "ParameterInformation"
			// populate the parameter information struct skip to the next line
			continue
		}

		if strings.Contains(line, "~O") {
			// fmt.Println("Populate other information")
			target = "OtherInformation"
			// populate the other information struct skip to the next line
			continue
		}

		if strings.Contains(line, "~A") {
			// fmt.Println("Populate well information")
			target = "DepthData"
			// populate the well information struct skip to the next line
			continue
		}

		utils.parseData(line, target)

	}

	// fmt.Printf("LasData: %s\n", lasData)
}

func main() {
	// setup server

	scanLasFile("mockData/test2.las")
	fmt.Printf("LasData: %s\n", lasData)
}
