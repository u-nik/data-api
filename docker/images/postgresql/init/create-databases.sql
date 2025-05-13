-- Create hydra user
CREATE USER hydra
  WITH PASSWORD 'password'
       SUPERUSER;

-- Create database for hydra
CREATE DATABASE hydra
  WITH OWNER = hydra
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'en_US.utf8'
       LC_CTYPE = 'en_US.utf8'
       CONNECTION LIMIT = -1;
