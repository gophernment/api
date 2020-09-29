APP?=api
PORT?=8080
CONTAINER_IMAGE?=docker.io/webdeva/${APP}
RELEASE=v1.0.0
COMMIT=`git rev-parse --short HEAD`
BUILD_TIME=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`

clean:
	rm -f ${APP}

build: clean
	go build \
		-ldflags "-s -w -X main.version=${RELEASE} \
		-X main.buildcommit=${COMMIT} -X main.buildtime=${BUILD_TIME}" \
		-o ${APP}

container: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

push: container
	docker push $(CONTAINER_IMAGE):$(RELEASE)

ship: container
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		$(APP):$(RELEASE)

run: build
	PORT=${PORT} ./${APP}

test:
	go test -v -race ./...