package utils

import (
	"strconv"
	"strings"
)

func ParseCSVToInt(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	var result []int

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		num, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}

		result = append(result, num)
	}

	return result, nil
}

func ParseCSVToString(input string) []string {
	parts := strings.Split(input, ",")
	var result []string

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		result = append(result, p)
	}

	return result
}
