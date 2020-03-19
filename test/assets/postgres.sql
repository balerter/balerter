CREATE TABLE "public"."users" (
    "id" int4 NOT NULL,
    "email" varchar(255) NOT NULL,
    "name" varchar(255) NOT NULL
);

INSERT INTO "public"."users" ("id", "email", "name") VALUES
(1, 'email1@domain.com', 'user1'),
(2, 'email2@domain.com', 'user2'),
(3, 'email3@domain.com', 'user3'),
(4, 'email4@domain.com', 'user4'),
(5, 'email5@domain.com', 'user5')
