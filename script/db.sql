CREATE TABLE if not exists book (
    id     SERIAL PRIMARY KEY,
    title  VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn   VARCHAR(20)  NOT NULL
);

INSERT INTO book (title, author, isbn)
VALUES ('Clean Code', 'Robert C. Martin', '978-0132350884');


drop table book