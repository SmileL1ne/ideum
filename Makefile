run:
	go run ./cmd/api
build:
	docker build -t forum .
docker-run:
	docker run -p 5000:5000 forum
clean:
	docker system prune -af