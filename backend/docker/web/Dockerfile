# Use the official Redis image from the Docker Hub
FROM nginx:latest

ADD ./docker/web/default.conf /etc/nginx/conf.d/default.conf
ADD ./docker/web/dist /usr/share/nginx/html

RUN echo "start nginx"
