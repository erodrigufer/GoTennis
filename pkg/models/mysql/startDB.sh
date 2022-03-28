#!/bin/sh

mariadb < sessionsTable.mysql && mariadb usersTable.mysql && echo "* DB correctly configured!"
