CREATE TABLE "activities" (
                              "user_id" INT,
                              "date" TIMESTAMP,
                              "name" VARCHAR
) PARTITION BY RANGE("date");

CREATE INDEX "activities_user_id_date" ON "activities" ("user_id", "date");

CREATE TABLE "activities_202111" PARTITION OF "activities" FOR VALUES FROM
    ('2021-11-01'::TIMESTAMP) TO ('2021-12-01'::TIMESTAMP);

CREATE TABLE "activities_202112" PARTITION OF "activities" FOR VALUES FROM
    ('2021-12-01'::TIMESTAMP) TO ('2022-01-01'::TIMESTAMP);