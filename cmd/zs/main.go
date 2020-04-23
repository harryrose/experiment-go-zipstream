package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/harryrose/experiment-go-zipstream/handler"
	"os"
)

var handlerFactory = handler.FileHandlerFactory{
	&handler.CSVHandlerConstructor{},
}

func main() {
	ctx := context.Background()
	if l := len(os.Args); l < 2 {
		fmt.Println("Usage:")
		fmt.Printf("\t%s <zipFile>\n", os.Args[0])
		os.Exit(1)
	}


	for _, zipFilename := range os.Args[1:] {
		err := handleZipFile(ctx, zipFilename)
		if err != nil {
			fmt.Printf("Error handling zip file %s: %v", zipFilename, err)
		}
	}
	fmt.Println("Done")
}

func handleZipFile(ctx context.Context, zipFilename string) error {
	// this would probably be coming from an HTTP endpoint,
	// but the concept is similar, provided we get a ReaderAt interface
	zipFile, err := os.Open(zipFilename)
	defer zipFile.Close()
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	fileInfo, err := os.Stat(zipFilename)
	if err != nil {
		return fmt.Errorf("error statting file: %w", err)
	}

	zipReader, err := zip.NewReader(zipFile, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("error opening zip reader: %w", err)
	}

	for _, fileWithinZip := range zipReader.File{
		err := handleFileInsideZip(ctx, fileWithinZip)
		if err != nil {
			return fmt.Errorf("error handling %s: %w", fileWithinZip.Name, err)
		}
	}

	return nil
}

func handleFileInsideZip(ctx context.Context, file *zip.File) error {
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer reader.Close()

	h := handlerFactory.Construct(ctx, file.Name)
	if h == nil {
		return fmt.Errorf("no handler for file")
	}

	for item := range h.Handle(ctx, reader) {
		marsh, _ := json.Marshal(item)
		fmt.Println(string(marsh))
	}
	return nil
}