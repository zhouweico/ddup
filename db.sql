CREATE TABLE public.users (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NULL,
    updated_at datetime(3) DEFAULT NULL,
    deleted_at datetime(3) DEFAULT NULL,
    username varchar(50) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(100) NOT NULL,
    nickname varchar(50) DEFAULT NULL,
    avatar varchar(255) DEFAULT NULL,
    bio varchar(500) DEFAULT NULL,
    status tinyint NOT NULL DEFAULT '1',
    last_login datetime DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_username (username),
    UNIQUE KEY idx_email (email),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;