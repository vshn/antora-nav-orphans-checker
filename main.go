package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Needs as argument a path to an Antora project
	var antoraPath string
	flag.StringVar(&antoraPath, "antoraPath", "/docs", "Path to Antora document sources")

	// Module to check. If none passed, it is assumed to be ROOT
	var module string
	flag.StringVar(&module, "module", "ROOT", "Module to analyze")

	// File to parse and analyze; if none is passed, it is assumed to
	// be a standard Antora navigation file at the ROOT module.
	var filename string
	flag.StringVar(&filename, "filename", "/modules/ROOT/nav.adoc", "File to analyze")

	flag.Parse()

	// Remove the trailing slash
	if strings.HasSuffix(antoraPath, "/") {
		antoraPath = strings.TrimSuffix(antoraPath, "/")
	}

	// Main block
	errors := check(antoraPath, module, filename)

	// If there are errors, print the list of orphan files
	if len(errors) > 0 {
		for _, file := range errors {
			fmt.Fprintf(os.Stdout, "File '%s' not in '%s'\n", file, filename))
		}
		os.Exit(1)
	}

	// No errors, all good!
	fmt.Println(fmt.Sprintf("No orphan files in '%s'", filename))
	os.Exit(0)
}

// Takes a path and assumes it's a valid Antora project. It lists all pages
// contained in the module and then walks the file system to verify that
// all relevant files are referenced in the file.
func check(antoraPath string, module string, filename string) []string {
	// We assume that the project follows a standard Antora layout
	startPath := antoraPath + "/modules/" + module + "/pages"
	allFiles, err := listAllFiles(startPath)
	if err != nil {
		fmt.Println("Cannot list files in provided path " + startPath)
		os.Exit(1)
	}

	// Remove the initial part of the path from all filenames in the list
	simplified := mapArray(allFiles, func(input string) string {
		return strings.Replace(input, startPath+"/", "", -1)
	})

	// Filter out some stuff we don't show on the navigation anyway
	filtered := filterArray(simplified, func(input string) bool {
		avoid := []string{"search.adoc", "index.adoc", "sitesearch.adoc", "changelog_from_commits.adoc", "changelog_legacy.adoc"}
		return !stringInSlice(input, avoid)
	})

	// Verify that all filtered files appear in the file at least once
	fullPath := filepath.Join(antoraPath, filename)
	regex := `xref:(.+)\[`

	// If the file being checked is not nav.adoc, then it is assumed to be a standard
	// Asciidoc file, where files are referenced using `include::...[]` instead of `xref:...[]`
	if !strings.HasSuffix(fullPath, "nav.adoc") {
		regex = `include::(.+)\[`
	}

	// Gather all errors
	errors := walk(fullPath, filtered, regex)
	return errors
}

// Checks the index file for the existence of at least one
// reference to each one of the files passed in the second parameter,
// using the regular expresssion as parameter.
func walk(fullPath string, files []string, regex string) []string {
	var errors []string
	re := regexp.MustCompile(regex)
	contents, err := os.ReadFile(fullPath)
	if err == nil {
		stringContents := string(contents)
		matches := re.FindAllString(stringContents, -1)
		errors = filterArray(files, func(file string) bool {
			return !substringInSlice(file, matches)
		})
	} else {
		errors = append(errors, "Cannot read file "+fullPath)
	}
	return errors
}

// Collect filenames recursively at the path provided as argument
func listAllFiles(startPath string) ([]string, error) {
	var files []string
	_, err := os.Stat(startPath)
	if err != nil {
		fmt.Println("Path does not exist " + startPath)
		return files, err
	}
	err = filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".adoc") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// Filter array using the function passed as argument
func filterArray(array []string, f func(string) bool) []string {
	filtered := make([]string, 0)
	for _, s := range array {
		if f(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// Apply the transformation function to all items of the array
func mapArray(array []string, f func(string) string) []string {
	mapped := make([]string, len(array))
	for i, e := range array {
		mapped[i] = f(e)
	}
	return mapped
}

// Checks whether a string appears in an array
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Checks whether a string appears in an array complete or as a substring
func substringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a || strings.Index(b, a) != -1 {
			return true
		}
	}
	return false
}
