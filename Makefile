proto:
	protoc -I. \
	-I=./third_party/googleapis \
	-I=./third_party/grpc-gateway \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=. --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=./docs --openapiv2_opt logtostderr=true \
	proto/user.proto
