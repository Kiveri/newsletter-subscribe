# подключение файла с переменными окружения
include .env

# получение значений переменных окружения
export

# подстановка текущей рабочей директории в переменную окружения
export PROJECT_ROOT=$(shell pwd)


# запуск окружения
env-up:
	# -d позволяет запускать в фоновом процессе не блокируя терминальную сессию
	@docker compose up -d newsletter-subscribe-postgres

# остановка окружения
env-down:
	@docker compose down newsletter-subscribe-postgres

# очистка окружения
env-cleanup:
	@read -p "Очистить все volume файлы окружения? Очищение окружения приведет к потере всех данных. [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down newsletter-subscribe-postgres && \
	  rm -rf out/pgdata && \
	  echo "Файлы окружения очищены."; \
	else \
	  echo "Очистка окружения отменена."; \
	fi


# создание миграции
migrate-create:
	@if [ -z "$(name)" ]; then \
	    echo "Отсутствует необходимый параметр name. Пример: make migrate-create name=init"; \
	    exit 1; \
	fi; \
	docker compose run --rm newsletter-subscribe-postgres-migrate \
		create \
		-ext sql \
		-dir /scripts/migrations \
		-seq "$(name)"

# применение миграций
migrate-up:
	@make migrate-action action=up

# применение миграций
migrate-down:
	@make migrate-action action=down

# вызов сервиса миграций
migrate-action:
	@if [ -z "$(action)" ]; then \
    	echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
    	exit 1; \
    fi; \
	docker compose run --rm newsletter-subscribe-postgres-migrate \
    	-path /scripts/migrations \
    	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@newsletter-subscribe-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    	"$(action)"
