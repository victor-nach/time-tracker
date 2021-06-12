package config

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	defaultPort   = "8080"
	defaultSecret    = "secret"
	defaultDbUrl  = "mongodb://localhost:27017"
	defaultDbName = "tracker"
)

// Secrets contain all the config that this application needs
type Secrets struct {
	Port        string `json:"port"`
	JWTSecret   string `json:"jwt_secret"`
	DBName      string `json:"db_name"`
	DBURL       string `json:"dburl"`
}

// LoadSecrets loads secrets from the environment and returns it
// if a .env file is present, it would be loaded first
// default values are also set
func LoadSecrets(filename ...string) *Secrets {
	f := ".env"
	if len(filename) > 0 {
		f = filename[0]
	}
	_ = godotenv.Load(f)
	secrets := &Secrets{}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}
	secrets.Port = port

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		jwtSecret = defaultSecret
	}
	secrets.JWTSecret = jwtSecret

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dbUrl = defaultDbUrl
	}
	secrets.DBURL = dbUrl

	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		dbName = defaultDbName
	}
	secrets.DBName = dbName

	return secrets
}
