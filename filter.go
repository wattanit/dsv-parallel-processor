package main

import (
	"strings"
)

func filter(line string, spec Spec) bool {

	inputSeparator := spec.Input.Separator

	cells := strings.Split(line, inputSeparator)

	for _, filter := range spec.Filters {
		//fmt.Println(filter)

		if len(cells) >= filter.Column {
			columnValue := cells[filter.Column]
			//fmt.Println(columnValue)
			//fmt.Println(isin(columnValue, filter.Values))
			switch filter.ColumnType {
			case "datetime":
				continue
			case "number":
				continue
			// string type is default
			default:
				{
					if !isin(columnValue, filter.Values) {
						return false
					}
				}
			}
		}
	}
	return true

}

func isin(a string, list []string) bool {
	for _, item := range list {
		if item == a {
			return true
		}
	}
	return false
}
