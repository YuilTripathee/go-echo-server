package main

import (
	"database/sql"
	"fmt"
)

var (
	// database information
	dbhostsip  = "127.0.0.1:3306" // IP address of database server
	dbusername = "root"           // username of the database user
	dbpassword = ""               // password of the database user
	dbname     = "database_name"  // name of the database
	dbcharset  = "utf8"           // database character set

	// list of relations in database
	sampleTable = "sample" // table name for airplane
)

// function to generate database connection string
func getConnectionString() string {
	return dbusername + ":" + dbpassword + "@tcp(" + dbhostsip + ")/" + dbname + "?charset=" + dbcharset
}

// function to return JSON array from a MySQL database
func getJSONFromDB(sqlString string) ([]map[string]interface{}, error) {
	sqlConnString := getConnectionString()
	db, err := sql.Open("mysql", sqlConnString) //
	if err != nil {
		return nil, err // return error if present
	}

	defer db.Close()                 // close database connection if error occurs
	rows, err := db.Query(sqlString) // typically returns rows querying the database
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns() // returns the column name and returns errors if rows are closed
	if err != nil {
		return nil, err
	}
	count := len(columns) // number of columns from where rows are rendered.
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	// JSON format builder
	for rows.Next() { // for each row
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i] // find value pointers
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns { // for each data in row buld a json data (with key and value4)
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v // assign data as entry["key"] = value
		}
		tableData = append(tableData, entry) // append the data object map[string]interface{} to array []map[string]interface{}
	}
	// for marshalling (i.e. serializing)
	// jsonData, err := json.Marshal(tableData)
	// if err != nil {
	// 	return jsonData, err
	// }

	return tableData, nil
}

// function to get data from MySQL database ==> gets you ready response data
// sample buildout query: SELECT * FROM table_name WHERE id = id_requested;
func getDataDBbyIndex(table string, index string, id string) (int, map[string]interface{}) {
	status := 200
	response := make(map[string]interface{})
	sqlQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s';", table, index, id)
	dbData, err := getJSONFromDB(sqlQuery)
	if err != nil {
		status = 500
		response = statMsgData[3]
	} else if len(dbData) == 0 {
		status = 500
		response = statMsgData[4]
	} else {
		response = statMsgData[2]
		response["data"] = dbData
	}
	return status, response
}
