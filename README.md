[![Build Status](https://travis-ci.org/subchen/frep.svg?branch=master)](https://travis-ci.org/subchen/frep)
[![License](http://img.shields.io/badge/License-Apache_2-red.svg?style=flat)](http://www.apache.org/licenses/LICENSE-2.0)


# frep

A template file replace tool written golang

```
Usage: frep [OPTIONS] input-file:[output-file] ...
   or: frep [ --version | --help ]

Transform template file using environment, arguments, json/yaml files

Options:
  -e, --env=[]      set variable name=value, can be passed multiple times
  --json=string     load variables from json object
  --load=[]         load variables from json/yaml files
  --overwrite       overwrite if destination file exists
  --testing         test mode, output transform result to console
  --delims={{:}}    template tag delimiters
  --version         show version information
  --help            show this help
```

# Downloads

[frep-1.1.0](https://github.com/subchen/frep/releases/tag/v1.1.0)

Linux amd64

```
curl -fSL https://github.com/subchen/frep/releases/download/v1.1.0/frep-linux-amd64.zip -o frep.zip
sudo unzip frep.zip -d /usr/bin/
```

MAC OS

```
curl -fSL https://github.com/subchen/frep/releases/download/v1.1.0/frep-darwin-amd64.zip -o frep.zip
sudo unzip frep.zip -d /usr/bin/
```

Windows

```
curl -fSL https://github.com/subchen/frep/releases/download/v1.1.0/frep-window-amd64.zip -o frep.zip
unzip frep.zip -d c:/windows/
```

# Examples

## Load template variables

Load from environment

```
export webroot=/usr/share/nginx/html
export port=8080
frep nginx.conf.in
```

Load from arguments

```
frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080
```

Load from JSON String

```
frep nginx.conf.in --json '{"webroot": "/usr/share/nginx/html", "port": 8080}'
```

Load from JSON file

```
cat > ctx.json << EOF
{
  "webroot": "/usr/share/nginx/html",
  "port": 8080,
  "servers": [
    "127.0.0.1:8081",
    "127.0.0.1:8082"
  ]
}
EOF

frep nginx.conf.in --load ctx.json
```

Load from Yaml file

```
cat > ctx.yaml << EOF
webroot: /usr/share/nginx/html
port: 8080
servers:
  - 127.0.0.1:8081
  - 127.0.0.1:8082
EOF

frep nginx.conf.in --load ctx.yaml
```

## Output

Output to default file (auto remove last file ext)

```
frep nginx.conf.in --overwrite
```

Output to specified file

```
frep nginx.conf.in:/etc/nginx.conf --overwrite -e port=8080
```

Output to console

```
frep nginx.conf.in --testing
```

Output multiple files

```
frep nginx.conf.in redis.conf.in ...
```

## Others

If your file uses `{{` and `}}` as part of it's syntax, you can change the template escape characters using the -delims.

```
frep --delims "<%:%>" ...
```

# Template file

Templates use Golang [text/template](http://golang.org/pkg/text/template/). You can access environment variables within a template

```
ENV.PATH = {{ .PATH }}
```

There are some built-in functions as well: https://github.com/Masterminds/sprig

nginx.conf.in

```
server {
    listen {{.port}} default_server;

    root {{.webroot | default "/usr/share/nginx/html"}};
    index index.html index.htm;

    location /api {
        access_log off;
        proxy_pass http://backend;
    }
}

upstream backend {
    ip_hash;
{{range .servers}}
    server {{.}};
{{end}}
}
```

