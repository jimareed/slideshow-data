package main

import (
	"errors"
	"fmt"
)

/* Data type */
type Data struct {
	Id          int
	ResourceId  string
	Name        string
	Description string
	Permissions string
}

var nextId = 4

var data = []Data{
	Data{Id: 1, ResourceId: "default", Name: "Slideshow", Description: "Overview", Permissions: ""},
	Data{Id: 2, ResourceId: "instructions", Name: "Instructions", Description: "Steps to use", Permissions: ""},
	Data{Id: 3, ResourceId: "emotional-intelligence", Name: "Emotional Intelligence", Description: "Sample slideshow", Permissions: ""},
}

func readData(userId string) []Data {

	filteredData := []Data{}

	for _, d := range data {
		d.Permissions = "read write"

		filteredData = append(filteredData, d)
	}

	return filteredData
}

func updateData(userId string, name string, description string) error {

	index := 0

	for _, d := range data {
		if name == d.Name {
			data[index].Description = description
			return nil
		}
		index++
	}

	return errors.New("data not found")
}

func duplicateData(userId string, resourceId string) (Data, error) {

	newData := Data{}

	for _, s := range data {
		if s.ResourceId == resourceId {
			newData = Data{Id: nextId, ResourceId: "2345", Name: "copy of " + s.Name, Description: s.Description, Permissions: ""}
			data = append(data, newData)
			nextId++
			return newData, nil
		}
	}

	return newData, fmt.Errorf("Duplicate error: invalid id %s", resourceId)
}
