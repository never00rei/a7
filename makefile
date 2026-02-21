APP_EXECUTABLE=a7
.PHONY: build build-cross build-target run clean

build:
	go build -o ${APP_EXECUTABLE} main.go

build-cross:
	GOOS=darwin GOARCH=amd64 go build -o ${APP_EXECUTABLE}-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -o ${APP_EXECUTABLE}-linux-amd64 main.go

build-target:
	@test -n "$(GOOS)" || (echo "GOOS is required (example: make build-target GOOS=linux GOARCH=arm64)"; exit 1)
	@test -n "$(GOARCH)" || (echo "GOARCH is required (example: make build-target GOOS=linux GOARCH=arm64)"; exit 1)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ${APP_EXECUTABLE}-$(GOOS)-$(GOARCH) main.go

run: build
	./${APP_EXECUTABLE}

clean:
	go clean
	rm -f ${APP_EXECUTABLE} ${APP_EXECUTABLE}-darwin-amd64 ${APP_EXECUTABLE}-linux-amd64
