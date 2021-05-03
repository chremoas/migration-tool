BEGIN;
LOCK TABLE users IN EXCLUSIVE MODE;
SELECT setval('users_user_id_seq', COALESCE((SELECT MAX(user_id)+1 FROM users), 1), false);
COMMIT;
