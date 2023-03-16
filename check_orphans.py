#!/usr/bin/env python3

# This script alerts of orphan files in the modules/ROOT/pages directory
# that are not referenced in modules/ROOT/nav.adoc.

from os import path, walk
import re

# Opens the input file referenced by its filename, and uses the
# regular expression to filter all included files. It then
# iterates over the list and appends an error if a file in
# the directory is not explicitly referenced in the input.
def check(input, regex):
    with open(input, "r") as file:
        contents = file.read()
    matches = re.findall(regex, contents)
    for f in adoc_files:
        if f not in matches:
            errors.append('File "{entry}" not in {input}'.format(entry=f, input=input))


# Global variables
pages_dir = "modules/ROOT/pages/"

# Replace the 'pages_dir' variable above form the input with an empty string and return
def remove(input):
    return input.replace(pages_dir, "")

# Get all Asciidoc files recursively
# https://stackoverflow.com/a/19309964
all_files = [path.join(dp, f) for dp, dn, fn in walk(path.expanduser(pages_dir)) for f in fn]

# Filter the path name from filenames
filtered = list(map(remove, all_files))

# Only care for certain files
adoc_files = [
    f
    for f in filtered
    if not f == "search.adoc"
    and not f == ".vale.ini"
    and not f == "index.adoc"
    and not re.search("^inc\_", f)
]
errors = []

# Perform checks on the nav.adoc file
nav_adoc = path.join(path.abspath(path.dirname(__file__)), "modules/ROOT/nav.adoc")
check(nav_adoc, r"xref:(.+)\[")

# Exit with error if some file is orphan
if len(errors) > 0:
    for e in errors:
        print(e)
    exit(1)

# All is well, bye bye
print("No orphan files in modules/ROOT/nav.adoc")
exit(0)
