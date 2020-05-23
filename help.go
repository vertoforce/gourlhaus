package gourlhaus

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gocarina/gocsv"
)

// downloadCSVZip takes a link, makes the request, decodes the zip contents (zipFile=true), and fills in the structs with the csv contents
func downloadCSV(ctx context.Context, url string, zipFile bool, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if !zipFile {
		return readCSV(ctx, resp.Body, out)
	}

	// Download all data
	allData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Decode zip file
	readerAt := bytes.NewReader(allData)
	zipDecoder, err := zip.NewReader(readerAt, int64(len(allData)))
	if err != nil {
		return fmt.Errorf("error opening zip file: %s", err)
	}
	for _, file := range zipDecoder.File {
		f, err := file.Open()
		if err != nil {
			return err
		}
		return readCSV(ctx, f, out)
	}

	return fmt.Errorf("no files in zip file")
}

// readCSV Reads a csv into a provided slice
func readCSV(ctx context.Context, reader io.Reader, out interface{}) error {
	// URLHaus Specific:
	// We need to read the first few comments lines
	// And the # in front of the header row
	bufReader := bufio.NewReader(reader)
	for i := 0; i < 8; i++ {
		bufReader.ReadLine()
	}
	bufReader.ReadByte()
	bufReader.ReadByte()

	if err := gocsv.Unmarshal(bufReader, out); err != nil {
		return err
	}
	return nil
}
