#!/bin/bash
cd /home/isucon/webapp/go/
cat prepare_script/update.sql | mysql -uisucon isucon
prepare_script/prepare_script

rm /tmp/mysql-slow.log
mysqladmin -uroot -proot flush-logs
