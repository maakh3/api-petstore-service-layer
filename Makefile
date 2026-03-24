.PHONY: docker_build docker_run docker_utest docker_ci mocks

docker_build:
	docker build -t api-petstore-service-layer .

docker_run: docker_build
	docker run --rm -p 8080:8080 --name api-petstore-service-layer api-petstore-service-layer

docker_utest:
	docker run --rm -v "$(PWD)":/app -w /app golang:1.25-alpine sh -c "go test ./..."

docker_ci: docker_build docker_utest

mocks:
	go generate ./...

utest: mocks
	go test ./...