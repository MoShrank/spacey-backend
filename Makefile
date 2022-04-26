test:
	- ENV_FILE_PATH="$(pwd)/.env" go test ./...

serve:
	- for file in ./services/*; do \
    	echo "Building $$file"; \
    	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $$file/bin/main ./$$file; \
	  done
	- docker-compose build
	- docker-compose up -d

cleanup:
	- docker-compose down --remove-orphans