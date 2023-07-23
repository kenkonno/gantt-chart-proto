docker-compose down

docker volume rm dbdata_gantt

docker volume create dbdata_gantt

docker-compose build postgres
