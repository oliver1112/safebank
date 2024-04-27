package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Define a Room struct as the innermost level of nesting
type Room struct {
	Type   string
	Number int
}

// Define a Detail struct that includes Room as a nested struct
type Detail struct {
	Building   string
	Apartment  int
	RoomDetail Room
}

// Define an Address struct that includes Detail as a nested struct
type Address struct {
	Street, City, State, Country string
	ZipCode                      int
	ResidenceDetail              Detail
}

// Define a User struct that includes Address as a nested struct
type User struct {
	FirstName, LastName string
	Age                 int
	Address             Address
}

func TestEncrypt(t *testing.T) {
	// Create an instance of User with nested Address, Detail, and Room
	user := User{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Address: Address{
			Street:  "1234 Maple Street",
			City:    "Somewhere",
			State:   "CA",
			Country: "USA",
			ZipCode: 90210,
			ResidenceDetail: Detail{
				Building:  "Building A",
				Apartment: 101,
				RoomDetail: Room{
					Type:   "Living Room",
					Number: 1,
				},
			},
		},
	}

	// Convert the User instance into JSON format
	userJSON, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(userJSON))
}
