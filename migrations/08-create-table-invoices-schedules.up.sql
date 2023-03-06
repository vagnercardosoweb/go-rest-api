DROP TABLE IF EXISTS invoices_schedules;

DROP TYPE IF EXISTS ENUM_INVOICES_SCHEDULES_STATUS;
CREATE TYPE ENUM_INVOICES_SCHEDULES_STATUS AS ENUM ('paid', 'unpaid');

CREATE TABLE invoices_schedules (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL,
    installment_number SMALLSERIAL NOT NULL,
    status ENUM_INVOICES_SCHEDULES_STATUS NOT NULL DEFAULT 'unpaid',
    paid_at DATE NULL DEFAULT NULL,
    unpaid_at DATE NULL DEFAULT NULL,
    due_date DATE NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
);

ALTER TABLE invoices_schedules
    DROP CONSTRAINT IF EXISTS invoices_schedules_id_pk,
    ADD CONSTRAINT invoices_schedules_id_pk PRIMARY KEY (id);

CREATE INDEX invoices_schedules_id ON invoices_schedules USING btree (id);
CREATE INDEX invoices_schedules_invoice_id ON invoices_schedules USING btree (invoice_id);
CREATE INDEX invoices_schedules_paid_at ON invoices_schedules USING btree (paid_at);
CREATE INDEX invoices_schedules_unpaid_at ON invoices_schedules USING btree (unpaid_at);
CREATE INDEX invoices_schedules_created_at ON invoices_schedules USING btree (created_at);
CREATE INDEX invoices_schedules_status ON invoices_schedules USING btree (status);

ALTER TABLE invoices_schedules
    ADD CONSTRAINT invoices_schedules_invoice_id_fk FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON
        DELETE CASCADE ON
        UPDATE CASCADE;
