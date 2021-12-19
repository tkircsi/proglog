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
						-cn="root" \
						test/client-csr.json | cfssljson -bare root-client

	cfssl gencert \
						-ca=client-ca.pem \
						-ca-key=client-ca-key.pem \
						-config=test/ca-config.json \
						-profile=client \
						-cn="nobody" \
						test/client-csr.json | cfssljson -bare nobody-client

	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: compile
compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

$(CONFIG_PATH)/model.conf:
	cp test/model.conf $(CONFIG_PATH)/model.conf

$(CONFIG_PATH)/policy.csv:
	cp test/policy.csv $(CONFIG_PATH)/policy.csv

.PHONY: test
test: $(CONFIG_PATH)/policy.csv $(CONFIG_PATH)/model.conf
# MallocNanoZone=0 https://github.com/golang/go/issues/49138
	MallocNanoZone=0 go test -v -race ./...