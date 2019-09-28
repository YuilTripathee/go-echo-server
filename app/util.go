package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

/*
	Utility modules for supporting rest API
	Note:
		1. Variables with lowercase starting letter are private.
		2. Variables with uppercase starting letter are exported or public.
*/

func init() {
	fmt.Println("Utility function loaded")
	fmt.Println("=======================")
	fmt.Println("Total CPUs in this machine :", runtime.GOMAXPROCS(runtime.NumCPU()))
	fmt.Println()
}

// GetJSONArrayData (Address of file in string form)
// function to get all the data from json file
// exported data is in the form of array of objects
func GetJSONArrayData(FileAddr string) []map[string]interface{} {
	// open the JSON file
	jsonFile, err := os.Open(FileAddr)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

	return result
}

// GetJSONObjectData (Address of file in string form)
// function to get all the data from json file
// exported data is in the form of object map[string]interface{}
func GetJSONObjectData(FileAddr string) map[string]interface{} {
	// open the JSON file
	jsonFile, err := os.Open(FileAddr)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

	return result
}
