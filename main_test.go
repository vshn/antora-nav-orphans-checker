package main

import (
	"testing"
)

func TestShouldFail(t *testing.T) {
	modules := []string{"ROOT", "TestModule"}
	errors := checkAntora("fixture", modules)

	if len(errors) == 2 {
		error := errors[0]
		if error != "third.adoc" {
			t.Fail()
		}
		error = errors[1]
		if error != "three.adoc" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
