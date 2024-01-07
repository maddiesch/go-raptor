CREATE TABLE "key_value" ("key" TEXT NOT NULL, "value" BLOB);

CREATE UNIQUE INDEX "key_value_key_index" ON "key_value" ("key");
