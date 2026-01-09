APP_EXECUTABLE=a7

build:
	GOARCH=amd64 GOOS=darwin go build -o ${APP_EXECUTABLE}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${APP_EXECUTABLE}-linux main.go

run: build
	./${APP_EXECUTABLE}

clean:
	go clean
	rm ${APP_EXECUTABLE}-darwin
	rm ${APP_EXECUTABLE}-linux
