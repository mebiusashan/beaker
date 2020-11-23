# Beaker

Beaker is a simple blog system.

## Version

0.1.1

## Platforms

- Linux
- MacOS

## Install

### Build with source

You need install golang with v1.5.0.

Run command in terminal.

```
git clone https://github.com/mebiusashan/beaker.git
cd beake
./build.sh
```

Now you can find `bin` folder in current folder.If your server's os is Linux, please run command `cd linux`. If your server's os is MacOS, run command `cd darwin`.

You can see 2 executable files, `beaker_server` and `beaker_admin`. So you build system is success.

### Install Server

#### Create Database with MySQL

```
mysql>CREATE DATEBASE beaker;
```

#### Install beaker

You need have install folder in you server, like `/beaker`. If you want install beaker to `/www/blog`, Run command in terminal:

```
cd /beaker
./install.sh
Beaker
Please input your MySQL HOST:
localhost
Please input your MySQL Port:
3306
Please input your MySql user:
root
Please input your MySql password:
root
Please input your MySql Database:
beaker
Please input your server path:
/www/blog
Please input your server user name:
beaker
Please input your server user password:
beaker123
Please input your server port:
9091
Please input your admin server port:
9092
Please input your domain: 
ashan.org
```

Now you install beaker is success, you can see `config.toml` file in `/www/blog/` and `admin/toml` file.

`config.toml` like this:

```
#config file

[website]
SITE_NAME = "Beaker"
SITE_URL = "ashan.org"
SITE_DES = "a simple blog system"
SITE_FOOTER = "Beaker is a simple blog system. [github.com/mebiusashan/beaker]"
SITE_KEYWORDS = "beaker, golang, blog"

INDEX_LIST_NUM = 30
TWEET_NUM_ONE_PAGE = 10
TEMP_FOLDER = "/www/blog/temp"
STATIC_FILE_FOLDER = "/www/blog/static/"

[server]
PORT = ":9091"
URL = "localhost"

[redis]
REDIS_IP = "127.0.0.1"
REDIS_PORT = "6379"
REDIS_PREFIX = "beaker_"
EXPIRE_TIME = 25200

[database]
DB_URL = "localhost:3306"
DB_USER = "root"
DB_PW = "root"
DB_NAME = "beaker"
MAX_IDLE_NUM = 10
MAX_OPEN_NUM = 100
```

`admin.toml`

```
[authinfo]
Name="beaker"
Password="beaker123"
ConfigPath="/www/blog/config.toml"
ServerKeyDir="/www/blog/keys/"
ClientKeyDir="/www/blog/keys/"

[server]
PORT = ":9092"
URL = "localhost"

[redis]
REDIS_IP = "127.0.0.1"
REDIS_PORT = "6379"
REDIS_PREFIX = "beaker_"
EXPIRE_TIME = 25200

[database]
DB_URL = "localhost:3306"
DB_USER = "root"
DB_PW = "root"
DB_NAME = "beaker"
MAX_IDLE_NUM = 10
MAX_OPEN_NUM = 100
```

You can change it.

### PM2 config

```
module.exports = {
  apps: [{
    name: 'blog web server',
    cwd: "/www/blog",
    script: './beaker_server',
    env: {
	"BEAKERPATH":"/www/blog/config.toml"
    },
    watch: true
  }, {
    name: 'blog admin server',
    cwd: "/www/blog",
    script: './beaker_admin',
    env: {
	"BEAKERADMINPATH":"/www/blog/admin.toml"
    },
    watch: true
  }]
}

```


## How to use

### Run server

You need installed beaker. If you beaker in `/www/blog`, run command in terminal:

```
export BEAKERPATH=/www/blog/config.toml
export BEAKERADMINPATH=/www/blog/admin/toml
./beaker_server
./beaker_admin
```

Open `localhost:9091` on brower.

### Run client cli

Copy beaker cli file to your computer. Beaker cli support platform with Linux, MacOS and Windows.

If your computer os is MacOS, run command in terminal:

```
mv beaker_mac beaker
./beaker
```

You need set your website info, list this:

```
./beaker config addw http://xxx.com:8888 -uroot -proot -ablog -d
```

## Contributing

Raising a good question is the first step to participate a open source community. You can report issues [here](https://github.com/mebiusashan/beaker/issues). 

## License

This content is released under the (https://github.com/mebiusashan/beaker/blob/master/LICENSE) AGPLv3 License.

![AGPLv3](https://img.shields.io/badge/license-AGPLv3-blue.svg)
