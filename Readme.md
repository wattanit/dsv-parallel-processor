# DSV Parallel Processor

## Spec file

DSV Parallel Processor takes input files and query specification via a spec file (conventionally named "spec.toml").

Example spec.toml

```text
[input]
# list all input file in a list
file_path = [
    "test_data/test_data2.txt"
]

# or specify just the input directory
directory = "test_data"

separator = "|"

[output]
# name the output file
outputFile = "output.tsv"
separator = "\t"

# each filter condition is listed below
[filter]

# specify column number to filter (0th-index)
column = 1

# specify column type, currently support number, string, datetime
columnType = string

# specify accept value as a list
values = [
    "OPTION1",
    "OPTION2"
]

# or from a file, one line per one value
value_list = "selected_value.txt"

# or as a range (for number and datetime)
condition = "<"
value = 1979-05-27T07:32:00
```
