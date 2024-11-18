FROM postgres:14.13-alpine

# Copy the database initialization script
COPY ./db/postgres_init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
