CREATE TABLE snippets(
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title varchar(100) NOT NULL ,
    content TEXT NOT NULL ,
    created DATETIME NOT NULL ,
    expires DATETIME NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(id);

CREATE TABLE users(
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL ,
    email VARCHAR(255) NOT NULL ,
    hashed_password CHAR(60) NOT NULL ,
    created DATETIME NOT NULL
);
ALTER TABLE  users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO  users(id, name, email, hashed_password, created)
VALUES (
        -- 指定首个插入的id(适配测试的逻辑)
        39,
        'Miku',
        'miku@vocaloid.com',
        '$2a$12$xP/sNUJvHoFrpezJXSv4m.Oy6.bOpq5n4JBGGcql9g7/N9KHsEH9a',
        '2007-8-31 00:00:00'
);