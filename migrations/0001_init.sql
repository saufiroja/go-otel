CREATE USER root WITH PASSWORD 'root';

DROP DATABASE IF EXISTS accountdb;
CREATE DATABASE accountdb;

GRANT ALL PRIVILEGES ON DATABASE accountdb TO root;