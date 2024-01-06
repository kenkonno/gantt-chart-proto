# Use the official Redis image from the Docker Hub
FROM redis:latest

# Expose Redis port
EXPOSE 6379

COPY ./docker/session/redis.conf /usr/local/etc/redis/redis.conf

# Command to start the Redis server
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf"]