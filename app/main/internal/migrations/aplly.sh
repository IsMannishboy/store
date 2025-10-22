#!/bin/bash
set -e
set -a
source /app/config/env.env
set +a 
DATABASE_URL="postgres://${postgres_user}:${postgres_pass}@${postgres_addr}:${postgres_port}/${postgres_db_name}?sslmode=${postgres_sslmode}"
echo "url: $DATABASE_URL"
migrate -path /app/internal/migrations/migrations -database "$DATABASE_URL" up
    
  
