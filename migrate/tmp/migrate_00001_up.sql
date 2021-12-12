CREATE TABLE "test.users" (
                         "user_id" INT,
                         "name" VARCHAR,
                         "age" INT,
                         "spouse" INT
);

CREATE UNIQUE INDEX "users_user_id" ON "test.users" ("user_id");



CREATE TABLE "test.activities" (
                              "user_id" INT,
                              "date" TIMESTAMP,
                              "name" VARCHAR
);

CREATE INDEX "activities_user_id_date" ON "test.activities" ("user_id", "date");