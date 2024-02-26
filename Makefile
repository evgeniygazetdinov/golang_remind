.PHONY: up down logs ps clean

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

ps:
	docker-compose ps

clean:
	docker-compose down -v

init-db:
	docker-compose exec postgres psql -U postgres -d sql_trainer -f /docker-entrypoint-initdb.d/init-db.sql