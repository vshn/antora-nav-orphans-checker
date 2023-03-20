package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFailNavInROOT(t *testing.T) {
	errors := check("fixture", "ROOT", "/modules/ROOT/nav.adoc")

	assert.Len(t, errors, 1)
	assert.Contains(t, errors, "third.adoc")
}

func TestShouldFailNavInAnotherModule(t *testing.T) {
	errors := check("fixture", "AnotherModule", "/modules/AnotherModule/nav.adoc")

	assert.Len(t, errors, 1)
	assert.Contains(t, errors, "one.adoc")
}

func TestShouldFailDocument(t *testing.T) {
	errors := check("fixture", "ROOT", "/document.adoc")

	assert.Len(t, errors, 1)
	assert.Contains(t, errors, "second.adoc")
}
