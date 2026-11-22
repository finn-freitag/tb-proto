.PHONY: generate
generate:
	protoc -I proto \
	  --go_out=gen --go_opt=paths=source_relative \
	  --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	  proto/tossbox/agent/v1/agent.proto
	mkdir -p gen/agent/v1
	mv gen/tossbox/agent/v1/*.go gen/agent/v1/ && rm -rf gen/tossbox
