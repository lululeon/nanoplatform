-- revoke all on database ${DB_NAME} from public;
-- revoke all on schema public from public;
-- create schema if not exists ${MAIN_SCHEMA};
-- create schema if not exists ${ALT_SCHEMA};
select now() as the_time;