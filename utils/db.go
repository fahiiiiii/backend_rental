package utils

import (
    "fmt"
    "strconv"
    "database/sql"
    "backend_rental/models"
    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/core/config"
    _ "github.com/lib/pq"
)

type DBConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
    SSLMode  string
}

func getDBConfig() (*DBConfig, error) {
    // Get database configuration directly using environment keys
    host, err := config.String("db::host")
    if err != nil || host == "" {
        host = "localhost"
    }

    portStr, err := config.String("db::port")
    port := 5432
    if err == nil && portStr != "" {
        if p, err := strconv.Atoi(portStr); err == nil {
            port = p
        }
    }

    user, err := config.String("db::user")
    if err != nil || user == "" {
        user = "fahimah" // default from your config
    }

    password, err := config.String("db::password")
    if err != nil || password == "" {
        password = "fahimah123" // default from your config
    }

    name, err := config.String("db::name")
    if err != nil || name == "" {
        name = "rental_db" // default from your config
    }

    sslmode, err := config.String("db::sslmode")
    if err != nil || sslmode == "" {
        sslmode = "disable"
    }

    fmt.Printf("Database configuration loaded: host=%s, port=%d, user=%s, dbname=%s, sslmode=%s\n",
        host, port, user, name, sslmode)

    return &DBConfig{
        Host:     host,
        Port:     port,
        User:     user,
        Password: password,
        Name:     name,
        SSLMode:  sslmode,
    }, nil
}

func createDatabaseIfNotExists(config *DBConfig) error {
    // Connect to PostgreSQL server (postgres database) to create new database
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
        config.Host,
        config.Port,
        config.User,
        config.Password,
        config.SSLMode,
    )

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error connecting to postgres database: %v", err)
    }
    defer db.Close()

    // Check if database exists
    var exists bool
    err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_database WHERE datname = $1)", config.Name).Scan(&exists)
    if err != nil {
        return fmt.Errorf("error checking if database exists: %v", err)
    }

    // Create database if it doesn't exist
    if !exists {
        fmt.Printf("Creating database %s...\n", config.Name)
        _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Name))
        if err != nil {
            return fmt.Errorf("error creating database: %v", err)
        }
        fmt.Printf("Database %s created successfully\n", config.Name)
    }

    return nil
}

func InitDB() error {
    // Get database configuration
    dbConfig, err := getDBConfig()
    if err != nil {
        return fmt.Errorf("failed to get database config: %v", err)
    }

    // Create database if it doesn't exist
    err = createDatabaseIfNotExists(dbConfig)
    if err != nil {
        return fmt.Errorf("failed to create database: %v", err)
    }

    // Create connection string for the actual database
    dbURL := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        dbConfig.Host,
        dbConfig.Port,
        dbConfig.User,
        dbConfig.Password,
        dbConfig.Name,
        dbConfig.SSLMode,
    )

    // Register the database driver
    err = orm.RegisterDriver("postgres", orm.DRPostgres)
    if err != nil {
        return fmt.Errorf("failed to register driver: %v", err)
    }

    // Register the default database
    err = orm.RegisterDataBase("default", "postgres", dbURL)
    if err != nil {
        return fmt.Errorf("failed to register database: %v", err)
    }

    // Register models
    orm.RegisterModel(new(models.Location))

    // Create tables
    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        return fmt.Errorf("failed to sync database: %v", err)
    }

    fmt.Println("Database initialized successfully")
    return nil
}








