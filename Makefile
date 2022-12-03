start-swagger:
	swag init --output ./docs

delete-swagger:
	rm -rf ./docs

run:
	go run main.go