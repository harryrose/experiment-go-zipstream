package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CSVHandlerConstructor provides a method for building a CSV Handler
type CSVHandlerConstructor struct {

}

// CanHandle returns true if the file is a csv file
func (c *CSVHandlerConstructor) CanHandle(_ context.Context, filename string) bool {
	return strings.HasSuffix(filename, ".csv")
}

// Construct returns a CSV Handler
func (c *CSVHandlerConstructor) Construct(_ context.Context) Handler {
	return &csvHandler{}
}

// csvHandler returns a function for handling CSV files
type csvHandler struct{

}

// Handle streams data from the reader one CSV row at a time and emits data to a channel.
// The channel is closed either on an error (including EOF) or when the context is closed
func (c *csvHandler) Handle(ctx context.Context, reader io.Reader) <-chan Item {
	out := make(chan Item)

	go func() {
		defer close(out)
		reader := csv.NewReader(reader)
		// first row is headings
		headings, err := reader.Read()
		if err != nil {
			// in production you'd do some proper error handling here...
			fmt.Printf("error occurred while reading from CSV: %v\n", err)
			return
		}

		for {
			row, err := reader.Read()
			if err != nil {
				fmt.Printf("error occurred while reading from CSV: %v\n", err)
				return
			}

			select {
				case <-ctx.Done():
					// context closed
					return

				case out <- toMap(headings, row):
					// do nothing
			}
		}
	}()

	return out
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func toMap(keys, values []string) map[string]string {
	out := make(map[string]string)

	for i := 0; i < min(len(keys),len(values)); i ++ {
		out[keys[i]] = values[i]
	}
	return out
}
