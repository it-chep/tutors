LOCAL_BIN := $(CURDIR)/bin

include ./.env

export PATH := $(PATH):$(LOCAL_BIN)

.PHONY: deps
deps:
	# todo: временно
	GOBIN=$(LOCAL_BIN) go install gitlab.ozon.ru/whc/go/libs/xo@v1.0.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: infra
infra:
	docker-compose up -d --build

.PHONY: minfra
minfra-up:
	sleep 2s && \
	$(LOCAL_BIN)/goose postgres "user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} dbname=${DB_NAME} sslmode=disable" -dir=./migrations up

.PHONY: minfra-down
minfra-down:
	$(LOCAL_BIN)/goose postgres "user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} dbname=${DB_NAME} sslmode=disable" -dir=./migrations reset

# If the first argument is "run"...
ifeq (migration,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  MIGRATION_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATION_ARGS):;@:)
endif
.PHONY: migration
migration:
	$(LOCAL_BIN)/goose create -dir=./migrations $(MIGRATION_ARGS) sql


# todo: временно
XO_OUTPUT_PATH=./pkg/xo
XO_TEMPLATE_PATH=./pkg/xo_templates
.PHONY: xo ## генерация dto базы данных
xo:
	rm -r $(XO_OUTPUT_PATH)
	mkdir -p $(XO_OUTPUT_PATH)
	xo "pgsql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" \
	-o $(XO_OUTPUT_PATH) --template-path $(XO_TEMPLATE_PATH) --schema public --suffix ".xo.go" --custom-type-package custom

	rm $(XO_OUTPUT_PATH)/goosedbversion.xo.go
	rm $(XO_OUTPUT_PATH)/xo_db.xo.go