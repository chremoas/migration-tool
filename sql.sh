TABLES="users alliances corporations characters authentication_codes authentication_scope_character_map authentication_scopes user_character_map"

rm -f insert.sql

for table in ${TABLES}; do
	echo "TRUNCATE TABLE ${table} CASCADE;" >> insert.sql
done

for table in ${TABLES}; do
	echo "Processing ${table}"
	./cockroach sql --insecure --host 10.42.2.27 -e "SELECT * from \"chremoas-aba\".${table};" --format=csv > chremoas-aba/${table}.csv
	echo "\\\copy ${table} FROM 'chremoas-aba/${table}.csv' CSV HEADER;" >> insert.sql
done

sed -i 's/,$/,""/g' chremoas-aba/characters.csv
sed -i 's/NULL//g' chremoas-aba/corporations.csv

#psql -h localhost -U postgres -d chremoas_aba -f insert.sql
#psql -h localhost -U postgres -d chremoas_aba -f fix_sequence.sql
