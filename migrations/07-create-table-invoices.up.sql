DROP TABLE IF EXISTS invoices;

DROP TYPE IF EXISTS ENUM_INVOICES_TYPE;
CREATE TYPE ENUM_INVOICES_TYPE AS ENUM ('income', 'expense');

CREATE TABLE invoices (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    category_id UUID NOT NULL,
    type ENUM_INVOICES_TYPE NOT NULL DEFAULT 'income',
    value_in_cents INTEGER NOT NULL DEFAULT 0,
    repeatable_in_days SMALLINT NOT NULL DEFAULT 0,
    total_installments SMALLSERIAL NOT NULL,
    payday SMALLSERIAL NOT NULL,
    start_at DATE NOT NULL,
    end_at DATE NULL DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
);

ALTER TABLE invoices
    DROP CONSTRAINT IF EXISTS invoices_id_pk,
    ADD CONSTRAINT invoices_id_pk PRIMARY KEY (id);

ALTER TABLE invoices
    DROP CONSTRAINT IF EXISTS invoices_payday_ck,
    ADD CONSTRAINT invoices_payday_ck CHECK ( payday >= 1 AND payday <= 31);

CREATE INDEX invoices_id ON invoices USING btree (id);
CREATE INDEX invoices_user_id ON invoices USING btree (user_id);
CREATE INDEX invoices_wallet_id ON invoices USING btree (wallet_id);
CREATE INDEX invoices_category_id ON invoices USING btree (category_id);
CREATE INDEX invoices_type ON invoices USING btree (type);
CREATE INDEX invoices_repeatable_in_days ON invoices USING btree (repeatable_in_days);
CREATE INDEX invoices_start_at ON invoices USING btree (start_at);
CREATE INDEX invoices_created_at ON invoices USING btree (created_at);
CREATE INDEX invoices_end_at ON invoices USING btree (end_at);

ALTER TABLE invoices
    ADD CONSTRAINT invoices_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE invoices
    ADD CONSTRAINT invoices_category_id_fk FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE invoices
    ADD CONSTRAINT invoices_wallet_id_fk FOREIGN KEY (wallet_id) REFERENCES wallets (id) ON DELETE RESTRICT ON UPDATE CASCADE;

COMMENT ON COLUMN invoices.payday IS 'Day of the month that will be paid, used to generate the due_date of the schedules';
COMMENT ON COLUMN invoices.repeatable_in_days IS '0 = never repeats and above 1 repeat based on the days';
