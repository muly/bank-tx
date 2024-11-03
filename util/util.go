package util

import (
	"math"
	"os"
)

// RoundToOneDecimal rounds a float64 to 1 decimal places
func RoundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

// LoadFileData loads file content as a string
func LoadFileData(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
