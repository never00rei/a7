APP_EXECUTABLE=a7
BUILD_DIR=build/bin/
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_TARGET=$(BUILD_DIR)/$(APP_EXECUTABLE)-$(GOOS)-$(GOARCH)

build:
	mkdir -p "${BUILD_DIR}"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_TARGET) main.go

run: build
	./$(BUILD_TARGET)

clean:
	go clean
	rm -f $(BUILD_TARGET)
