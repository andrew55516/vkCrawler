CREATE TABLE IF NOT EXISTS "likes_edges_weighted"
(
    "id"         bigserial PRIMARY KEY,
    "source"      bigint     NOT NULL,
    "target"      bigint     NOT NULL,
    "weight"      bigint     NOT NULL,
    "start" date NOT NULL DEFAULT (now()),
    "end"   date NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "comments_edges_weighted"
(
    "id"         bigserial PRIMARY KEY,
    "source"      bigint     NOT NULL,
    "target"      bigint     NOT NULL,
    "weight"      bigint     NOT NULL,
    "start" date NOT NULL DEFAULT (now()),
    "end"   date NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "all_edges_weighted"
(
    "id"         bigserial PRIMARY KEY,
    "source"      bigint     NOT NULL,
    "target"      bigint     NOT NULL,
    "weight"      bigint     NOT NULL,
    "start" date NOT NULL DEFAULT (now()),
    "end"   date NOT NULL DEFAULT (now())
);

ALTER TABLE "likes_edges_weighted" ADD FOREIGN KEY ("source") REFERENCES "likes_nodes" ("id");
ALTER TABLE "likes_edges_weighted" ADD FOREIGN KEY ("target") REFERENCES "likes_nodes" ("id");
ALTER TABLE "comments_edges_weighted" ADD FOREIGN KEY ("source") REFERENCES "comments_nodes" ("id");
ALTER TABLE "comments_edges_weighted" ADD FOREIGN KEY ("target") REFERENCES "comments_nodes" ("id");
ALTER TABLE "all_edges_weighted" ADD FOREIGN KEY ("source") REFERENCES "all_nodes" ("id");
ALTER TABLE "all_edges_weighted" ADD FOREIGN KEY ("target") REFERENCES "all_nodes" ("id");