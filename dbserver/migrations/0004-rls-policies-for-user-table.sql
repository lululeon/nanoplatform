create or replace function current_user_id() returns uuid as $$
  select nullif(current_setting('jwt.claims.user_id', true), '')::uuid;
$$ language sql stable;

create or replace function current_user_role() returns text as $$
  select nullif(current_setting('jwt.claims.role', true), '')::text;
$$ language sql stable;

alter table ${MAIN_SCHEMA}.users enable row level security;

drop policy if exists access_self on ${MAIN_SCHEMA}.users;
create policy access_self on ${MAIN_SCHEMA}.users
  using("id" = current_user_id());

alter table ${MAIN_SCHEMA}.unbuilts enable row level security;

drop policy if exists access_own_unbuilts on ${MAIN_SCHEMA}.unbuilts;
create policy access_own_unbuilts on ${MAIN_SCHEMA}.unbuilts
  using("originator" = current_user_id());
