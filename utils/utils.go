package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func CreatePath(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	return nil
}

func OpenFile(path string, filename string) (string, error) {
	var content []byte

	f := filepath.Join(path, filename)

	content, err := os.ReadFile(f)
	if err != nil {
		return string(content), fmt.Errorf("could not open file: %w", err)
	}

	return string(content), nil

}

func SaveFile(path string, filename string, content string) error {

	f := filepath.Join(path, filename)

	// This may seem ambiguous so to explain:
	// os.O_APPEND - adds to the original if it exists.
	// os.O_CREATE - creates the file if it doesn't exist.
	// os.O_WRONLY - this only allows writing to the file.
	file, err := os.OpenFile(f,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0755)
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("could not write content to file: %w", err)
	}

	return nil
}

func RandomStringFromSlice(slice []string) string {
	var str string

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	randIndex := seed.Intn(len(slice))

	str = slice[randIndex]

	return str
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func SanitizeSpecialChars(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	result := re.ReplaceAllString(input, "_")
	return result
}
