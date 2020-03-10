APP?=api
PORT?=8080

clean:
	rm -f ${APP}

build: clean
	go build \
		-ldflags "-s -w -X version=${RELEASE} \
		-X buildcommit=${COMMIT} -X buildtime=${BUILD_TIME}" \
		-o ${APP}

run: build
	PORT=${PORT} ./${APP}

test:
	go test -v -race ./...