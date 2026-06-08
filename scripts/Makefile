.PHONY: stop start restart rebuild

SERVICES := auth-service menu-service staff-service order-service api-gateway

stop:
	sudo docker stop $(SERVICES)

start:
	sudo docker compose up -d $(SERVICES)

restart: rebuild stop start

rebuild:
	cd ../api-gateway && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server
	cd ../services/auth-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server
	cd ../services/menu-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server
	cd ../services/order-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server
	cd ../services/staff-service && CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server
postgre:
	sudo docker exec -it restaurant-postgres psql -U restaurant_user -d restaurant_db
