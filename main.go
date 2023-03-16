package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Needs as argument a path to an Antora project
	if len(os.Args) < 2 {
		fmt.Println("Please provide a path to an Antora project")
		os.Exit(1)
	}
	path := os.Args[1]

	// Remove the trailing slash
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	// We assume that the project follows a standard Antora layout
	startPath := path + "/modules/ROOT/pages"
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

	// Verify that all filtered files appear in the nav.adoc file at least once
	navPath := path + "/modules/ROOT/nav.adoc"
	regex := `xref:(.+)\[`
	errors := check(navPath, filtered, regex)

	// If there are errors, print the list of orphan files
	if len(errors) > 0 {
		for _, v := range errors {
			fmt.Println(v)
		}
		os.Exit(1)
	}

	// No errors, all good!
	fmt.Println("No orphan files in modules/ROOT/nav.adoc")
	os.Exit(0)
}

// Checks the file nav.adoc for the existence of at least one
// reference to each one of the files passed in the second parameter,
// using the regular expresssion as parameter.
func check(navPath string, files []string, regex string) []string {
	var errors []string
	re := regexp.MustCompile(regex)
	contents, err := os.ReadFile(navPath)
	if err == nil {
		stringContents := string(contents)
		matches := re.FindAllString(stringContents, -1)
		for _, file := range files {
			match := fmt.Sprintf("xref:%s[", file)
			if !stringInSlice(match, matches) {
				errors = append(errors, fmt.Sprintf("File '%s' not in 'nav.adoc'", file))
			}
		}
	} else {
		errors = append(errors, "Cannot read file "+navPath)
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
