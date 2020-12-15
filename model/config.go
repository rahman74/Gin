package model

// Config struct holding all sub details
type Config struct {
	Server                       ServerConfig
	DatabaseConfig                     DatabaseConfig
	Log LogConfig
	Jwt JwtConfig
}

// ServerConfig ...
type ServerConfig struct {
	Port        int
	Environment string
}

// DatabaseConfig ...
type DatabaseConfig struct {
	Host     string
	Port     int
	DbName   string
	Username string
	Password string
}

// LogConfig ...
type LogConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// JwtConfig ...
type JwtConfig struct {
	Issuer                     string
	ExpirationDurationInMinute int
	SignatureKey               string
	ExpirationDurationVerifyInMinute int
}