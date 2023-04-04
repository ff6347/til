package main

import (
	"bytes"
	"os"
	"path/filepath"
	//"strings"
	"testing"
	//"time"
	"io"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_saveToCsvFile(t *testing.T) {
	// Create a temporary CSV file
	fileName := "test_output_save.csv"
	dataPath, _ := xdg.DataFile(filepath.Join("til", fileName))

	// Content to be written
	data := [][]string{
		{"2022-01-01 00:00:00", "Test entry 1"},
		{"2022-01-02 00:00:00", "Test entry 2"},
	}

	// Call the function to test it
	err := saveToCsvFile(dataPath, data)
	require.NoError(t, err)

	// Read the content from the temporary CSV file
	var buffer bytes.Buffer
	f, err := os.Open(dataPath)
	require.NoError(t, err)
	_, err = buffer.ReadFrom(f)
	require.NoError(t, err)
	_ = f.Close()

	//////
	// Updated expected content of the file (removed quotes)
	var expected bytes.Buffer
	expected.WriteString("2022-01-01 00:00:00,Test entry 1\n")
	expected.WriteString("2022-01-02 00:00:00,Test entry 2\n")
	//////

	expectedStr := expected.String()
	outputStr := buffer.String()
	assert.Equal(t, expectedStr, outputStr) // assert.Equal will check if the expected output is equal to the received output

	// Cleanup the temporary file
	err = os.Remove(dataPath)
	require.NoError(t, err) // require.NoError will stop the test execution if there's an error
}

func Test_listCsvContents(t *testing.T) {
	// Create a temporary CSV file
	fileName := "test_output.csv"
	dataPath, _ := xdg.DataFile(filepath.Join("til", fileName))

	// Fill the temporary file with content
	data := [][]string{
		{"2022-01-01 00:00:00", "Test entry 1"},
		{"2022-01-02 00:00:00", "Test entry 2"},
	}
	err := saveToCsvFile(dataPath, data)
	require.NoError(t, err) // require.NoError will stop the test execution if there's an error

	// Call the function to test it
	var buffer bytes.Buffer
	buffer.WriteString("Timestamp\t\t\tTIL\n")
	buffer.WriteString("2022-01-01 00:00:00\t")
	buffer.WriteString("Test entry 1\n")
	buffer.WriteString("2022-01-02 00:00:00\t")
	buffer.WriteString("Test entry 2\n")

	expected := buffer.String()

	// Capture the output of the function
	output := captureOutput(func() {
		err = listCsvContents(dataPath)
		require.NoError(t, err) // require.NoError will stop the test execution if there's an error
	})

	assert.Equal(t, expected, output) // assert.Equal will check if the expected output is equal to the received output
}

func captureOutput(f func()) string {
	// Set up a pipe to capture the output
	r, w, _ := os.Pipe()
	origOut := os.Stdout
	defer func() {
		// Restore the original stdout after capturing
		os.Stdout = origOut
	}()
	os.Stdout = w

	// Invoke the provided function
	f()

	// Close the write end of the pipe
	_ = w.Close()

	// Buffer to capture the output
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}

func TestMain(m *testing.M) {
	retCode := m.Run()

	// Cleanup
	fileNames := []string{
		"test_output.csv",
		"test_output_save.csv",
	}

	for _, fileName := range fileNames {
		dataPath, _ := xdg.DataFile(filepath.Join("til", fileName))
		os.Remove(dataPath)
	}

	os.Exit(retCode)
}
