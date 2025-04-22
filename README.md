# input_output_bound

1. git clone https://github.com/LionJr/input-output-bound
2. docker run -d --name CONTAINER_NAME \
              -p PORT:6379 redis \
               redis-server --requirepass PASSWORD
3. Create .env file with:
    - HTTP_HOST
    - HTTP_PORT
    - REDIS_HOST
    - REDIS_PASSWORD
    - REDIS_DB
4. go run cmd/main.go
