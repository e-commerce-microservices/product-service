CREATE TABLE IF NOT EXISTS category (
    "id" serial8 PRIMARY KEY,
    "name" varchar(64) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS product (
    "id" serial8 PRIMARY KEY,
    "name" varchar(128) NOT NULL,
    "description" text NOT NULL,
    "price" numeric(8) NOT NULL,
    "thumbnail" text NOT NULL,
    "inventory" integer NOT NULL DEFAULT 0,
    "supplier_id" serial8 NOT NULL,
    "category_id" serial8 NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE product
ADD
    FOREIGN KEY ( "category_id" ) REFERENCES category ("id");
