up:
	docker compose up -d --build

down:
	docker compose down

rebuild: clean up

clean:
	docker compose down -v

logs:
	docker compose logs -f app

test:
	go test -v ./...

fmt:
	go fmt ./...