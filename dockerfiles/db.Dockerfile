FROM postgres:14.13-alpine

# Copy the database initialization script
ADD ./db/postgres_init.sql /docker-entrypoint-initdb.d

EXPOSE 5432
