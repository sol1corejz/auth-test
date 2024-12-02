package config

import (
	"flag"
	"os"
)

var (
	DatabaseURI string
)

func ParseFlags() {

	flag.StringVar(&DatabaseURI, "d", "", "database uri")
	flag.Parse()

	if databaseURI := os.Getenv("DATABASE_URI"); databaseURI != "" {
		DatabaseURI = databaseURI
	}

}
