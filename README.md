[![Build Status](https://travis-ci.org/jmartin82/cnfgo.svg?branch=master)](https://travis-ci.org/jmartin82/cnfgo)
[![codecov](https://codecov.io/gh/jmartin82/cnfgo/branch/master/graph/badge.svg)](https://codecov.io/gh/jmartin82/cnfgo)
# cnfgo

cnfgo is a lightweight Golang library for integrating configs files like (json, yml, toml) and environment variables into one config struct.

## Usage

First define a configuration structure:

```golang
type MysqlConfiguration struct {
    Host     string `env:"MYSQL_HOST"`
    Username string `env:"MYSQL_USERNAME"`
    Password string `env:"MYSQL_PASSWORD"`
    Database string `env:"MYSQL_DATABASE"`
    Port     int    `env:"MYSQL_PORT"`
}

type RedisConfiguration struct {
    Host string `env:"REDIS_HOST"`
    Port int    `env:"REDIS_PORT"`
}

type Configuration struct {
    Port         int `env:"APP_PORT"`
    Mysql        MysqlConfiguration
    Redis        RedisConfiguration
}
```

Then fill your YAML (config.yaml) file:

```yaml
---
Port: 3001
Mysql:
  Host: 192.168.0.1
  Username: root
  Password: test
  Database: cnfgo
  Port: 3306
Redis:
  Host: localhost
  Port: 6379
```

From your code:

```golang
import 	"github.com/jmartin82/cnfgo/pkg/cnfgo"

configuration := Configuration{}
err := cnfgo.Parse("config.yaml", &configuration)
if err != nil {
	panic(err)
}
```

Finally in prod:

```bash
MYSQL_PASSWORD=SuperSecret MYSQL_HOST=mysql.prod.service REDIS_HOST=mysql.prod.service your_app
```
