-- Widens translation_list.created_at and translation.created_at from date to timestamp,
-- so entries synced from multiple devices on the same day can still be ordered by time,
-- not just date. init.sql already creates new databases with the timestamp column
-- directly; this script brings an existing database (still on the old date column) up
-- to match. Re-running it after upgrading is a no-op since the target types already match.
--
-- Run against the running Postgres container, e.g.:
--   docker exec -i postgres psql -U golang -d production < content-service/migrations/001_translation_created_at_to_timestamp.sql

begin;

alter table translation_list
    alter column created_at type timestamp,
    alter column created_at set default now ();

alter table translation
    alter column created_at type timestamp,
    alter column created_at set default now ();

commit;
