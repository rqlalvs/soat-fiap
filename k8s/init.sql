
SET GLOBAL default_authentication_plugin = 'mysql_native_password';

CREATE DATABASE IF NOT EXISTS soat_fiap CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'appuser'@'%' IDENTIFIED BY 'apppassword';
GRANT ALL PRIVILEGES ON soat_fiap.* TO 'appuser'@'%' WITH GRANT OPTION;

FLUSH PRIVILEGES;
