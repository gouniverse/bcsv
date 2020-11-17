package bcsv

import (
	"encoding/base64"
	"errors"
	"log"
	"reflect"
	"strings"
)

// MarshalToFile a struct into BCSV string
func MarshalToFile(filePath string, i interface{}) (bool, error) {
	bcsvString, error := MarshalToString(i)

	if error != nil {
		return false, error
	}

	writeError := filePutContents(filePath, bcsvString, 0)

	if writeError == nil {
		return true, nil
	}

	return false, error
}

// MarshalToString a struct into BCSV string
func MarshalToString(i interface{}) (string, error) {
	data := reflect.ValueOf(i)

	if data.Kind() != reflect.Slice {
		return "", errors.New("only slices can be marshalled")
	}

	if data.IsNil() {
		return "", nil
	}

	if data.Len() < 1 {
		return "", nil
	}

	firstRow := data.Index(0)

	if firstRow.Kind() == reflect.Interface {
		firstRow = firstRow.Elem()
	}

	out := [][]string{}

	var colNames []string
	if firstRow.Type().Kind() == reflect.Struct {
		colNames = reflectColumnsAsStringArrayForBcsvTags(firstRow.Type())
		out = append(out, colNames)
	} else {
		//colNames = firstRow.Interface().([]string)
	}

	// DEBUG: log.Println(colNames)

	for i := 0; i < data.Len(); i++ {
		var values []string
		if firstRow.Type().Kind() == reflect.Struct {
			values = reflectValuesAsStringArrayForBcsvTags(data.Index(i), firstRow.Type())
		} else {
			values = data.Index(i).Interface().([]string)
		}
		out = append(out, values)
	}

	// DEBUG: log.Println(out)

	bcsvRows := []string{}
	for _, row := range out {
		rowValues := []string{}
		for _, cell := range row {
			bcsvValue := base64.StdEncoding.EncodeToString([]byte(cell))
			rowValues = append(rowValues, bcsvValue)
		}
		bcsvRows = append(bcsvRows, strings.Join(rowValues, ","))
	}

	log.Println(bcsvRows)

	bcsvString := strings.Join(bcsvRows, "\n")
	return bcsvString, nil
}

// UnmarshalFromString unmarshals a string
func UnmarshalFromString(bcsvContents string) ([][]string, error) {
	rows := [][]string{}

	bcsvRows := strings.Split(bcsvContents, "\n")

	for _, bcsvRow := range bcsvRows {
		row := []string{}
		bcsvValues := strings.Split(bcsvRow, ",")
		for _, bcsvValue := range bcsvValues {
			decodedString, decodeError := base64.StdEncoding.DecodeString(bcsvValue)
			if decodeError != nil {
				return rows, decodeError
			}
			row = append(row, string(decodedString))
		}

		rows = append(rows, row)
	}

	return rows, nil
}

// UnmarshalFromFile unmarshals a file
func UnmarshalFromFile(filePath string) ([][]string, error) {
	bcsvContents, contentsError := fileGetContents(filePath)

	if contentsError != nil {
		return [][]string{}, contentsError
	}

	rows, unmarshalError := UnmarshalFromString(bcsvContents)

	if unmarshalError != nil {
		return [][]string{}, unmarshalError
	}

	return rows, nil
}
