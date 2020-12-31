package main

import (
	"errors"
	"fmt"
)

/* Stub DataItem type */
type DataItem struct {
	Id          int
	Name        string
	Description string
	ResourceId  string
	Permissions string
}

type casbinEnforcer struct {
	i int
}

/* Data type */
type Data struct {
	enforcer *casbinEnforcer
}

var nextId = 4

var dataItems = []DataItem{
	DataItem{Id: 1, Name: "Slideshow", Description: "Overview", ResourceId: "default", Permissions: ""},
	DataItem{Id: 2, Name: "Instructions", Description: "Steps to use", ResourceId: "instructions", Permissions: ""},
	DataItem{Id: 3, Name: "Emotional Intelligence", Description: "Sample slideshow", ResourceId: "emotional-intelligence", Permissions: ""},
}

func Init(modelFile string, policyFile string) Data {
	d := Data{}
	e := casbinEnforcer{0}

	d.enforcer = &e

	return d
}

func (data Data) ReadData(userEmail string) []DataItem {

	filteredData := []DataItem{}

	for _, d := range dataItems {

		d.Permissions = ""

		hasRead := true
		hasWrite := true

		if hasRead {
			d.Permissions = "read"

			if hasWrite {
				d.Permissions += " "
			}
		}
		if hasWrite {
			d.Permissions += "write"
		}

		if hasRead || hasWrite {
			filteredData = append(filteredData, d)
		}
	}

	return filteredData
}

func (data Data) NewData(userId string, resourceId string) (DataItem, error) {

	newData := DataItem{}

	newData.Id = nextId
	newData.Name = fmt.Sprintf("New Slideshow")
	newData.Description = fmt.Sprintf("New Slideshow")
	newData.ResourceId = resourceId

	dataItems = append(dataItems, newData)
	nextId++

	return newData, nil
}

func (data Data) UpdateData(userId string, id int, name string, description string) error {
	index := 0

	for _, d := range dataItems {
		result := true
		if result {
			if id == d.Id {
				dataItems[index].Name = name
				dataItems[index].Description = description
				return nil
			}
		}
		index++
	}

	return errors.New("data not found")
}

func (data Data) DeleteData(userId string, id int) error {

	for i, d := range dataItems {
		result := true
		if result {
			if d.Id == id {
				if i != len(dataItems)-1 {
					dataItems[i] = dataItems[len(dataItems)-1]
				}
				dataItems = dataItems[:len(dataItems)-1]
				return nil
			}
		}
	}

	return errors.New("insufficient privileges or data not found")
}
