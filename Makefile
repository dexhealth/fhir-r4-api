run:
	@docker start mongo
	@go run ./cmd

# test need to be re-writtern for r4
test:
	@docker stop mongo
	@docker rm mongo
	@docker run --name mongo \
		-p 27017:27017 \
		-e MONGO_INITDB_ROOT_USERNAME=root \
		-e MONGO_INITDB_ROOT_PASSWORD=example \
		-d mongo
	@go test -v ./cmd/testing
