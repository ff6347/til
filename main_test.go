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
func Test_setupDatabase(t *testing.T) {
	// Create a temporary database file
	fileName := "test_til.db"
	dbPath, _ := xdg.DataFile(filepath.Join("til", fileName))

	// Call the function to be tested
	db, err := setupDatabase(dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Check if the schema was correctly created
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='entries'").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "entries", name)

	// Cleanup
	err = db.Close()
	require.NoError(t, err)
	err = os.Remove(dbPath)
	require.NoError(t, err)
}

func Test_saveToDb(t *testing.T) {
	// Create a temporary database file
	fileName := "test_save_til.db"
	dbPath, _ := xdg.DataFile(filepath.Join("til", fileName))

	db, err := setupDatabase(dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		err = db.Close()
		require.NoError(t, err)
		err = os.Remove(dbPath)
		require.NoError(t, err)
	}()

	// Content to be saved
	content := [][]string{
		{"2022-01-01 00:00:00", "Test entry 1"},
		{"2022-01-02 00:00:00", "Test entry 2"},
	}

	// Call the function to save the data to the database
	err = saveToDb(db, content)
	require.NoError(t, err)

	// Check if the content was saved correctly
	rows, err := db.Query("SELECT timestamp, content FROM entries")
	require.NoError(t, err)

	i := 0
	var timestamp, savedContent string
	for rows.Next() {
		err = rows.Scan(&timestamp, savedContent)
		require.NoError(t, err)
		assert.Equal(t, content[i][0], timestamp)
		assert.Equal(t, content[i][1], savedContent)
		i++
	}

	assert.Equal(t, len(content), i)

	err = rows.Close()
	require.NoError(t, err)
}

func Test_listDbContents(t *testing.T) {
	// Create a temporary database file
	fileName := "test_list_til.db"
	dbPath, _ := xdg.DataFile(filepath.Join("til", fileName))

	db, err := setupDatabase(dbPath)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() {
		err = db.Close()
		require.NoError(t, err)
		err = os.Remove(dbPath)
		require.NoError(t, err)
	}()

	// Fill the table with content
	content := [][]string{
		{"2022-01-01 00:00:00", "Test entry 1"},
		{"2022-01-02 00:00:00", "Test entry 2"},
	}
	err = saveToDb(db, content)
	require.NoError(t, err)

	// Call the function to be tested
	var buffer bytes.Buffer
	buffer.WriteString("Timestamp\t\t\tTIL\n")
	buffer.WriteString("2022-01-01 00:00:00\t")
	buffer.WriteString("Test entry 1\n")
	buffer.WriteString("2022-01-02 00:00:00\t")
	buffer.WriteString("Test entry 2\n")

	expected := buffer.String()

	// Capture the output of the function
	output := captureOutput(func() {
		err = listDbContents(db)
		require.NoError(t, err)
	})

	assert.Equal(t, expected, output)
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
		"test_til.db",
		"test_save_til.db",
		"test_list_til.db",
  }

	for _, fileName := range fileNames {
		dataPath, _ := xdg.DataFile(filepath.Join("til", fileName))
		os.Remove(dataPath)
	}

	os.Exit(retCode)
}
