-- Adds translation.updated_at, needed now that SyncList upserts entries by id instead of doing
-- a blind delete-and-reinsert per sync. Backfills every pre-existing row to updated_at =
-- created_at - the honest neutral default ("no edit is known to have happened since creation"),
-- not a guess at the true last-edit time, since that history doesn't exist: every prior edit
-- went through the old delete-and-reinsert. init.sql already creates new databases with this
-- column directly; this script brings an existing database up to match, same precedent as
-- 001_translation_created_at_to_timestamp.sql. Re-running it after upgrading is a no-op.
--
-- Run against the running Postgres container, e.g.:
--   docker exec -i postgres psql -U golang -d production < content-service/migrations/002_translation_add_updated_at.sql

begin;

alter table translation
    add column if not exists updated_at timestamp not null default now ();

update translation set updated_at = created_at;

commit;
