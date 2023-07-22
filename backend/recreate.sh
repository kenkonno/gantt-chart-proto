docker-compose down

docker volume rm dbdata

docker volume create dbdata

docker-compose build postgres
