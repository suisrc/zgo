.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = zgo
SERVER_BIN = ./cmd/app/${APP}
RELEASE_ROOT = release
RELEASE_SERVER = release/${APP}

run: start

dev: debug

# 初始化mod
mod:
	go mod init github.com/suisrc/${APP}

install:
	go install ./cmd/app

# 修正依赖
tidy:
	go mod tidy

build:
	go build -ldflags "-w -s" -o $(SERVER_BIN) ./cmd/app

# go get -u github.com/go-delve/delve/cmd/dlv
# 使用“--”将标志传递给正在调试的程
# 调试过程中,会产生一些僵尸进程,这个时候,可以通过杀死父进程解决
# ps -ef | grep defunct | more (kill -9 pid 是无法删除进程)
debug:
	dlv debug --headless --api-version=2 --listen=127.0.0.1:2345 cmd/app/main.go -- web -c ./configs/config.toml
	#echo c | dlv debug --accept-multiclient --api-version=2 --listen=127.0.0.1:2345 cmd/app/main.go -- web -c ./configs/config.toml

start:
	go run cmd/app/main.go web -c ./configs/__config.toml

# go get -u github.com/swaggo/swag/cmd/swag
swagger:
	swag init --generalInfo ./app/swagger.go --output ./app/swagger

# go get -u github.com/google/wire/cmd/wire
wire:
	wire gen ./app/injector

# go get github.com/facebookincubator/ent/cmd/entc
# cd ./app/model && entc init User
# entc init --target ./app/model/ent/schema
# generate the schema for User under <project>/ent/schema/
entc:
	go generate ./app/model/ent

# go get -u github.com/nicksnyder/go-i18n/v2/goi18n
# goi18n -help
i18n:
	goi18n extract -outdir locales -sourceLanguage zh-CN
i18n-update:
	cd locales && goi18n merge -sourceLanguage zh-CN active.zh-CN.toml translate.en-US.toml
i18n-merge:
	cd locales && goi18n merge -sourceLanguage en-US active.en-US.toml translate.en-US.toml

# go get -u github.com/mdempsky/gocode
code:
	gocode -s -debug

test:
	@go test -v $(shell go list ./...)

clean:
	rm -rf data release $(SERVER_BIN) ./app/test/data ./cmd/app/data

pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}

cli:
	go run cmd/cli.go init

# demo
demo-run:
	go run cmd/demo/main.go demo -c ./configs/config.toml
demo-swagger:
	swag init --generalInfo ./demo/swagger.go --output ./demo/swagger
demo-wire:
	wire gen ./demo/injector
demo-entc:
	go generate ./demo/model/ent

# database
db-mysql:
	go run cmd/db/main.go mysql -m doc/model.md -o configs/model.mysql.sql
db-entc:
	go run cmd/db/main.go entc -m doc/model.md -o app/model/ent/schema
db-entc-del:
	go run cmd/db/main.go entc-del -o app/model/ent -es $(es)