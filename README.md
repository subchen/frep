# go-template-replace
A template replace tools written golang

# Command-line Options

Output console using ENV

```
export IPADDR=0.0.0.0
export PORT=8080
$ gtr -t nginx.conf.in --stdout
```

Output console using arguments

```
$ gtr -t nginx.conf.in --stdout -e IPADDR=0.0.0.0 -e PORT=8080
```

Output to default file

```
$ gtr -t nginx.conf.in --overwrite -e IPADDR=0.0.0.0 -e PORT=8080
```

If your file uses `{{` and `}}` as part of it's syntax, you can change the template escape characters using the -delims.

```
$ gtr --delims "<%:%>" ...
```
