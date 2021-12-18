CONFIG_PATH="/Users/tibcsi/Workspace/go/distributed-systems/proglog/certs/.proglog"

export CONFIG_DIR=$(subst $\",,${CONFIG_PATH})

.PHONY: init
init:
	rm -rf ${CONFIG_PATH}
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert -initca test/client-ca-csr.json | cfssljson -bare client-ca

	cfssl gencert \
						-ca=ca.pem \
						-ca-key=ca-key.pem \
						-config=test/ca-config.json \
						-profile=server \
						test/server-csr.json | cfssljson -bare server

	cfssl gencert \
						-ca=client-ca.pem \
						-ca-key=client-ca-key.pem \
						-config=test/ca-config.json \
						-profile=client \
						test/client-csr.json | cfssljson -bare client

	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: compile
compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

.PHONY: test
test:
# MallocNanoZone=0 https://github.com/golang/go/issues/49138
	MallocNanoZone=0 go test -v -race ./...