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
export SLIDESHOW_API_ID=--your Auth0 API Identifier--
export SLIDESHOW_DOMAIN=--your Auth0 domain--
```

Build and run
```
go run .
```

### Local docker build
```
docker build --tag slideshow-data-image .
docker run --name slideshow-data -p 8080:8080 -d slideshow-data-image -e SLIDESHOW_API_ID='--your Auth0 API Identifier--' -e SLIDESHOW_DOMAIN='--your Auth0 domain--'

docker stop slideshow-data
docker rm slideshow-data
docker rmi slideshow-image
```

### Docker compose
```
docker-compose up -d
docker-compose down
```

### Run public image
```
docker run -p 8080:8080 -d --name slideshow-data jimareed/slideshow-data

docker rm -f slideshow-data
```
