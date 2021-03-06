# DSV Parallel Processor

## Spec file

DSV Parallel Processor takes input files and query specification via a spec file (conventionally named "spec.toml").

Example spec.toml

```text
[[input]]
# filePaths = [                          # list all input file in a list
#     "test_data/test_data2.txt"
# ]

directory = "test_data"                  # or specify just the input directory
separator = "|"

[[output]]
outputFile = "output.tsv"                # name the output file
separator = "\t"

# each filter condition is listed below

# example of string filter
[[filters]]
column = 16                              # specify column to filter (0th-index)
columnType = "string"                    # available type are string, number, datetime
values = [                               # list accepted value as a list
    "OPTION2",
    "OPTION1"
]

# valueFile = "filter.txt"               # or read value from a file, one line per one value

[[filters]]
column = 1
columnType = "string"
valueFile = "account_list.txt"

# Example of number filter
# [[filters]]
# column = 6                             # specify column to filter (0th-index)
# columnType = "number"                  # available type are string, number, datetime
# condition = "<"                        # available condition "<", "<=", ">", ">=", "=="
# value = "250"                          # condition value to check

# Example of datetime filter
# [[filters]]
# column = 3                             # specify column to filter (0th-index)
# columnType = "datetime"                # available type are string, number, datetime
# condition = "<"                        # available condition "<", "<=", ">", ">=", "=="
# datetimeFormat = "02/01/2006"          # specify datetime format using Golang's notation.
# value = "01/01/2015"                   # condition value to check

# Golang datetime format can be found at https://programming.guide/go/format-parse-string-time-date-example.html
```

## Run with Docker

The provided ```filter_csv.sh``` script will run the program as a Docker container. 

```
./filter_csv.sh spec.toml data_dir
```

Please note the following details:
1. All input and output files are mounted to the container at /data directory. Therefore, all *data_dir* paths in the spec file must be replaced by ```/data/```

