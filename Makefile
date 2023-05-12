proto/EchoService.pb.go: proto/EchoService.proto
	protoc $< --go_out=plugins=grpc:.

dist/rootCA.crt:
	mkdir -p dist/
	certificate/generate_ca.sh dist/rootCA

dist/server.crt:
	mkdir -p dist/
	certificate/generate_certificate.sh dist/server dist/rootCA $$(hostname -f)

dist/client.crt:
	mkdir -p dist/
	certificate/generate_certificate.sh dist/client dist/rootCA "${USER}@"$$(hostname -d)

.PHONY: run-client
run-client: proto/EchoService.pb.go dist/rootCA.crt dist/server.crt dist/client.crt
	go run client/client.go \
		--ca_cert="dist/rootCA.crt" \
		--self_cert="dist/client.crt" \
		--self_key="dist/client.key" \
		--server_cert="dist/server.crt"

.PHONY: run-server
run-server: proto/EchoService.pb.go proto/EchoService.pb.go dist/rootCA.crt dist/server.crt
	go run server/server.go \
		--ca_cert="dist/rootCA.crt" \
		--self_cert="dist/server.crt" \
		--self_key="dist/server.key"
