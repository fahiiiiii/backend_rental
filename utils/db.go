package utils

import (
    "fmt"
    "strconv"
    "database/sql"
    "io/ioutil"  
    "encoding/json"  
    // "backend_rental/services"
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
    // orm.RegisterModel(new(models.RentalProperty))
    // orm.RegisterModel(new(models.PropertyDetails))

    // Create tables
    err = orm.RunSyncdb("default", false, true)
    if err != nil {
        return fmt.Errorf("failed to sync database: %v", err)
    }
    err = loadRentalPropertyData()
    if err != nil {
        return fmt.Errorf("failed to load rental property data: %v", err)
    }
    // service := &services.PropertyDetailsServiceDB{}
    err = LoadPropertyDetailsFromJSON()
    if err != nil {
        return fmt.Errorf("failed to load property details: %v", err)
    }
    fmt.Println("Database initialized successfully")
    return nil
}


func loadRentalPropertyData() error {
    // Read the RentalProperty.json file
    data, err := ioutil.ReadFile("data/RentalProperty.json")
    if err != nil {
        // If file doesn't exist, it's not an error - just skip loading
        fmt.Println("RentalProperty.json not found. Skipping data loading.")
        return nil
    }

    var properties []models.RentalProperty
    err = json.Unmarshal(data, &properties)
    if err != nil {
        return fmt.Errorf("failed to parse RentalProperty.json: %v", err)
    }

    fmt.Printf("Loaded %d properties from JSON\n", len(properties))

    // Get database connection
    db, err := orm.GetDB("default")
    if err != nil {
        return fmt.Errorf("failed to get database connection: %v", err)
    }

    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to start transaction: %v", err)
    }

    // Prepare to defer transaction rollback or commit
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        tx.Commit()
    }()

    // Clear existing data
    _, err = tx.Exec("DELETE FROM rental_property")
    if err != nil {
        return fmt.Errorf("failed to clear existing data: %v", err)
    }

    // Prepare insert statement
    stmt, err := tx.Prepare("INSERT INTO rental_property (city_id, property_id, name, property_type, bedrooms, bathrooms, amenities) VALUES ($1, $2, $3, $4, $5, $6, $7)")
    if err != nil {
        return fmt.Errorf("failed to prepare insert statement: %v", err)
    }
    defer stmt.Close()

    // Insert new data
    for _, prop := range properties {
        _, err = stmt.Exec(prop.CityID, prop.PropertyID, prop.Name, prop.PropertyType, prop.Bedrooms, prop.Bathrooms, prop.Amenities)
        if err != nil {
            return fmt.Errorf("failed to insert property %v: %v", prop.PropertyID, err)
        }
    }

    fmt.Printf("Successfully inserted %d properties\n", len(properties))
    return nil
}

func LoadPropertyDetailsFromJSON() error {
    // Read the PropertyDetails.json file
    data, err := ioutil.ReadFile("data/PropertyDetails.json")
    if err != nil {
        return fmt.Errorf("failed to read PropertyDetails.json: %v", err)
    }
 
    var propertyDetails []models.PropertyDetails
    err = json.Unmarshal(data, &propertyDetails)
    if err != nil {
        return fmt.Errorf("failed to parse PropertyDetails.json: %v", err)
    }
 
    fmt.Printf("Loaded %d property details from JSON\n", len(propertyDetails))
 
    // Get database connection
    db, err := orm.GetDB("default")
    if err != nil {
        return fmt.Errorf("failed to get database connection: %v", err)
    }
 
    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to start transaction: %v", err)
    }
 
    // Prepare to defer transaction rollback or commit
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        tx.Commit()
    }()
 
    // Clear existing data
    _, err = tx.Exec("DELETE FROM property_details")
    if err != nil {
        return fmt.Errorf("failed to clear existing data: %v", err)
    }
 
    // Prepare insert statement
    stmt, err := tx.Prepare(`
        INSERT INTO property_details 
        (property_id, description, review_score, review_count, review_score_word, image_type, image_urls) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `)
    if err != nil {
        return fmt.Errorf("failed to prepare insert statement: %v", err)
    }
    defer stmt.Close()
 
    // Insert new data
    for _, detail := range propertyDetails {
        // Convert ImageUrls to JSON string
        imageUrlsJSON, err := json.Marshal(detail.ImageUrls)
        if err != nil {
            return fmt.Errorf("failed to marshal image URLs: %v", err)
        }
 
        _, err = stmt.Exec(
            detail.PropertyID, 
            detail.Description, 
            detail.ReviewScore, 
            detail.ReviewCount, 
            detail.ReviewScoreWord, 
            detail.ImageType, 
            string(imageUrlsJSON),
        )
        if err != nil {
            return fmt.Errorf("failed to insert property detail %v: %v", detail.PropertyID, err)
        }
    }
 
    fmt.Printf("Successfully inserted %d property details\n", len(propertyDetails))
    return nil
 }
