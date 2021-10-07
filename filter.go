package main

import (
	"strconv"
	"strings"
	"time"
)

func filter(line string, spec Spec) bool {

	inputSeparator := spec.Input.Separator

	cells := strings.Split(line, inputSeparator)

	for _, filter := range spec.Filters {

		if len(cells) < filter.Column {
			return false
		}

		columnValue := cells[filter.Column]

		switch filter.ColumnType {
		case "datetime":
			if !compareDatetimeField(columnValue, filter) {
				return false
			}
		case "number":
			if !compareNumberField(columnValue, filter) {
				return false
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

func compareNumberField(columeValue string, filter SpecFilter) bool {
	colValue, err := strconv.ParseFloat(columeValue, 64)
	if err != nil {
		return false
	}
	filterValue, err := strconv.ParseFloat(filter.Value, 64)
	check(err)

	switch filter.Condition {
	case "<":
		return colValue < filterValue
	case "<=":
		return colValue <= filterValue
	case ">":
		return colValue > filterValue
	case ">=":
		return colValue >= filterValue
	case "==":
		return colValue == filterValue
	default:
		return colValue == filterValue
	}
}

func compareDatetimeField(columnValue string, filter SpecFilter) bool {
	colValue, err := time.Parse(filter.DatetimeFormat, columnValue)
	if err != nil {
		return false
	}

	filterValue, err := time.Parse(filter.DatetimeFormat, filter.Value)
	check(err)

	switch filter.Condition {
	case "<":
		return colValue.Before(filterValue)
	case "<=":
		return colValue.Before(filterValue) || colValue.Equal(filterValue)
	case ">":
		return colValue.After(filterValue)
	case ">=":
		return colValue.After(filterValue) || colValue.Equal(filterValue)
	case "==":
		return colValue.Equal(filterValue)
	default:
		return colValue.Equal(filterValue)
	}

}
