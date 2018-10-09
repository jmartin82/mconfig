[![Build Status](https://travis-ci.org/jmartin82/cnfgo.svg?branch=master)](https://travis-ci.org/jmartin82/cnfgo)
[![codecov](https://codecov.io/gh/jmartin82/cnfgo/branch/master/graph/badge.svg)](https://codecov.io/gh/jmartin82/cnfgo)
# cnfgo

cnfgo is a lightweight Golang library for integrating configs files like (json, yml, toml) and environment variables into one config struct.

## Features
* Load multiple types of files (yaml, json, toml).
* Autofill environment variables with `.env` files.
* Works with nested structs.
* Combine config files and environment variables.
* Map environment variables to struct via tag.
* Auto casting.
* Easy to use (only one function exposed).
* Extendable (allows add more providers).
* Customizable (via facade). 

## Usage

Config is designed to be very simple and straightforward to use. All you can do with it is load configurations to a predifined struct.

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

## Supported types

All environment variables are string but you can specify the correct type in your configuration struct.

The valid types are:

```golang
type TestTypes struct {
	IsInt        int      `env:"FOO_VAR"`
	IsUint       uint     `env:"FOO_VAR"`
	IsString     string   `env:"FOO_VAR"`
	IsFloat      float64  `env:"FOO_VAR"`
	IsBool       bool     `env:"FOO_VAR"`
	PtrIsInt     *int     `env:"FOO_VAR"`
	PtrIsInt8    *int8    `env:"FOO_VAR"`
	PtrIsInt16   *int16   `env:"FOO_VAR"`
	PtrIsInt32   *int32   `env:"FOO_VAR"`
	PtrIsInt64   *int64   `env:"FOO_VAR"`
	PtrIsUint    *uint    `env:"FOO_VAR"`
	PtrIsUint8   *uint8   `env:"FOO_VAR"`
	PtrIsUint16  *uint16  `env:"FOO_VAR"`
	PtrIsUint32  *uint32  `env:"FOO_VAR"`
	PtrIsUint64  *uint64  `env:"FOO_VAR"`
	PtrIsString  *string  `env:"FOO_VAR"`
	PtrIsFloat32 *float32 `env:"FOO_VAR"`
	PtrIsFloat64 *float64 `env:"FOO_VAR"`
    PtrIsBool    *bool    `env:"FOO_VAR"`
    Struct       struct
    PtrStruct    *struct
}
```

In case of an casting error, the field will keep the original value.

## Customization

By default only the Parse function it's exposed. With that, you can read json, yaml, toml. Load .env files and combine the config read from the file with your system environment variables.

But if you don't want the full feature call, you always can use the facade to adjust the functionality to your necessities.

Read only env vars from your system and load into your config struct.

```golang
import 	"github.com/jmartin82/cnfgo/pkg/cnfgo"

configuration := Configuration{}
err := cnfgo.ConfigManager.ReadFromEnvironment(&configuration);
if err != nil {
	panic(err)
}
```

Read only the config from a toml.

```golang
import 	"github.com/jmartin82/cnfgo/pkg/cnfgo"

configuration := Configuration{}
err := cnfgo.ConfigManager.ReadFromFile("config.toml", &configuration)
if err != nil {
	panic(err)
}
```

Use differnt files to load the enviroment variables.

```golang
import 	"github.com/jmartin82/cnfgo/pkg/cnfgo"

configuration := Configuration{}
cnfgo.ConfigManager.SetEnvFiles("enviroment.txt","env.txt")
err := cnfgo.Parse("config.yaml", &configuration)
if err != nil {
	panic(err)
}
```

## Extensibility

You can add more file provider in order to read different kind of config files.

Interface:

```golang
type Unmarshaler interface {
	IsCapable(filename string) bool
	Unmarshal(source []byte, configuration interface{}) error
}
```
Code:

```golang
import 	"github.com/jmartin82/cnfgo/pkg/cnfgo"

configuration := Configuration{}
cnfgo.ConfigManager.AddFileUnmashaler(NewXMLFormatUnmarshaler())
err := cnfgo.Parse("config.yaml", &configuration)
if err != nil {
	panic(err)
}
```

## Current limitations

* Only the exposed fields can be setted via file config or enviroment variable.
* The config struct should be a pointer.
