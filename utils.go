package main

import (
	"errors"
	"fmt"
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
