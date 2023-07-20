package main

import (
	"log"
	"os"
)

func DieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Env(key string, def string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		if def != "" {
			return def
		}

		log.Fatalf("Config key %s is required but not set!", key)
	}

	return value
}

func SliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
