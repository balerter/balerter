CREATE TABLE default.users
(
    id Int64,
    name String,
    balance Float64
)
ENGINE = MergeTree() PARTITION BY (id) ORDER BY (id);

INSERT INTO default.users (id, name, balance) VALUES (1, 'John', 10.10);
INSERT INTO default.users (id, name, balance) VALUES (2, 'Bill', -10.10);
INSERT INTO default.users (id, name, balance) VALUES (3, 'Mark', 0);