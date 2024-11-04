package util

import (
	"math"
	"os"
	"strconv"
	"strings"
)

// RoundToOneDecimal rounds a float64 to 1 decimal places
func RoundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

// RoundToTwoDecimal rounds a float64 to 2 decimal places
func RoundToTwoDecimal(value float64) float64 {
	return math.Round(value*100) / 100
}

// LoadFileData loads file content as a string
func LoadFileData(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ParseFloat parses float values
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(s, "$", ""), ",", ""), 64)
}
