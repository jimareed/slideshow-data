package data

import (
	"errors"
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

/* DataItem type */
type DataItem struct {
	Id          int
	Name        string
	Description string
	ResourceId  string
	Permissions string
}

/* Data type */
type Data struct {
	enforcer *casbin.Enforcer
}

var nextId = 4

var dataItems = []DataItem{
	DataItem{Id: 1, Name: "Slideshow", Description: "Overview", ResourceId: "default", Permissions: ""},
	DataItem{Id: 2, Name: "Instructions", Description: "Steps to use", ResourceId: "instructions", Permissions: ""},
	DataItem{Id: 3, Name: "Emotional Intelligence", Description: "Sample slideshow", ResourceId: "emotional-intelligence", Permissions: ""},
}

func Init(modelFile string, policyFile string) Data {
	d := Data{}

	e, err := casbin.NewEnforcer(modelFile, policyFile)
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
		return d
	}

	d.enforcer = e

	return d
}

func (data Data) ReadData(userEmail string) []DataItem {

	filteredData := []DataItem{}

	for _, d := range dataItems {
		id := fmt.Sprintf("%d", d.Id)

		d.Permissions = ""

		hasRead, err := data.enforcer.Enforce(userEmail, id, "read:data")
		if err != nil {
			log.Fatalf("Enforce error: %v", err)
		}
		hasWrite, err := data.enforcer.Enforce(userEmail, id, "write:data")
		if err != nil {
			log.Fatalf("Enforce error: %v", err)
		}
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
	id := fmt.Sprintf("%d", newData.Id)

	_, err := data.enforcer.AddPolicy(userId, id, "write:data")
	if err != nil {
		log.Fatalf("error adding policy: %v", err)
		return newData, err
	}

	dataItems = append(dataItems, newData)
	nextId++

	return newData, nil
}

func (data Data) UpdateData(userId string, id int, name string, description string) error {
	index := 0

	sid := fmt.Sprintf("%d", id)

	for _, d := range dataItems {
		result, err := data.enforcer.Enforce(userId, sid, "write:data")
		if err != nil {
			log.Fatalf("Enforce error: %v", err)
			return err
		}
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

	sid := fmt.Sprintf("%d", id)

	for i, d := range dataItems {
		result, err := data.enforcer.Enforce(userId, sid, "write:data")
		if err != nil {
			log.Fatalf("Enforce error: %v", err)
			return err
		}
		if result {
			if d.Id == id {
				if i != len(dataItems)-1 {
					dataItems[i] = dataItems[len(dataItems)-1]
				}
				dataItems = dataItems[:len(dataItems)-1]
				_, err := data.enforcer.RemovePolicy(userId, sid, "write:data")
				if err != nil {
					log.Fatalf("error removing policy: %v", err)
					return err
				}
				return nil
			}
		}
	}

	return errors.New("insufficient privileges or data not found")
}
