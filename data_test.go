package main

import (
	"testing"
)

const checkMark = "\u2713"
const xMark = "\u2717"

func TestReadData(t *testing.T) {

	t.Log("Alice:")

	data := readData("alice@example.com")

	if len(data) == 3 {
		t.Log("Should have access to three data items.", checkMark)
	} else {
		t.Fatal("Should have access to three data items.", xMark)
	}
}

func TestDuplicateData(t *testing.T) {

	t.Log("A user:")

	_, err := duplicateData("", "default")
	if err == nil {
		t.Log("Should be able to duplicate the default Slideshow.", checkMark)
	} else {
		t.Fatal("Should be able to duplicate the default Slideshow.", xMark, err)
	}
}

func TestUpdateData(t *testing.T) {

	t.Log("Alice:")

	data := readData("alice@example.com")

	for _, d := range data {
		if d.Name == "Slideshow" {
			if d.Description == "Overview" {
				t.Log("Should be able to lookup the description for data1 ", checkMark)
			} else {
				t.Fatal("Should have access to two data items.", xMark, d.Description)
			}
		}
	}

	err := updateData("alice@example.com", "Slideshow", "Slideshow Updated")

	if err == nil {
		t.Log("Should be able to update the description for data1 ", checkMark)
	} else {
		t.Fatal("Should be able to update the description for data1.", xMark, err)
	}

	data = readData("alice@example.com")

	for _, d := range data {
		if d.Name == "Slideshow" {
			if d.Description == "Slideshow Updated" {
				t.Log("Should be able to get the updated description for data1 ", checkMark)
			} else {
				t.Fatal("Should be able to get the updated description for data1.", xMark, d.Description)
			}
		}
	}

}
