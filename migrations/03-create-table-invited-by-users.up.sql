DROP TABLE IF EXISTS invited_by_users;

CREATE TABLE invited_by_users (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    guest_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE invited_by_users
    DROP CONSTRAINT IF EXISTS invited_by_users_id_pk,
    ADD CONSTRAINT invited_by_users_id_pk PRIMARY KEY (id);

ALTER TABLE invited_by_users
    DROP CONSTRAINT IF EXISTS invited_by_users_user_id_guest_id_uk,
    ADD CONSTRAINT invited_by_users_user_id_guest_id_uk UNIQUE (user_id, guest_id);

CREATE INDEX invited_by_users_id ON invited_by_users USING btree (id);
CREATE INDEX invited_by_users_user_id ON invited_by_users USING btree (user_id);
CREATE INDEX invited_by_users_guest_id ON invited_by_users USING btree (guest_id);

ALTER TABLE invited_by_users
    ADD CONSTRAINT invited_by_users_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE invited_by_users
    ADD CONSTRAINT invited_by_users_guest_id_fk FOREIGN KEY (guest_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;
