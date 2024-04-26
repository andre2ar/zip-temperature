run:
	go run main.go

live-reload:
	air -c .air.linux.conf

up:
	docker-compose up