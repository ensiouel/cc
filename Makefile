docker-build:
	docker build --tag cc .

docker-run:
	docker run -p 8081:8080 cc