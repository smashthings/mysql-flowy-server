package main

import "os"
import "fmt"

var (
	RequiredEnvVars []string = []string{
		"DB_DSN",
	}
)

func CheckReqEnvVars() {
	for _, x := range RequiredEnvVars {
		if os.Getenv(x) == "" {
			panic(fmt.Sprintf("Could not find required environment variable %s, exiting!", x))
		}
	}
}

func SliceContainsString(s []string, str string) bool {
	for _, x := range s {
		if str == x {
			return true
		}
	}
	return false
}