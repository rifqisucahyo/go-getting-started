package models

import (
	"github.com/jinzhu/gorm"
)

func getDataRows(db *gorm.DB, qry string, args ...interface{}) ([]map[string]interface{}, int, error) {

	rows, err := db.Raw(qry, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var numrow int = 0
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	parseData := make([]map[string]interface{}, 0, 0)
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		tmp_struct := make(map[string]interface{})

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			tmp_struct[col] = v
		}
		parseData = append(parseData, tmp_struct)
		numrow++
	}

	rows.Close()
	return parseData, numrow, nil
}
