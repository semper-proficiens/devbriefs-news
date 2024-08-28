package utils

import (
	"fmt"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// CalculateSimilarity calculates the similarity score between two strings using Levenshtein distance.
func CalculateSimilarity(a, b string) float64 {
	distance := levenshtein.DistanceForStrings([]rune(a), []rune(b), levenshtein.DefaultOptions)
	maxLen := maxInt(len(a), len(b))
	fmt.Printf("a: %s, b: %s, distance: %d, maxLen: %d\n", a, b, distance, maxLen)
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(distance)/float64(maxLen)
}

// maxInt is a helper function to find the maximum of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
