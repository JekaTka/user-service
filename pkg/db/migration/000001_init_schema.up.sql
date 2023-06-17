CREATE TABLE "users" (
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("id");
CREATE INDEX ON "users" ("email");