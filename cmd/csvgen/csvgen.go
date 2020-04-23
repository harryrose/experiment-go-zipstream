package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage:")
		fmt.Printf("\t%s <filename> <noCols> <noRows>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	cols, err := strconv.Atoi(os.Args[2])
	if err != nil  || cols <= 0{
		fmt.Println("number of columns must be a number > 0")
		os.Exit(2)
	}

	rows, err := strconv.Atoi(os.Args[3])
	if err != nil || rows <= 0 {
		fmt.Println("number of rows must be a number > 0")
		os.Exit(3)
	}

	if _, err := os.Stat(filename); err == nil {
		// file exists, chicken out
		fmt.Printf("file %s exists. will not overwrite\n", filename)
		os.Exit(4)
	}

	if err := writeFile(filename, cols, rows); err != nil {
		fmt.Printf("error writing to file, %s: %v\n",filename, err)
		os.Exit(5)
	}

	fmt.Printf("%s written\n",filename)
}

type WriteStringer interface{
	WriteString(string) (int, error)
}

func writeRow(w WriteStringer, prefix string, cols int, rowPrefix string) error {
	_, err := w.WriteString(fmt.Sprintf("%s%s%d",rowPrefix, prefix,0))
	if err != nil {
		return err
	}

	for col := 1; col < cols; col ++ {
		_, err :=w.WriteString(fmt.Sprintf(",%s%d",prefix,col))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filename string, cols, rows int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// emit the column headers
	if err := writeRow(file, "col", cols, ""); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	for row := 0; row < rows; row ++ {
		if err := writeRow(file, fmt.Sprintf("r%dc",row), cols, "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	return nil
}
