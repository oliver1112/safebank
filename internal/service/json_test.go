package service

import (
	"encoding/json"
	"fmt"
	"log"
)

// Define a struct
type Person struct {
	Name    string
	Age     int
	Address string
}

func main() {
	// Create an instance of the struct
	person := Person{
		Name:    "John Doe",
		Age:     30,
		Address: "123 Elm St",
	}

	// Convert struct to JSON string
	jsonData, err := json.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}

	// Print the JSON string
	fmt.Println(string(jsonData))
}
