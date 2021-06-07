# Go CSV Viewer

A commaand line app for viewing large CSV files.

Uses the SeeSV Go library for reading and navigating CSV files.

## Build

```shell
go mod tidy
go build
```

## Run

Open a CSV file.

```shell
./go-csv-viewer /path/to/myfile.csv
```

Open a CSV file that has 1 extra metadata line at the top of the file that should be skipped, with the column headers on line 2.

```shell
./go-csv-viewer -s 1 /path/to/myfile.csv
```
