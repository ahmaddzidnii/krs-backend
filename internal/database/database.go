package database

import (
	"fmt"
	"time"

	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

// ConnectionPoolConfig holds connection pool configuration
type ConnectionPoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConnectionPoolConfig returns default connection pool configuration
func DefaultConnectionPoolConfig() ConnectionPoolConfig {
	return ConnectionPoolConfig{
		MaxIdleConns:    100,
		MaxOpenConns:    120,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     config.GetEnv("DB_HOST", "localhost"),
		User:     config.GetEnv("DB_USER", "gorm"),
		Password: config.GetEnv("DB_PASSWORD", "gorm"),
		DBName:   config.GetEnv("DB_NAME", "gorm"),
		Port:     config.GetEnv("DB_PORT", "9920"),
		SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
	}
}

// buildDSN builds PostgreSQL DSN from config
func (c DatabaseConfig) buildDSN() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		c.User, c.Password, c.DBName, c.Host, c.Port, c.SSLMode)
}

// InitDatabase initializes database connection with proper error handling
func InitDatabase() (*gorm.DB, error) {
	dbConfig := LoadDatabaseConfig()
	poolConfig := DefaultConnectionPoolConfig()

	return InitDatabaseWithConfig(dbConfig, poolConfig)
}

// InitDatabaseWithConfig initializes database with custom configuration
func InitDatabaseWithConfig(dbConfig DatabaseConfig, poolConfig ConnectionPoolConfig) (*gorm.DB, error) {
	dsn := dbConfig.buildDSN()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	if err := configureConnectionPool(db, poolConfig); err != nil {
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}

	return db, nil
}

// configureConnectionPool sets up database connection pool
func configureConnectionPool(db *gorm.DB, config ConnectionPoolConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from GORM: %w", err)
	}

	// SetMaxIdleConns: Jumlah koneksi yang "diam" atau tidak terpakai di dalam pool.
	// Nilai yang baik adalah sama dengan MaxOpenConns untuk aplikasi dengan beban tinggi
	// agar tidak ada proses pembuatan koneksi baru saat ada lonjakan traffic.
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	// SetMaxOpenConns: Jumlah maksimum koneksi yang boleh dibuka ke database.
	// Ini adalah 'rem' utama Anda untuk melindungi database.
	// Nilai 100-150 adalah titik awal yang kuat untuk 30k users.
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	// SetConnMaxLifetime: Waktu maksimum sebuah koneksi boleh digunakan kembali.
	// Penting untuk load balancing di environment cloud (AWS RDS, Google Cloud SQL)
	// dan untuk menangani restart database dengan mulus.
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// SetConnMaxIdleTime: Waktu maksimum sebuah koneksi boleh idle sebelum ditutup.
	// Ini membantu membebaskan resource di database selama periode sepi.
	// Harus lebih pendek dari ConnMaxLifetime.
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return nil
}

// InitDatabaseForTesting initializes database for testing with minimal connection pool
func InitDatabaseForTesting() (*gorm.DB, error) {
	dbConfig := LoadDatabaseConfig()

	// Use minimal connection pool for testing
	poolConfig := ConnectionPoolConfig{
		MaxIdleConns:    2,
		MaxOpenConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	return InitDatabaseWithConfig(dbConfig, poolConfig)
}

// TestConnection tests database connectivity
func TestConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// CloseDatabase closes database connection gracefully
func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
