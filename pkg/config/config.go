package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type MariaDBConfig struct {
	DBHost               string
	DBPort               string
	DBUser               string
	DBName               string
	DBPassword           string
	TLSConfig            string
	AllowNativePasswords bool
	MultiStatements      bool
	MaxOpenConns         uint
	MaxIdleConns         uint
	ConnMaxLifetime      uint
}

type RedisConfig struct {
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RedisDBNum      uint8
	RedisExpiration uint
}

type LoggingConfig struct {
	ServerLogPath string
	ErrLogPath    string
	LogMaxSize    int
	LogMaxBackup  int
	LogMaxAge     int
	LogCompress   bool
}

type SSLConfig struct {
	KeyPath  string
	CertPath string
}

type FileHandlingConfig struct {
	MaxFileSize      int64
	AllowedExtension []string
	UploadDir        string
}

type SecurityConfig struct {
	AllowedCharacters string
	MinPasswordLength int
	SaltSize          int
	IterationCount    int
	KeySize           int
	EncryptionKey     string
}

type AppConfig struct {
	MariaDBConfig
	RedisConfig
	LoggingConfig
	SSLConfig
	FileHandlingConfig
	SecurityConfig
	AppMode     string
	AppLanguage string
	AppTimezone string
	AppPort     string
	GRPCPort    string
	AppHost     string
	AppRoot     string
	DebugMode   bool
}

func New(envPath string) *AppConfig {

	if err := godotenv.Load(envPath); err != nil {
		log.Println("Failed to locate .env file, program will proceed with provided env if any is provided")
	}

	config := &AppConfig{
		MariaDBConfig: MariaDBConfig{
			DBHost:               getEnv("DB_HOST", ""),
			DBPort:               getEnv("DB_PORT", "3306"),
			DBUser:               getEnv("DB_USER", ""),
			DBPassword:           getEnv("DB_PASSWORD", ""),
			DBName:               getEnv("DB_NAME", ""),
			TLSConfig:            getEnv("DB_TLS_CONFIG", "false"),
			AllowNativePasswords: getEnvAsBool("DB_ALLOW_NATIVE_PASSWORDS", true),
			MultiStatements:      getEnvAsBool("DB_MULTI_STATEMENTS", false),
			MaxOpenConns:         uint(getEnvAsInt("DB_MAX_OPEN_CONNS", 20)),
			MaxIdleConns:         uint(getEnvAsInt("DB_MAX_IDLE_CONNS", 5)),
			ConnMaxLifetime:      uint(getEnvAsInt("DB_CONN_MAX_LIFETIME", 5)),
		},
		RedisConfig: RedisConfig{
			RedisHost:       getEnv("REDIS_HOST", ""),
			RedisPort:       getEnv("REDIS_PORT", "6378"),
			RedisPassword:   getEnv("REDIS_PASSWORD", ""),
			RedisDBNum:      uint8(getEnvAsInt("REDIS_DB_NUM", 0)),
			RedisExpiration: uint(getEnvAsInt("REDIS_EXPIRATION", 0)),
		},
		LoggingConfig: LoggingConfig{
			ServerLogPath: getEnv("LOG_PATH", "./log/server.log"),
			ErrLogPath:    getEnv("ERR_LOG_PATH", "./log/error.log"),
			LogMaxSize:    getEnvAsInt("LOG_MAX_SIZE", 30),
			LogMaxBackup:  getEnvAsInt("LOG_MAX_BACKUP", 5),
			LogMaxAge:     getEnvAsInt("LOG_MAX_AGE", 30),
			LogCompress:   getEnvAsBool("LOG_COMPRESS", true),
		},
		SSLConfig: SSLConfig{
			KeyPath:  getEnv("KEY_PATH", ""),
			CertPath: getEnv("CERT_PATH", ""),
		},
		FileHandlingConfig: FileHandlingConfig{
			MaxFileSize:      int64(getEnvAsInt("MAX_FILE_SIZE", 1024*20)), // 20 MB Max
			AllowedExtension: getEnvAsSlice("ALLOWED_FILE_EXTENSIONS", []string{"jpg", "jpeg", "png"}, ","),
			UploadDir:        getEnv("UPLOAD_DIR", "assets/vendors"),
		},
		SecurityConfig: SecurityConfig{
			AllowedCharacters: getEnv("ALLOWED_CHARACTERS", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]~"),
			MinPasswordLength: getEnvAsInt("MIN_PASSWORD_LENGTH", 16),
			SaltSize:          getEnvAsInt("SALT_SIZE", 16),
			IterationCount:    getEnvAsInt("ITERATION_COUNT", 4096),
			KeySize:           getEnvAsInt("KEY_SIZE", 32),
			EncryptionKey:     getEnv("ENCRYPTION_KEY", ""),
		},
		AppMode:     getEnv("APP_MODE", "devs"),
		AppLanguage: getEnv("APP_LANG", "en"),
		AppTimezone: getEnv("APP_TIMEZONE", "Asia/Jakarta"),
		AppPort:     getEnv("APP_PORT", ""),
		GRPCPort:    getEnv("GRPC_PORT", ""),
		AppHost:     getEnv("APP_HOST", ""),
		AppRoot:     getEnv("APP_ROOT", "/api/v1"),
		DebugMode:   getEnvAsBool("DEBUG", false),
	}

	return config
}

// Simple helper function to read an environment or return a default value.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if nextValue := os.Getenv(key); nextValue != "" {
		return nextValue
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value.
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a slice of a specific type or return default value.
func getEnvAsSlice[T any](name string, defaultVal []T, sep string) []T {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	vals := strings.Split(valStr, sep)
	result := make([]T, len(vals))

	for i, v := range vals {
		switch any(result).(type) {
		case []string:
			result[i] = any(v).(T)
		case []int:
			intVal, _ := strconv.Atoi(v)
			result[i] = any(intVal).(T)
		case []bool:
			boolVal, _ := strconv.ParseBool(v)
			result[i] = any(boolVal).(T)
		default:
			return defaultVal
		}
	}

	return result
}
