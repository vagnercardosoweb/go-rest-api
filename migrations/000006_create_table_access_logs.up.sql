CREATE TABLE access_logs (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT NOT NULL,
    total_failures SMALLINT NOT NULL DEFAULT 0,
    total_success SMALLSERIAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE access_logs
    DROP CONSTRAINT IF EXISTS access_logs_id_pk,
    ADD CONSTRAINT access_logs_id_pk PRIMARY KEY (id);

CREATE INDEX access_logs_id ON access_logs USING btree (id);
CREATE INDEX access_logs_user_id ON access_logs USING btree (user_id);
CREATE INDEX access_logs_total_failures ON access_logs USING btree (total_failures);
CREATE INDEX access_logs_created_at ON access_logs USING btree (created_at);

ALTER TABLE access_logs
    ADD CONSTRAINT access_logs_user_id_pk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;
