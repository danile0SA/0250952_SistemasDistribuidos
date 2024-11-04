# Configuración de la ruta de destino de los archivos de configuración
CONFIG_PATH=C:\Users\danie\Documents\UP Daniel\Computo Distribuido\Go_Server\GO_Module\0250952_SistemasDistribuidos\test

.PHONY: init

# Inicializar carpeta de configuración
init:
	if (!(Test-Path "$(CONFIG_PATH)")) { New-Item -ItemType Directory -Force -Path "$(CONFIG_PATH)" }

.PHONY: gencert
# gencert - Genera certificados
gencert: init
	cfssl gencert -initca test\ca-csr.json | cfssljson -bare ca

	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test\ca-config.json -profile=server test\server-csr.json | cfssljson -bare server

	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test\ca-config.json -profile=client test\client-csr.json | cfssljson -bare client

	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test\ca-config.json -profile=client -cn="root" test\client-csr.json | cfssljson -bare root-client

	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test\ca-config.json -profile=client -cn="nobody" test\client-csr.json | cfssljson -bare nobody-client

	Move-Item -Path *.pem, *.csr -Destination "$(CONFIG_PATH)" -Force

compile:
	protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.

$(CONFIG_PATH)\model.conf:
	Copy-Item -Path test\model.conf -Destination "$(CONFIG_PATH)\model.conf" -Force

$(CONFIG_PATH)\policy.csv:
	Copy-Item -Path test\policy.csv -Destination "$(CONFIG_PATH)\policy.csv" -Force

test: $(CONFIG_PATH)\policy.csv $(CONFIG_PATH)\model.conf
	go test -race ./...

compile_rpc:
	protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=.
