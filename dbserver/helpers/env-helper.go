package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env string

	PgUser            string
	PgPwd             string
	PgDb              string
	PgHost            string
	PgPort            string
	PgUrl             string
	MigrationsDir     string
	MainSchema        string
	SupertokensSchema string
}

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

func LoadConfig() *Config {
	// must be supplied (e.g. in local dev) if you want to load an env file
	envpath := getEnv("ENVPATH")

	if !IsBlank(envpath) {
		err := godotenv.Load(envpath)
		if err != nil {
			msg := fmt.Sprintf("Could not load env file at path [%s]", envpath)
			log.Fatal(msg)
		}
	}

	env := getEnvWithDefault("APP_ENV", "prod")
	pguser := getEnv("POSTGRES_USER")
	pgpwd := getEnv("POSTGRES_PASSWORD")
	pgdb := getEnv("POSTGRES_DB")
	pghost := getEnv("PG_HOST")
	pgport := getEnvWithDefault("PG_PORT", "5432")
	mainSchema := getEnv("MAIN_SCHEMA")
	supertokensSchema := getEnv("SUPERTOKENS_SCHEMA")
	migDir := getEnv("MIGPATH")

	PGURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pguser, pgpwd, pghost, pgport, pgdb)
	if strings.EqualFold(env, "local") {
		PGURL += "?sslmode=disable"
	}

	config := Config{
		Env:               env,
		PgUser:            pguser,
		PgPwd:             pgpwd,
		PgDb:              pgdb,
		PgHost:            pghost,
		PgPort:            pgport,
		MainSchema:        mainSchema,
		SupertokensSchema: supertokensSchema,
		PgUrl:             PGURL,
		MigrationsDir:     migDir,
	}

	return &config
}

func HydrateSQLTemplate(templateStr string, config Config) string {
	replacer := strings.NewReplacer(
		"${DB_NAME}", config.PgDb,
		"${MAIN_SCHEMA}", config.MainSchema,
		"${SUPERTOKENS_SCHEMA}", config.SupertokensSchema,
		"${APP_USER}", config.PgUser,
		"${APP_USER_PASS}", config.PgPwd,
	)
	migrationStr := replacer.Replace(templateStr)

	return migrationStr
}
