-- fix needed because quoting from postgraphile docs:
-- " the keys must contain either one or two period (.) characters, and the prefix (the bit before the first period) must not be used by any Postgres extension.""
-- " Variables without periods will be interpreted as internal Postgres settings, such as role, and will be applied by Postgres.""
-- Therefore, the prior `jwt.claims.role` is reserved for application-specific use, and we need a postgres-specific role pick-up, which is what we need to complete RLS setup.
create or replace function current_user_role() returns text as $$
  select current_user;
$$ language sql stable;

-- decided it'd be nice to have a default / fallback role (think guest auth):
create user ${GUESTUSER} with password '${GUESTUSER_PASSWORD}';
grant usage on schema ${MAIN_SCHEMA} to ${GUESTUSER};
grant usage on schema public to ${GUESTUSER};
grant SELECT on all tables in schema ${MAIN_SCHEMA} to ${GUESTUSER};

-- access to current_user_id() and current_user_role() checks
grant execute on all functions in schema public to ${GUESTUSER};
