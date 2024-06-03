CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go


proto-gen:
	./scripts/genproto.sh

DB_URL := "postgres://postgres:123@localhost:5432/dennic_session_service?sslmode=disable"


swag-gen:
	swag init -g ./api/router.go -o api/docs


