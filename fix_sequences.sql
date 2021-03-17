BEGIN;
LOCK TABLE users IN EXCLUSIVE MODE;
SELECT setval('users_user_id_seq', COALESCE((SELECT MAX(user_id)+1 FROM users), 1), false);
COMMIT;

BEGIN;
LOCK TABLE authentication_scopes IN EXCLUSIVE MODE;
SELECT setval('authentication_scopes_authentication_scope_id_seq', COALESCE((SELECT MAX(authentication_scope_id)+1 FROM authentication_scopes), 1), false);
COMMIT;
