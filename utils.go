package main

import (
	"fmt"
	"os"
	"strings"
)

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isU2F(serialNumber string) bool {
	// Hardware token serial numbers are not ARNs and look like this: GAHT12345678
	if strings.HasPrefix(serialNumber, "arn:") {
		split := strings.Split(serialNumber, ":")
		if len(split) > 5 && strings.HasPrefix(split[5], "u2f/") {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
