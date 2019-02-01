package main

import (
	"os"
)

type Config struct {
	AppUrl              string
	LocalVideoDirectory string
	HttpVideoDirectory  string
	Port                string
	UsePostgres         bool
	PostgresUser        string
	PostgresPassword    string
	PostgresAddress     string
	PostgresDatabase    string
}

func envDefault(name string, backup string) string {
	found, exists := os.LookupEnv(name)
	if exists {
		return found
	}
	return backup
}

func FillFromEnv() Config {
	return Config{
		AppUrl:              envDefault("CREAMY_APP_URL", "http://localhost:3000"),
		LocalVideoDirectory: envDefault("CREAMY_VIDEO_DIR", "dummyvideos"),
		HttpVideoDirectory:  envDefault("CREAMY_HTTP_VIDEO_DIR", "/static/videos/"),
		Port:                envDefault("CREAMY_HTTP_PORT", "3000"),
		UsePostgres:         envDefault("CREAMY_POSTGRES", "false") == "true",
		PostgresUser:        envDefault("CREAMY_POSTGRES_USER", "postgres"),
		PostgresPassword:    envDefault("CREAMY_POSTGRES_PASSWORD", "postgres"),
		PostgresDatabase:    envDefault("CREAMY_POSTGRES_DATABASE", "postgres"),
		PostgresAddress:     envDefault("CREAMY_POSTGRES_ADDRESS", "localhost:5432"),
	}
}
