package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davealexis/seesv"
	"github.com/eiannone/keyboard"
	"github.com/olekukonko/tablewriter"
)

func main() {
	skipLinesArg := flag.Int("s", 0, "Number of lines at the top of the file to skip.")
	flag.Parse()

	skipLines := *skipLinesArg

	if flag.NArg() == 0 {
		fmt.Println("Please specify a file to open")
		fmt.Println()
		return
	}

	filepath := flag.Args()[0]
	fileInfo, err := os.Stat(filepath)

	if os.IsNotExist(err) || fileInfo.IsDir() {
		fmt.Println("Could not open the file ", filepath)
		fmt.Println()
		return
	}

	var csvFile seesv.DelimitedFile

	err = csvFile.Open(filepath, skipLines, true)

	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	DisplayRows := 20

	fmt.Println()
	log.Println(csvFile.RowCount, " rows")

	currentRow := int64(0)
	displayRows(&csvFile, currentRow, DisplayRows, 0)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	stop := false
	columnShift := 0
	columnCount := len(csvFile.Headers) - 1

	for stop == false {
		fmt.Print("Press ESC to quit >> ")
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyPgup:
			if currentRow < int64(DisplayRows) {
				currentRow = 0
			} else {
				currentRow -= int64(DisplayRows)
			}
		case keyboard.KeyPgdn:
			currentRow += int64(DisplayRows)

			if currentRow >= csvFile.RowCount {
				currentRow = csvFile.RowCount - int64(DisplayRows)
			}
		case keyboard.KeyArrowRight:
			columnShift += 2
			if columnShift+15 > columnCount {
				columnShift = columnCount - 15
			}
		case keyboard.KeyArrowLeft:
			columnShift -= 2
			if columnShift < 0 {
				columnShift = 0
			}
		case keyboard.KeyEsc:
			fmt.Print("\033[H\033[2J")
			stop = true
		case 0:
			switch char {
			case 'g':
				currentRow = 0
			case 'G':
				currentRow = csvFile.RowCount - int64(DisplayRows)
			case '/':
				fmt.Print("/")
				var input string
				fmt.Scan(&input)
				input = strings.TrimSuffix(input, "\r\n")
				var line int64
				line, err := strconv.ParseInt(strings.TrimSuffix(input, "\r\n"), 10, 64)
				fmt.Println("Go to:", line)

				if err == nil {
					currentRow = int64(line)
					if currentRow >= csvFile.RowCount {
						currentRow = csvFile.RowCount - int64(DisplayRows)
					}
				}
			}
		}

		if stop == false {
			displayRows(&csvFile, currentRow, DisplayRows, columnShift)
		}
	}
}

func displayRows(csv *seesv.DelimitedFile, start int64, rows int, colShift int) {
	rowNum := start
	end := start + int64(rows)
	columnCount := 15
	columnScrollIndicator := ""

	if colShift == 0 {
		columnScrollIndicator = ""
	} else {
		columnScrollIndicator = "<"
	}

	if colShift+columnCount < len(csv.Headers)-1 {
		columnScrollIndicator += " >"
	}
	// w, h, err := terminal.GetSize(int(os.Stdout.Fd()))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(h, w)

	fmt.Print("\033[H\033[2J")
	fmt.Println(start, "->", end, "of", csv.RowCount)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	rowHeaders := append([]string{columnScrollIndicator, "Row"}, csv.Headers[colShift:columnCount+colShift]...)

	table.SetHeader(rowHeaders)

	for v := range csv.Rows(start, -1) {
		table.Append(append([]string{"", strconv.FormatInt(rowNum, 10)}, v[colShift:columnCount+colShift]...))

		rowNum++

		if rowNum >= end {
			break
		}
	}

	table.Render()
}
