package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"log"
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
	XSRFKeyB64          string
	XSRFKey             []byte
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
	cfg := appConfig{
		AppURL:              envDefault("CREAMY_APP_URL", ""),
		LocalVideoDirectory: envDefault("CREAMY_VIDEO_DIR", "dummyvideos"),
		HTTPVideoDirectory:  envDefault("CREAMY_HTTP_VIDEO_DIR", "/static/videos/"),
		Port:                envDefault("CREAMY_HTTP_PORT", "3000"),
		UsePostgres:         envDefault("CREAMY_POSTGRES", "false") == "true",
		PostgresUser:        envDefault("CREAMY_POSTGRES_USER", "postgres"),
		PostgresPassword:    envDefault("CREAMY_POSTGRES_PASSWORD", "postgres"),
		PostgresDatabase:    envDefault("CREAMY_POSTGRES_DATABASE", "postgres"),
		PostgresAddress:     envDefault("CREAMY_POSTGRES_ADDRESS", "localhost:5432"),
		XSRFKeyB64:          envDefault("CREAMY_XSRF_KEY_B64", ""),
		FilesystemKey:       0x69, // hardcoded for now
		ReadOnly:            envDefault("CREAMY_READ_ONLY", "false") == "true",
	}

	if cfg.XSRFKeyB64 == "" && !cfg.ReadOnly {
		cfg.XSRFKeyB64 = randomXSRFKeyB64()
	}
	if cfg.XSRFKeyB64 != "" {
		var err error
		cfg.XSRFKey, err = base64.StdEncoding.DecodeString(cfg.XSRFKeyB64)
		if err != nil {
			log.Fatal("CREAMY_XSRF_KEY_B64 is set to an invalid value:", err)
		}
	}

	return cfg
}

func randomXSRFKeyB64() string {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("CREAMY_XSRF_KEY_B64 is unset and an error was encountered during generation:", err)
	}
	str := base64.StdEncoding.EncodeToString(bytes)
	log.Printf("CREAMY_XSRF_KEY_B64 is not specified, using CREAMY_XSRF_KEY_B64=%v (XSRF will be invalid upon restart)", str)
	return str
}