// package utils

// import (
//     "fmt"
//     "strconv"
//     "database/sql"
//     "backend_rental/models"
//     "github.com/beego/beego/v2/client/orm"
//     "github.com/beego/beego/v2/core/config"
//     _ "github.com/lib/pq"
// )

// type DBConfig struct {
//     Host     string
//     Port     int
//     User     string
//     Password string
//     Name     string
//     SSLMode  string
// }

// func getDBConfig() (*DBConfig, error) {
//     // Get database configuration directly using environment keys
//     host, err := config.String("db::host")
//     if err != nil || host == "" {
//         host = "localhost"
//     }

//     portStr, err := config.String("db::port")
//     port := 5432
//     if err == nil && portStr != "" {
//         if p, err := strconv.Atoi(portStr); err == nil {
//             port = p
//         }
//     }

//     user, err := config.String("db::user")
//     if err != nil || user == "" {
//         user = "fahimah" // default from your config
//     }

//     password, err := config.String("db::password")
//     if err != nil || password == "" {
//         password = "fahimah123" // default from your config
//     }

//     name, err := config.String("db::name")
//     if err != nil || name == "" {
//         name = "rental_db" // default from your config
//     }

//     sslmode, err := config.String("db::sslmode")
//     if err != nil || sslmode == "" {
//         sslmode = "disable"
//     }

//     fmt.Printf("Database configuration loaded: host=%s, port=%d, user=%s, dbname=%s, sslmode=%s\n",
//         host, port, user, name, sslmode)

//     return &DBConfig{
//         Host:     host,
//         Port:     port,
//         User:     user,
//         Password: password,
//         Name:     name,
//         SSLMode:  sslmode,
//     }, nil
// }

// func createDatabaseIfNotExists(config *DBConfig) error {
//     // Connect to PostgreSQL server (postgres database) to create new database
//     connStr := fmt.Sprintf(
//         "host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
//         config.Host,
//         config.Port,
//         config.User,
//         config.Password,
//         config.SSLMode,
//     )

//     db, err := sql.Open("postgres", connStr)
//     if err != nil {
//         return fmt.Errorf("error connecting to postgres database: %v", err)
//     }
//     defer db.Close()

//     // Check if database exists
//     var exists bool
//     err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_database WHERE datname = $1)", config.Name).Scan(&exists)
//     if err != nil {
//         return fmt.Errorf("error checking if database exists: %v", err)
//     }

//     // Create database if it doesn't exist
//     if !exists {
//         fmt.Printf("Creating database %s...\n", config.Name)
//         _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Name))
//         if err != nil {
//             return fmt.Errorf("error creating database: %v", err)
//         }
//         fmt.Printf("Database %s created successfully\n", config.Name)
//     }

//     return nil
// }

// func InitDB() error {
//     // Get database configuration
//     dbConfig, err := getDBConfig()
//     if err != nil {
//         return fmt.Errorf("failed to get database config: %v", err)
//     }

//     // Create database if it doesn't exist
//     err = createDatabaseIfNotExists(dbConfig)
//     if err != nil {
//         return fmt.Errorf("failed to create database: %v", err)
//     }

//     // Create connection string for the actual database
//     dbURL := fmt.Sprintf(
//         "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
//         dbConfig.Host,
//         dbConfig.Port,
//         dbConfig.User,
//         dbConfig.Password,
//         dbConfig.Name,
//         dbConfig.SSLMode,
//     )

//     // Register the database driver
//     err = orm.RegisterDriver("postgres", orm.DRPostgres)
//     if err != nil {
//         return fmt.Errorf("failed to register driver: %v", err)
//     }

//     // Register the default database
//     err = orm.RegisterDataBase("default", "postgres", dbURL)
//     if err != nil {
//         return fmt.Errorf("failed to register database: %v", err)
//     }

//     // Register models
//     orm.RegisterModel(new(models.Location))

//     // Create tables
//     err = orm.RunSyncdb("default", false, true)
//     if err != nil {
//         return fmt.Errorf("failed to sync database: %v", err)
//     }

//     fmt.Println("Database initialized successfully")
//     return nil
// }