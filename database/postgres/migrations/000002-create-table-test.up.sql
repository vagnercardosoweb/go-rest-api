CREATE TABLE test (
  id uuid primary key not null default generate_uuid_v4(),
  created_at timestamptz not null default now()
)
