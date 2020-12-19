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

	t.Log("Alice:")

	d, err := duplicateData("alice@example.com", "default")
	if err == nil {
		t.Log("Should be able to duplicate the default Slideshow.", checkMark)
	} else {
		t.Fatal("Should be able to duplicate the default Slideshow.", xMark, err)
	}

	_ = deleteData("alice@example.com", d.Id)

}

func TestDeleteData(t *testing.T) {

	t.Log("Alice:")

	userId := "alice@example.com"

	d, err := duplicateData(userId, "default")
	if err != nil {
		t.Fatal("Should be able to duplicate the default Slideshow.", xMark, err)
	}
	_, err = duplicateData(userId, "default")
	if err != nil {
		t.Fatal("Should be able to duplicate the default Slideshow.", xMark, err)
	}

	data := readData(userId)

	if len(data) == 5 {
		t.Log("Should have access to five data items.", checkMark)
	} else {
		t.Fatal("Should have access to five data items.", xMark, len(data))
	}

	err = deleteData(userId, d.Id)
	if err == nil {
		t.Log("Should be able to delete a data item.", checkMark)
	} else {
		t.Fatal("Should be able to delete a data item.", xMark)
	}

	data = readData(userId)

	for _, d2 := range data {
		if d2.Id == d.Id {
			t.Fatal("Should not find deleted id in list of data items.", xMark)
		}
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
