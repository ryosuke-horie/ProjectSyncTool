.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen@v0.23.1
	tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

.PHONY: deploy
deploy:
	wrangler deploy

.PHONY: init-db
init-db:
	wrangler d1 execute d1-project-management-homemade-db --file=./schema.sql --remote

.PHONY: init-db-local
init-db-local:
	wrangler d1 execute d1-project-management-homemade-db --file=./schema.sql --local
