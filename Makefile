.PHONY: all
.PHONY: write_api
.PHONY: read_api
.PHONY: consumer

all:
	make consumer
	make write_api
	make read_api

write_api:
	go run write_api/main.go &

read_api:
	go run read_api/main.go &

consumer:
	go run consumer/post_event/main.go &
	go run consumer/delete_event/main.go &
	go run consumer/checksum/main.go &
