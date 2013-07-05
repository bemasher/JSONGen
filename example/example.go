package main

import (
	"encoding/json"
	"fmt"
	"jsongen"
	"log"
	"os"
)

func main() {
	testFile, err := os.Open("test.json")
	if err != nil {
		log.Fatal("Error opening filter file:", err)
	}
	defer testFile.Close()

	jsonDecoder := json.NewDecoder(testFile)
	var data interface{}
	err = jsonDecoder.Decode(&data)
	if err != nil {
		log.Fatal("Error decoding filters:", err)
	}

	t := jsongen.Parse("Test", data)

	fmt.Println(t.Format())

	fmt.Println()
	fmt.Printf("%#v\n", t.Fields["non-homogeneous"])
}
