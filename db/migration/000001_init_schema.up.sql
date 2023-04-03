CREATE TABLE IF NOT EXISTS "posts" (
                         "id" bigserial PRIMARY KEY,
                         "link" varchar NOT NULL unique,
                         "owner" varchar NOT NULL,
                         "likes" bigint NOT NULL,
                         "comments" bigint NOT NULL,
                         "reposts" bigint NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "comments" (
                            "id" bigserial PRIMARY KEY,
                            "post_id" bigint NOT NULL,
                            "owner" varchar NOT NULL,
                            "thread_owner" varchar NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "likes" (
                            "id" bigserial PRIMARY KEY,
                            "post_id" bigint NOT NULL,
                            "owner" varchar NOT NULL
);

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
ALTER TABLE "likes" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
