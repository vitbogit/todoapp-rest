package filters

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
)

const (
	SortTypeAsc  = "asc"
	SortTypeDesc = "desc"
)

func IsValidSortType(candidate string) bool {
	return candidate == SortTypeAsc || candidate == SortTypeDesc
}

type Filtering struct {
	SortType   string   `json:"sort_type"`
	SortColumn string   `json:"sort_column"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	Filters    []Filter `json:"filters"`
}

type Filter struct {
	FieldName string      `json:"field_name"`
	Equals    interface{} `json:"equals"`
}

func (f Filtering) Filter(query string, argsInQuery int, allowedColumns []string, choosenColumns ...string) (string, []interface{}, error) {
	if len(choosenColumns) == 0 {
		return query, nil, nil
	}

	for _, column := range choosenColumns {
		if !slices.Contains(allowedColumns, column) {
			return "", nil, fmt.Errorf("choosen column not allowed")
		}
	}
	for _, filter := range f.Filters {
		if !slices.Contains(allowedColumns, filter.FieldName) {
			return "", nil, fmt.Errorf("choosen field in filters not allowed")
		}
	}

	doSort := IsValidSortType(f.SortType)

	if doSort {
		if len(choosenColumns) == 0 {
			return "", nil, errors.New("no columns specified")
		}
		if len(f.SortColumn) == 0 {
			f.SortColumn = choosenColumns[0]
		}
	}

	tableName := "unfiltered"
	newquery := "select "
	for i, column := range choosenColumns {
		newquery += fmt.Sprintf("%s.%s", tableName, column)
		if i != len(choosenColumns)-1 {
			newquery += ", "
		}
	}
	newquery += fmt.Sprintf(" from (%s) %s", query, tableName)
	narg := argsInQuery + 1
	args := make([]interface{}, 0, len(f.Filters))
	if len(f.Filters) > 0 {
		newquery += " where "
		for i, filter := range f.Filters {
			newquery += fmt.Sprintf("%s.%s = $%d", tableName, filter.FieldName, narg)
			//filterValue, err := ConvertInterfaceToString(filter.Equals)
			// if err != nil {
			// 	return "", nil, err
			// }
			args = append(args, filter.Equals)
			if i != len(f.Filters)-1 {
				newquery += " and "
			}
			narg += 1
		}
	}

	if doSort {
		newquery += fmt.Sprintf(" order by %s %s", f.SortColumn, f.SortType)
	}

	if f.Limit > 0 {
		newquery += fmt.Sprintf(" limit %d", f.Limit)
	}

	if f.Offset > 0 {
		newquery += fmt.Sprintf(" offset %d", f.Offset)
	}

	return newquery, args, nil

}

func ConvertInterfaceToString(interf interface{}) (string, error) {
	if val, ok := interf.(string); ok {
		return val, nil
	}
	if val, ok := interf.(int64); ok {
		str := strconv.FormatInt(val, 10)
		return str, nil
	}
	if val, ok := interf.(int); ok {
		str := strconv.Itoa(val)
		return str, nil
	}
	if val, ok := interf.(bool); ok {
		str := strconv.FormatBool(val)
		return str, nil
	}
	return "", errors.New("can`t convert interface value to string")
}
