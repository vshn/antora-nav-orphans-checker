package main

import (
	"testing"
)

func TestShouldFail(t *testing.T) {
	errors := check("fixture")

	if len(errors) == 1 {
		error := errors[0]
		if error != "third.adoc" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
