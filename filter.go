package main

import (
	"strconv"
	"strings"
)

func filter(line string, spec Spec) bool {

	inputSeparator := spec.Input.Separator

	cells := strings.Split(line, inputSeparator)

	for _, filter := range spec.Filters {
		//fmt.Println(filter)

		if len(cells) < filter.Column {
			return false
		}

		columnValue := cells[filter.Column]
		//fmt.Println(columnValue)
		//fmt.Println(isin(columnValue, filter.Values))
		switch filter.ColumnType {
		case "datetime":
			continue
		case "number":
			{
				colValue, err := strconv.ParseFloat(columnValue, 64)
				if err != nil {
					return false
				}
				if !compareNumberField(colValue, filter) {
					return false
				}
			}

		// string type is default
		default:
			{
				if !isin(columnValue, filter.Values) {
					return false
				}
			}
		}
	}
	return true
}

func compareNumberField(columeValue float64, filter SpecFilter) bool {
	filterValue, err := strconv.ParseFloat(filter.Value, 64)
	check(err)

	switch filter.Condition {
	case "<":
		return columeValue < filterValue
	case "<=":
		return columeValue <= filterValue
	case ">":
		return columeValue > filterValue
	case ">=":
		return columeValue >= filterValue
	case "==":
		return columeValue == filterValue
	default:
		return columeValue == filterValue
	}
}
