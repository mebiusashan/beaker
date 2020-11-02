#! /bin/bash

echo "Beaker"
echo "Please input your MySQL HOST:"
read MySQL_Host

echo "Please input your MySQL Port:"
read MySql_Port

echo "Please input your MySql user:"
read MySQL_Name

echo "Please input your MySql password:"
read MySQL_Password

echo "Please input your MySql Database:"
read MySql_Database

echo "Please input your server path:"
read ServerPath

echo "Please input your server user name:"
read UserName

echo "Please input your server user password:"
read UserPassword

echo "Please input your server port:"
read ServerPort

echo "Please input your admin server port:"
read AdminPort

echo "Please input your domain: "
read Domain

mysql -h$MySQL_Host –u$MySQL_Name –p$MySQL_Password –D$MySql_Database<./beaker.sql

cp ./config.toml $ServerPath
cp ./admin.toml $ServerPath
mkdir $ServerPath/keys

# config.toml
sed -i 's/-SITE_URL/'$Domain'/g' $ServerPath/config.toml
sed -i 's/-TEMP_FOLDER/'$ServerPath'/g' $ServerPath/config.toml
sed -i 's/-STATIC_FILE_FOLDER/'$ServerPath'/g' $ServerPath/config.toml
sed -i 's/-PORT/'$ServerPort'/g' $ServerPath/config.toml

sed -i 's/-DB_URL/'$MySQL_Host':'$MySql_Port'/g' $ServerPath/config.toml
sed -i 's/-DB_USER/'$MySQL_Name'/g' $ServerPath/config.toml
sed -i 's/-DB_PW/'$MySQL_Password'/g' $ServerPath/config.toml
sed -i 's/-DB_NAME/'$MySql_Database'/g' $ServerPath/config.toml

# admin.toml
sed -i 's/-Name/'$UserName'/g' $ServerPath/admin.toml
sed -i 's/-Password/'$UserPassword'/g' $ServerPath/admin.toml
sed -i 's/-ConfigPath/'$ServerPath'/g' $ServerPath/admin.toml
sed -i 's/-ServerKey/'$ServerPath'/g' $ServerPath/admin.toml
sed -i 's/-PORT/'$AdminPort'/g' $ServerPath/admin.toml

sed -i 's/-DB_URL/'$MySQL_Host':'$MySql_Port'/g' $ServerPath/admin.toml
sed -i 's/-DB_USER/'$MySQL_Name'/g' $ServerPath/admin.toml
sed -i 's/-DB_PW/'$MySQL_Password'/g' $ServerPath/admin.toml
sed -i 's/-DB_NAME/'$MySql_Database'/g' $ServerPath/admin.toml

