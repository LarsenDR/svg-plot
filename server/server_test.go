package server

import (
	"fmt"
	"testing"
)

func TestClientData(t *testing.T) {

	dt, err := GetClientData("testdata.json")
	if err != nil {
		fmt.Printf("In Test:  %#v\n", err)
	}

	fmt.Printf("In Test: %v\n", dt)

}
