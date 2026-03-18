.PHONY: docker_build docker_run docker_unit_test docker_ci

docker_build:
	docker build -t api-petstore-service-layer .

docker_run: docker_build
	docker run --rm -p 8080:8080 --name api-petstore-service-layer api-petstore-service-layer

docker_unit_test:
	docker run --rm -v "$(PWD)":/app -w /app golang:1.25-alpine sh -c "go test ./..."

docker_ci: docker_build docker_unit_test