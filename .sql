CREATE DATABASE ad_zh;

use ad_zh;


CREATE TABLE users (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    address varchar(255) NULL,
    telephone varchar(100) NULL,
    full_name  varchar(255) NULL,
    image_url varchar(255) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);