[![Build Status](https://travis-ci.org/subchen/frep.svg?branch=master)](https://travis-ci.org/subchen/frep)
[![License](http://img.shields.io/badge/License-Apache_2-red.svg?style=flat)](http://www.apache.org/licenses/LICENSE-2.0)


# frep

Generate file using template from environment, arguments, json/yaml/toml config files.

```
NAME:
   frep - Generate file using template

USAGE:
   frep [options] input-file[:output-file] ...

VERSION:
   1.3.x

AUTHORS:
   Guoqiang Chen <subchen@gmail.com>

OPTIONS:
   -e, --env name=value    set variable name=value, can be passed multiple times
       --json jsonstring   load variables from json object string
       --load file         load variables from json/yaml/toml file
       --no-sys-env        exclude system environments, default false
       --overwrite         overwrite if destination file exists
       --dryrun            just output result to console instead of file
       --strict            exit on any error during template processing
       --delims value      template tag delimiters (default: {{:}})
       --help              print this usage
       --version           print version information

EXAMPLES:
   frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080
   frep nginx.conf.in:/etc/nginx.conf -e webroot=/usr/share/nginx/html -e port=8080
   frep nginx.conf.in --json '{"webroot": "/usr/share/nginx/html", "port": 8080}'
   frep nginx.conf.in --load config.json --overwrite
   echo "{{ .Env.PATH }}"  | frep -
```

## Downloads

v1.3.7 Release: https://github.com/subchen/frep/releases/tag/v1.3.7

- Linux

    ```
    curl -fSL https://github.com/subchen/frep/releases/download/v1.3.7/frep-1.3.7-linux-amd64 -o /usr/local/bin/frep
    chmod +x /usr/local/bin/frep
    
    # centos / redhat
    yum install https://github.com/subchen/frep/releases/download/v1.3.7/frep-1.3.7-68.x86_64.rpm
    
    # ubuntu
    curl -fSL https://github.com/subchen/frep/releases/download/v1.3.7/frep_1.3.7-68_amd64.deb -o frep_1.3.7-68_amd64.deb
    dpkg -i frep_1.3.7-68_amd64.deb
    ```

- macOS

    ```
    brew install subchen/tap/frep
    ```

- Windows

    ```
    wget https://github.com/subchen/frep/releases/download/v1.3.7/frep-1.3.7-windows-amd64.exe
    ```

## Docker

You can run frep using docker container

```
docker run -it --rm subchen/frep --help
```


## Examples

### Load template variables

- Load from environment

    ```
    export webroot=/usr/share/nginx/html
    export port=8080
    frep nginx.conf.in
    ```

- Load from arguments

    ```
    frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080
    ```

- Load from JSON String

    ```
    frep nginx.conf.in --json '{"webroot": "/usr/share/nginx/html", "port": 8080}'
    ```

- Load from JSON file

    ```
    cat > config.json << EOF
    {
      "webroot": "/usr/share/nginx/html",
      "port": 8080,
      "servers": [
        "127.0.0.1:8081",
        "127.0.0.1:8082"
      ]
    }
    EOF

    frep nginx.conf.in --load config.json
    ```

- Load from YAML file

    ```
    cat > config.yaml << EOF
    webroot: /usr/share/nginx/html
    port: 8080
    servers:
      - 127.0.0.1:8081
      - 127.0.0.1:8082
    EOF

    frep nginx.conf.in --load config.yaml
    ```

- Load from TOML file

    ```
    cat > config.toml << EOF
    webroot = /usr/share/nginx/html
    port = 8080
    servers = [
       "127.0.0.1:8081",
       "127.0.0.1:8082"
    ]
    EOF

    frep nginx.conf.in --load config.toml
    ```

### Input/Output

- Input from file

    ```
    // input file: nginx.conf
    frep nginx.conf.in
    ```

- Input from console(stdin)

    ```
    // input from stdin pipe
    echo "{{ .Env.PATH }}" | frep -
    ```

- Output to default file (Removed last file ext)

    ```
    // output file: nginx.conf
    frep nginx.conf.in --overwrite
    ```

- Output to the specified file

    ```
    // output file: /etc/nginx.conf
    frep nginx.conf.in:/etc/nginx.conf --overwrite -e port=8080
    ```

- Output to console(stdout)

    ```
    frep nginx.conf.in --dryrun
    frep nginx.conf.in:-
    ```

- Output multiple files

    ```
    frep nginx.conf.in redis.conf.in ...
    ```

## Template

Templates use Golang [text/template](http://golang.org/pkg/text/template/).

You can access environment variables within a template

```
Env.PATH = {{ .Env.PATH }}
```

If your template file uses `{{` and `}}` as part of it's syntax,
you can change the template escape characters using the `--delims`.

```
frep --delims "<%:%>" ...
```

There are some built-in functions as well: Masterminds/sprig v2.14.1
- github: https://github.com/Masterminds/sprig
- doc: http://masterminds.github.io/sprig/

More [funcs](https://github.com/subchen/frep/blob/master/func.go) added:
- toJson
- toYaml
- toToml
- toBool
- fileSize
- fileLastModified
- fileGetBytes
- fileGetString
- fileExists
- include
- countRune
- pipeline compatible regex functions from sprig 
    - reReplaceAll
    - reReplaceAllLiteral
    - reSplit

Sample of nginx.conf.in

```
server {
    listen {{ .port }} default_server;

    root {{ .webroot | default "/usr/share/nginx/html" }};
    index index.html index.htm;

    location /api {
        {{ include "shared/log.nginx" | indent 8 | trim }}
        proxy_pass http://backend;
    }
}

upstream backend {
    ip_hash;
{{- range .servers }}
    server {{.}};
{{- end }}
}
```
