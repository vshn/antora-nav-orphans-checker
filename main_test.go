package main

import (
	"testing"
)

func TestShouldFailNavInROOT(t *testing.T) {
	errors := check("fixture", "ROOT", "/modules/ROOT/nav.adoc")

	if len(errors) == 1 {
		error := errors[0]
		if error != "third.adoc" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestShouldFailNavInAnotherModule(t *testing.T) {
	errors := check("fixture", "AnotherModule", "/modules/AnotherModule/nav.adoc")

	assert.Len(t, errors, 1)
	assert.Contains(t, errors, "one.adoc")
}

func TestShouldFailDocument(t *testing.T) {
	errors := check("fixture", "ROOT", "/document.adoc")

	if len(errors) == 1 {
		error := errors[0]
		if error != "second.adoc" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
