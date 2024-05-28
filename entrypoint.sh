#!/bin/bash

docker-entrypoint.sh mysqld &

until mysql -uroot -p"$MYSQL_ROOT_PASSWORD" -e 'USE psy_data;' &> /dev/null
do
	echo "waiting.."
	sleep 5
done

echo "schema initialized"
exit 0
