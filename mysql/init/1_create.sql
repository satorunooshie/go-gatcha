CREATE DATABASE golang_db;
use golang_db;
CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    token VARCHAR(200) NOT NULL,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp on update current_timestamp
);

CREATE TABLE IF NOT EXISTS rarities (
    id int NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    probability float NOT NULL
);

CREATE TABLE IF NOT EXISTS chara (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rarity_id int,
    FOREIGN KEY (rarity_id) REFERENCES rarities(id) ON DELETE CASCADE
);

INSERT INTO
    users (name, token)
VALUES
("username", "secret");

INSERT INTO
    rarities
VALUES
(1, "SSR", 0.01),
(2, "SR", 0.1),
(3, "R", 0.89);

INSERT INTO
    chara (name, rarity_id)
VALUES
("super1", 1),
("super2", 1),
("frequest1", 2),
("frequent2", 2),
("frequent3", 2),
("frequent4", 2),
("common1", 3),
("common2", 3),
("common3", 3),
("common4", 3),
("common5", 3),
("common6", 3),
("common7", 3),
("common8", 3);
