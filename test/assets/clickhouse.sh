#!/bin/bash
set -e

clickhouse client -n <<-EOSQL
    CREATE TABLE default.users (date Date, name String, age Int32) ENGINE = MergeTree(date, (date), 8192);
    INSERT INTO users  (date, name, age) VALUES ('2020-01-01', 'user1', 42), ('2020-01-01', 'user2', 42), ('2020-01-01', 'user3', 42)
EOSQL