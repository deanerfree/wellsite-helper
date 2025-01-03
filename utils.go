package main

import (
	"bufio"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type EmptyStruct struct{}

type StandardInformation struct {
	Fields map[string]Data
}

type Data struct {
	Name        string
	Value       interface{}
	Description string
	Unit        string
}

type VersionInformation struct {
	Fields map[string]Data
}

type WellInformation struct {
	Fields map[string]Data
}

type CurveInformation struct {
	Fields     map[string]Data
	CurveOrder []string // Order of curve names for mapping data lines in ~A section
}

type Depth map[string]float64

type AsciiLog struct {
	Fields []map[string]string // List of data rows, each row is a map of field values
}

type LasData struct {
	VersionInformation   VersionInformation
	WellInformation      WellInformation
	CurveInformation     CurveInformation
	OtherInformation     EmptyStruct
	ParameterInformation EmptyStruct
	WellData             AsciiLog
}

func handleStandardInformation(line string) (Data, error) {
	if line == "" {
		return Data{}, errors.New("empty line")
	}

	index := strings.Index(line, ".")
	descriptionIndex := strings.LastIndex(line, ":")

	if index == -1 || descriptionIndex == -1 {
		return Data{}, errors.New("invalid line format")
	}

	key := strings.TrimSpace(line[:index])
	unit := strings.TrimSpace(line[index+1 : index+2])
	data := strings.TrimSpace(line[index+2 : descriptionIndex])
	description := strings.TrimSpace(line[descriptionIndex+1:])

	newEntry := Data{
		Name:        key,
		Value:       data,
		Description: description,
		Unit:        unit,
	}

	// fmt.Printf("Added entry for %s - Key: %s, Data: %v, Description: %s\n", key, newEntry.name, newEntry.value, newEntry.description)

	return newEntry, nil
}

func handleCurveInformation(line string) Data {
	substrings := []string{"~Version", "~WELL", "#", "~Curve", "~Parameter", "~Other", "~A"}
	for _, substring := range substrings {
		if strings.Contains(line, substring) {
			return Data{}
		}
	}

	index := strings.Index(line, ".")
	descriptionIndex := strings.LastIndex(line, ":")

	if index == -1 {
		return Data{}
	}

	// if no description index is found, set it to the end of the line
	newEntry := Data{
		Name:        strings.TrimSpace(line[:index]),
		Value:       strings.TrimSpace(line[index+2:]),
		Description: strings.TrimSpace(line[descriptionIndex+1:]),
		Unit:        strings.TrimSpace(line[index+1 : index+2]),
	}

	return newEntry
}

func handleData(line string, parsedData *struct {
	VersionInformation   VersionInformation
	WellInformation      WellInformation
	CurveInformation     CurveInformation
	OtherInformation     EmptyStruct
	ParameterInformation EmptyStruct
	WellData             AsciiLog
},
) map[string]string {
	// Split the line by whitespace to get each data value
	values := strings.Fields(line)
	if len(values) != len(parsedData.CurveInformation.CurveOrder) {
		fmt.Println("Data line does not match curve information")
		return nil
	}

	// Map each value to the curve name and add it as a new entry in wellData.Fields
	dataEntry := make(map[string]string)
	for i, value := range values {
		curveName := parsedData.CurveInformation.CurveOrder[i]
		dataEntry[curveName] = value
	}

	return dataEntry
}

func ParseData(line string, target string, parsedData *struct {
	VersionInformation   VersionInformation
	WellInformation      WellInformation
	CurveInformation     CurveInformation
	OtherInformation     EmptyStruct
	ParameterInformation EmptyStruct
	WellData             AsciiLog
},
) ([]string, error) {
	substrings := []string{"~Version", "~WELL", "#", "~Curve", "~Parameter", "~Other", "~A"}
	for _, substring := range substrings {
		if strings.Contains(line, substring) {
			return nil, errors.New("invalid line")
		}
	}

	// Populate the appropriate field based on the target
	switch target {
	case ("VersionInformation"):
		newEntry, errors := handleStandardInformation(line)
		if errors != nil {
			fmt.Println(errors)
			return nil, errors
		}
		fmt.Printf("New Version information: %+v\n", newEntry)

		parsedData.VersionInformation.Fields[newEntry.Name] = newEntry
		// fmt.Printf("Added entry for %s - Key: %s, Data: %v, Description: %s\n", target, newEntry.name, newEntry.value, newEntry.description)
	case "WellInformation":
		newEntry, errors := handleStandardInformation(line)
		if errors != nil {
			fmt.Println(errors)
			return nil, errors
		}

		fmt.Printf("New Well information: %+v\n", newEntry)

		parsedData.WellInformation.Fields[newEntry.Name] = newEntry
		// fmt.Printf("Added entry for %s - Key: %s, Data: %v, Description: %s\n", target, newEntry.name, newEntry.value, newEntry.description)
	case "CurveInformation":
		newEntry := handleCurveInformation(line)
		parsedData.CurveInformation.Fields[newEntry.Name] = newEntry
		parsedData.CurveInformation.CurveOrder = append(parsedData.CurveInformation.CurveOrder, newEntry.Name)

	case "DepthData":
		data := handleData(line, parsedData)
		parsedData.WellData.Fields = append(parsedData.WellData.Fields, data)
		// fmt.Printf("Added data entry: %+v\n", data)
	default:
		fmt.Printf("Unknown target: %s\n", target)
	}

	return nil, nil
}

func ScanLasFile(lasFile string) (struct {
	VersionInformation   VersionInformation
	WellInformation      WellInformation
	CurveInformation     CurveInformation
	OtherInformation     EmptyStruct
	ParameterInformation EmptyStruct
	WellData             AsciiLog
}, error,
) {
	parsedData := struct {
		VersionInformation   VersionInformation
		WellInformation      WellInformation
		CurveInformation     CurveInformation
		OtherInformation     EmptyStruct
		ParameterInformation EmptyStruct
		WellData             AsciiLog
	}{
		VersionInformation: VersionInformation{
			Fields: make(map[string]Data),
		},
		WellInformation: WellInformation{
			Fields: make(map[string]Data),
		},
		CurveInformation: CurveInformation{
			Fields:     make(map[string]Data),
			CurveOrder: []string{},
		},
		WellData: AsciiLog{
			Fields: []map[string]string{},
		},
	}
	// Check the file extension
	if filepath.Ext(lasFile) != ".las" {
		fmt.Println("Invalid file extension")

		return parsedData, errors.New("invalid file extension")
	}

	file, err := os.Open(lasFile)
	if err != nil {
		fmt.Println(err)
		return parsedData, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	target := ""

	// setup empty struct

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

		ParseData(line, target, &parsedData)

	}

	// fmt.Printf("LasData: %s\n", LasData)
	return parsedData, nil
}

func ScanSurveyFile(surveyFile *multipart.FileHeader) (SurveyData, error) {
	fmt.Print("Scanning survey file\n")
	var headers []string
	var surveyData SurveyData

	longestLine := 0
	row := 0
	// Check the file extension

	parsedData := SurveyData{}

	file, err := surveyFile.Open()
	if err != nil {
		fmt.Println(err)
		return parsedData, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// columns := []string{}

	for scanner.Scan() {
		row++
		line := scanner.Text()
		values := strings.Fields(line)

		if strings.Contains(strings.ToLower(line), "measured") || strings.Contains(strings.ToLower(line), "md") {
			headers = values
			continue
		}

		if line == "EOF" {
			break
		}

		if len(values) > longestLine {
			longestLine = len(values)
		}

	}
	fmt.Println(headers)
	fmt.Println(row)

	return parsedData, nil
}
