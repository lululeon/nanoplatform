-- init script
revoke all on database ${DB_NAME} from public;
revoke all on schema public from public;
create schema if not exists ${MAIN_SCHEMA};

-- for gen_random_uuid():
create extension if not exists pgcrypto;

-- change tracking
create or replace function update_updated_at()
  returns trigger as $$
    begin
      NEW.updated_at = now(); 
      return new;
    end;
  $$ language 'plpgsql';


-- user info
create table if not exists ${MAIN_SCHEMA}.users(
  id uuid primary key default gen_random_uuid(),
  username text unique not null,
  avatar_url text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);


drop trigger if exists users_updated_at
  on ${MAIN_SCHEMA}.users;
create trigger users_updated_at before update
  on ${MAIN_SCHEMA}.users for each row execute procedure 
  update_updated_at();


-- projects info
create table if not exists ${MAIN_SCHEMA}.unbuilts(
  id uuid primary key default gen_random_uuid(),
  originator uuid not null,
  title text not null,
  elevator_pitch text,
  logo_url text,
  open_source boolean not null default false,
  tags text[],
  current_status text not null check (current_status in ('unbuilt', 'scaffolding', 'building', 'poc', 'live')),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  foreign key (originator) references ${MAIN_SCHEMA}.users(id) on update cascade on delete cascade
);

drop trigger if exists unbuilts_updated_at
  on ${MAIN_SCHEMA}.unbuilts;
create trigger unbuilts_updated_at before update
  on ${MAIN_SCHEMA}.unbuilts for each row execute procedure 
  update_updated_at();

create index if not exists unbuilts_originator_title_idx on ${MAIN_SCHEMA}.unbuilts(originator, title);
create index if not exists unbuilts_title_idx on ${MAIN_SCHEMA}.unbuilts(title);
