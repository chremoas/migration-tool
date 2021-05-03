pg_dump -h 10.42.1.30 -d chremoas_aba > chremoas-aba.backup.sql
pg_dump -a -h 10.42.1.30 -d chremoas_aba > chremoas-aba.sql

#psql -h localhost -U postgres -d chremoas_aba -f insert.sql
#psql -h localhost -U postgres -d chremoas_aba -f fix_sequence.sql

COPY public.alliances (alliance_id, alliance_name, alliance_ticker, inserted_dt, updated_dt) FROM stdin;
COPY public.corporations (corporation_id, corporation_name, alliance_id, inserted_dt, updated_dt, corporation_ticker) FROM stdin;
COPY public.characters (character_id, character_name, inserted_dt, updated_dt, corporation_id, token) FROM stdin;
COPY public.authentication_codes (character_id, authentication_code, is_used) FROM stdin;
COPY public.authentication_scopes (authentication_scope_id, authentication_scope_name) FROM stdin;
COPY public.authentication_scope_character_map (character_id, authentication_scope_id) FROM stdin;
COPY public.users (user_id, chat_id) FROM stdin;
COPY public.user_character_map (user_id, character_id) FROM stdin;
