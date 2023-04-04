package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getEnvWithDefaultAndBlankableFlag(key string, defaultValue string, canBeBlank bool) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if !canBeBlank && defaultValue == "" {
			msg := fmt.Sprintf("No value for key %s which cannot be blank.", key)
			log.Fatal(msg)
		}
		return defaultValue
	}
	return value
}

func getEnvWithDefault(key string, defaultValue string) string {
	return getEnvWithDefaultAndBlankableFlag(key, defaultValue, false)
}

func getEnv(key string) string {
	return getEnvWithDefault(key, "")
}

func main() {
	envpath := os.Getenv("ENVPATH") // must be supplied at invocation

	err := godotenv.Load(envpath)
	if err != nil {
		msg := fmt.Sprintf("Could not load env file at path [%s]", envpath)
		log.Fatal(msg)
	}

	env := getEnvWithDefault("APP_ENV", "prod")
	pguser := getEnv("POSTGRES_USER")
	pgpwd := getEnv("POSTGRES_PASSWORD")
	pgdb := getEnv("POSTGRES_DB")
	pghost := getEnv("PG_HOST")
	pgport := getEnvWithDefault("PG_PORT", "5432")

	mainSchema := getEnv("MAIN_SCHEMA")
	altSchema := getEnv("ALT_SCHEMA")

	// TODO: determine next migration
	byteArr, errFile := os.ReadFile("./dbserver/migrations/init.sql")
	templateStr := string(byteArr)

	if errFile != nil {
		log.Fatal("Cannot read migration file!")
	}

	// parse/replace
	replacer := strings.NewReplacer(
		"${DB_NAME}", pgdb,
		"${MAIN_SCHEMA}", mainSchema,
		"${ALT_SCHEMAN}", altSchema,
		// TODO: shld really be sep user; proceeding as-is for simplicity for now
		"${APP_USER}", pguser,
		"${APP_USER_PASS}", pgpwd,
	)
	migrationStr := replacer.Replace(templateStr)

	// connect to db
	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pguser, pgpwd, pghost, pgport, pgdb)
	if strings.EqualFold(env, "local") {
		pgUrl += "?sslmode=disable"
	}

	fmt.Print(migrationStr)
}
