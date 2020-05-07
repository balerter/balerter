CREATE TABLE "public"."users"
(
    "id"      int NOT NULL,
    "name"    varchar,
    "balance" float
);

INSERT INTO "public"."users" ("id", "name", "balance")
VALUES ('1', 'John', '10.2'),
       ('2', 'Mark', '12'),
       ('3', 'Peter', '-15.4');
