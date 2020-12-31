# slideshow-data
backend data for slideshow editor

## Setup

### Build & Run

Grab the dependencies
```
go get
```

Set environment variables
```
export DATA_API_ID=--your Auth0 API Identifier--
export DATA_DOMAIN=--your Auth0 domain--
```

Build and run
```
go run .
```

### Docker compose
```
docker build --tag slideshow-data-image .
docker-compose up -d
docker-compose down
```

### Run public image
```
docker run -p 8080:8080 -d --name slideshow-data jimareed/slideshow-data

docker rm -f slideshow-data
```
