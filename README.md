# Backend Rental

## Overview
Backend Rental is a Go-based backend service that handles property-related data. The project uses the Beego framework and provides several API endpoints to fetch and generate property-related JSON data.

## Prerequisites
- Go (latest version recommended)
- Beego framework (`go get github.com/beego/beego/v2`)
- Docker & Docker Compose

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/fahiiiiii/backend_rental
   cd backend_rental
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Running the Project
You can run the project using either of the following commands:

### Running with Docker
To build and activate the Docker container, use the following commands:
```sh
docker-compose build
docker-compose up
```

### Running Locally
```sh
# Recommended way
go run main.go

# Alternative way using Bee tool
bee run
```

## Viewing the Database
To access the database inside the Docker container, run:
```sh
docker exec -it backend_rental-db-1 psql -U fahimah -d rental_db
```

Once inside PostgreSQL, you can run the following commands to inspect the tables and data:
```sql
\dt  -- List all tables
SELECT * FROM location;
SELECT * FROM property_details;
SELECT * FROM rental_property;
```

## API Endpoints & Usage
The following API endpoints are available:

### Data Fetching APIs
| Endpoint | Controller | Description | cURL Command |
|----------|------------|-------------|--------------|
| `/v1/city` | `CityController` | Fetches city data and saves it to `data/cities.json` | `curl http://localhost:8080/v1/city` |
| `/v1/properties` | `PropertyController` | Fetches properties for all cities and saves to `data/properties.json` | `curl http://localhost:8080/v1/properties` |
| `/v1/property-details` | `PropertyDetailController` | Fetches property details and saves to `data/property_details.json` | `curl http://localhost:8080/v1/property-details` |
| `/v1/property-description` | `PropertyDescriptionController` | Fetches property descriptions and saves to `data/property_desc_image.json` | `curl http://localhost:8080/v1/property-description` |
| `/v1/property-images` | `PropertyImageController` | Fetches property images and saves to `data/property_images.json` | `curl http://localhost:8080/v1/property-images` |

### Data Generation APIs
| Endpoint | Controller | Description | cURL Command |
|----------|------------|-------------|--------------|
| `/generate-rental-property` | `RentalPropertyController` | Generates `RentalProperty.json` from other JSON files | `curl http://localhost:8080/generate-rental-property` |
| `/generate-property-details` | `PropertyDetailsControllerJSON` | Generates `PropertyDetails.json` from other JSON files | `curl http://localhost:8080/generate-property-details` |

### Data Serving APIs
| Endpoint | Controller | Description | cURL Command |
|----------|------------|-------------|--------------|
| `/v1/property/list` | `PropertyListController` | Serves property list data | `curl http://localhost:8080/v1/property/list` |
| `/v1/property/details` | `PropertyDetailControllerDB` | Serves property details data | `curl http://localhost:8080/v1/property/details` |

## Contribution
Feel free to fork the repository and submit pull requests for improvements.

## License
This project is open-source and available under the [MIT License](LICENSE).

