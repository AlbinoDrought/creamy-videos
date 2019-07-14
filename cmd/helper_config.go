package cmd

import (
	"os"
)

type appConfig struct {
	AppURL              string
	LocalVideoDirectory string
	HTTPVideoDirectory  string
	Port                string
	UsePostgres         bool
	PostgresUser        string
	PostgresPassword    string
	PostgresAddress     string
	PostgresDatabase    string
	FilesystemKey       byte
	ReadOnly            bool
}

func envDefault(name string, backup string) string {
	found, exists := os.LookupEnv(name)
	if exists {
		return found
	}
	return backup
}

func makeConfig() appConfig {
	return appConfig{
		AppURL:              envDefault("CREAMY_APP_URL", "http://localhost:3000"),
		LocalVideoDirectory: envDefault("CREAMY_VIDEO_DIR", "dummyvideos"),
		HTTPVideoDirectory:  envDefault("CREAMY_HTTP_VIDEO_DIR", "/static/videos/"),
		Port:                envDefault("CREAMY_HTTP_PORT", "3000"),
		UsePostgres:         envDefault("CREAMY_POSTGRES", "false") == "true",
		PostgresUser:        envDefault("CREAMY_POSTGRES_USER", "postgres"),
		PostgresPassword:    envDefault("CREAMY_POSTGRES_PASSWORD", "postgres"),
		PostgresDatabase:    envDefault("CREAMY_POSTGRES_DATABASE", "postgres"),
		PostgresAddress:     envDefault("CREAMY_POSTGRES_ADDRESS", "localhost:5432"),
		FilesystemKey:       0x69, // hardcoded for now
		ReadOnly:            envDefault("CREAMY_READ_ONLY", "false") == "true",
	}
}
