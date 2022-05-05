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
