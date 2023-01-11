log2web
=======

```
log2web [-p PORTNO] LOGPATH
```

- Parameter
    - PORTNO: default is 8000
    - LOGPATH: the path showed to web
- Feature
    - Showing the last 1024 bytes of the logfile to web via http
    - When users type F5, reload the log. (Not watching filesystem)

Install
-------

Download binaries from [Releases](https://github.com/hymkor/log2web/releases) and unzip the executable.

If you have scoop-install,
```
scoop install https://raw.githubusercontent.com/hymkor/log2web/master/log2web.json
```
OR
```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install log2web
```
