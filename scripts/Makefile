.PHONY: all build stop start restart rebuild logs db proto \
        build-gateway build-auth build-menu build-schedule build-table \
        build-order build-notification build-user \
        fe-customer fe-admin fe-kitchen

# ── Docker service names (matches docker-compose.yml) ─────────────────────────
APP_SERVICES := auth-service menu-service schedule-service table-service \
                order-service notification-service user-service api-gateway

# ── Full stack ─────────────────────────────────────────────────────────────────
start:
	cd .. && sudo docker compose up -d $(APP_SERVICES)

stop:
	cd .. && sudo docker compose stop $(APP_SERVICES)

restart: rebuild stop start

logs:
	cd .. && sudo docker compose logs -f $(APP_SERVICES)

# ── Compile all Go binaries (Linux static) ────────────────────────────────────
rebuild: build-gateway build-auth build-menu build-schedule build-table \
         build-order build-notification build-user

build-gateway:
	cd ../api-gateway && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-auth:
	cd ../services/auth-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-menu:
	cd ../services/menu-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-schedule:
	cd ../services/schedule-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-table:
	cd ../services/table-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-order:
	cd ../services/order-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-notification:
	cd ../services/notification-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

build-user:
	cd ../services/user-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

# ── Targeted restart (rebuild single service + docker restart) ─────────────────
restart-gateway: build-gateway
	cd .. && sudo docker compose restart api-gateway

restart-auth: build-auth
	cd .. && sudo docker compose restart auth-service

restart-menu: build-menu
	cd .. && sudo docker compose restart menu-service

restart-schedule: build-schedule
	cd .. && sudo docker compose restart schedule-service

restart-table: build-table
	cd .. && sudo docker compose restart table-service

restart-order: build-order
	cd .. && sudo docker compose restart order-service

restart-notification: build-notification
	cd .. && sudo docker compose restart notification-service

restart-user: build-user
	cd .. && sudo docker compose restart user-service

# ── Proto generation ──────────────────────────────────────────────────────────
proto:
	export PATH=$$PATH:~/go/bin && bash generate-proto.sh

# ── Database access ───────────────────────────────────────────────────────────
db:
	sudo docker exec -it restaurant-postgres psql -U restaurant_user -d restaurant_db

# ── Frontend dev servers ──────────────────────────────────────────────────────
fe-customer:
	cd ../restaurant-app && npm run dev

fe-admin:
	cd ../restaurant-app-admin && npm run dev

fe-kitchen:
	cd ../restaurant-app-kitchen && npm run dev
