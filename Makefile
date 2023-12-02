BINARY_NAME=weather

# Use INSTALL_PATH from the environment, or default to "/usr/bin/"
INSTALL_PATH ?= "/usr/bin/"

all: build run

install:
	cd weather && go build -o "${INSTALL_PATH}${BINARY_NAME}"

build:
	go build -o "${BINARY_NAME}"

run:
	go build -o ${BINARY_NAME}
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}