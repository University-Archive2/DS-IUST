build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app main.go
docker-build:
	docker build -t alinazariii/ds-ingestion:$(VERSION) -f Dockerfile .
docker-push:
	docker push alinazariii/ds-ingestion:$(VERSION)
