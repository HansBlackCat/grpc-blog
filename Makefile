PROTO_LOC=proto/blog.proto
SERVER_CMD=./server/cmd/*.go
CLIENT_CMD=./client/cmd/*.go

GO_OUT=plugins=grpc:.

PROTOC=protoc
GOC=go

.PHONY: protoc server client

protoc:
	${PROTOC} ${PROTO_LOC} --go_out=${GO_OUT}

server:
	${GOC} run ${SERVER_CMD}

client:
	${GOC} run ${CLIENT_CMD}

	