-- +goose Up
-- +goose StatementBegin
CREATE TABLE "sessions" (
    "id" INTEGER PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,

    "uuid" TEXT NOT NULL,
    "last_ip" TEXT NOT NULL,

    "user_id" INTEGER NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "sessions";
-- +goose StatementEnd
