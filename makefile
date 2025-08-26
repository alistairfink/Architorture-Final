.PHONY: db
db:
	docker run -d --name architorture-db -e POSTGRES_PASSWORD=replace_with_password -v "$(CURDIR)/Backend/DatabaseImporter/init.sql:/docker-entrypoint-initdb.d/init.sql" \
	-e PGDATA=/var/lib/postgresql/data/pgdata -p 5430:5432 postgres

.PHONY: backend
backend:
	cd Backend && go run main.go

.PHONY: frontend
frontend:
	docker run -it --rm -v "$(CURDIR)/Frontend:/app" -w /app -p 3000:3000 node:10.24.1-alpine3.11 sh -c "npm install && npm start"