-- to enable RLS, need an role that is not the owner of tables etc, with which to enforce rls
create user ${APPUSER} with password '${APPUSER_PASSWORD}';
grant usage on schema ${MAIN_SCHEMA} to ${APPUSER};
grant usage on schema public to ${APPUSER};

-- restrict to main actions allowed
grant SELECT, INSERT, UPDATE on all tables in schema ${MAIN_SCHEMA} to ${APPUSER};

-- access to current_user_id() and current_user_role() checks
grant execute on all functions in schema public to ${APPUSER};
