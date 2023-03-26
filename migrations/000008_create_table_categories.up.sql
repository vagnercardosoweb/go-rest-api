CREATE TABLE categories (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    parent_id UUID NULL DEFAULT NULL,
    name VARCHAR(60) NOT NULL,
    sort_order SMALLSERIAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
);

ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS categories_id_pk,
    ADD CONSTRAINT categories_id_pk PRIMARY KEY (id);

ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS categories_user_id_name_uk,
    ADD CONSTRAINT categories_user_id_name_uk UNIQUE (user_id, name);

CREATE INDEX categories_id ON categories USING btree (id);
CREATE INDEX categories_name ON categories USING btree (name);
CREATE INDEX categories_user_id ON categories USING btree (user_id);
CREATE INDEX categories_parent_id ON categories USING btree (parent_id);
CREATE INDEX categories_sort_order ON categories USING btree (sort_order);
CREATE INDEX categories_created_at ON categories USING btree (created_at);

ALTER TABLE categories
    ADD CONSTRAINT categories_parent_id_fk FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE
        RESTRICT ON UPDATE CASCADE;

ALTER TABLE categories
    ADD CONSTRAINT categories_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE;
