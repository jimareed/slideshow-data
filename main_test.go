package main

import (
	"testing"
)

const checkMark = "\u2713"
const xMark = "\u2717"

func TestGet(t *testing.T) {

	t.Log("A user:")

	s, err := Get("default")
	if err == nil {
		t.Log("Should be able to get the default Slideshow.", checkMark)
	} else {
		t.Fatal("Should be able to get the default Slideshow.", xMark, err)
	}

	if s.Name == "" {
		t.Fatal("Which should have a valid name.", xMark, s.Name)
	} else {
		t.Log("Which should have a valid name.", checkMark)
	}

	s1, err := Get("instructions")

	if s1.Name == s.Name {
		t.Fatal("Which should have a different name from the Instructions slideshow.", xMark, s.Name, s1.Name)
	} else {
		t.Log("Which should have a different name from the Instructions slideshow.", checkMark)
	}

}

func TestDuplicate(t *testing.T) {

	t.Log("A user:")

	_, err := Duplicate("default")
	if err == nil {
		t.Log("Should be able to duplicate the default Slideshow.", checkMark)
	} else {
		t.Fatal("Should be able to duplicate the default Slideshow.", xMark, err)
	}
}
