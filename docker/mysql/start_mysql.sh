#!/usr/bin/env bash
docker run --name=mysql1 -d -e MYSQL_ROOT_PASSWORD=workhard -p 3306:3306 mysql/mysql-server