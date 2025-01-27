# Backend Rental

## Overview
Backend Rental is a Go-based backend service that handles property-related data. The project uses the Beego framework and provides several API endpoints to fetch and generate property-related JSON data.

## Prerequisites
- Go (latest version recommended)
- Beego framework (`go get github.com/beego/beego/v2`)

## Installation
1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd backend_rental
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Running the Project
You can run the project using either of the following commands:
```sh
# Recommended way
go run main.go

# Alternative way using Bee tool
bee run
```

## API Endpoints & Usage
The following API endpoints are available:

### Data Fetching APIs
| Endpoint | Controller | Description | cURL Command |
|----------|------------|-------------|--------------|
| `/v1/city` | `CityController` | Fetches city data and saves it to `data/cities.json` | `curl -X GET http://localhost:8080/v1/city` |
| `/v1/properties` | `PropertyController` | Fetches properties for all cities and saves to `data/properties.json` | `curl -X GET http://localhost:8080/v1/properties` |
| `/v1/property-details` | `PropertyDetailController` | Fetches property details and saves to `data/property_details.json` | `curl -X GET http://localhost:8080/v1/property-details` |
| `/v1/property-description` | `PropertyDescriptionController` | Fetches property descriptions and saves to `data/property_desc_image.json` | `curl -X GET http://localhost:8080/v1/property-description` |
| `/v1/property-images` | `PropertyImageController` | Fetches property images and saves to `data/property_images.json` | `curl -X GET http://localhost:8080/v1/property-images` |

### Data Generation APIs
| Endpoint | Controller | Description | cURL Command |
|----------|------------|-------------|--------------|
| `/generate-rental-property` | `RentalPropertyController` | Generates `RentalProperty.json` from other JSON files | `curl -X GET http://localhost:8080/generate-rental-property` |
| `/generate-property-details` | `PropertyDetailsControllerJSON` | Generates `PropertyDetails.json` from other JSON files | `curl -X GET http://localhost:8080/generate-property-details` |

### Data Serving APIs
| Endpoint | Controller | Description | cURL Command |
|----------|------------|-------------|--------------|
| `/v1/property/list` | `PropertyListController` | Serves property list data | `curl -X GET http://localhost:8080/v1/property/list` |
| `/v1/property/details` | `PropertyDetailControllerDB` | Serves property details data | `curl -X GET http://localhost:8080/v1/property/details` |

## Contribution
Feel free to fork the repository and submit pull requests for improvements.

## License
This project is open-source and available under the [MIT License](LICENSE).

