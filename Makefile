APP?=api
PORT?=8080
CONTAINER_IMAGE?=docker.io/webdeva/${APP}

clean:
	rm -f ${APP}

build: clean
	go build \
		-ldflags "-s -w -X version=${RELEASE} \
		-X buildcommit=${COMMIT} -X buildtime=${BUILD_TIME}" \
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