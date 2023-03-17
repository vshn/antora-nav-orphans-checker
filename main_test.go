package main

import (
	"fmt"
	"testing"
)

func TestShouldFail(t *testing.T) {
	path := "fixture"
	errors := check(path)

	if len(errors) == 1 {
		error := errors[0]
		fmt.Println(error)
		if error != "third.adoc" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
