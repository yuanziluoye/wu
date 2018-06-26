# wu (呜~)

A minimal **W**atch **U**tility who can run and restart specified command in
response to file changes automatically.

This utility is intended to provide a tiny tool to automate the Edit-Build-Run
loop of development. Although it is quite similar to watch tasks of Grunt or Gulp,
`wu` is designed to be just a single command with simplest interfaces to work with.

# Install

To install `wu` from source code, you have to install Golang's tool chain first.
Then run:

```
go get github.com/yuanziluoye/wu
go install github.com/yuanziluoye/wu
```

Precompiled version can be found [here](https://github.com/shanzi/wu/releases).

# Edit config.yaml
```
worker:
    -
        Directory: .
        Patterns: ['*.go', '*.js']
        Command: ['go','build']
    -
        Directory: .
        Patterns: ['*.c']
        Command: ['/bin/sh', '-c', echo "START" && sleep 5 && echo "END"]

events: ['CREATE', 'REMOVE', 'WRITE', 'RENAME']

logPath: ./log/app.log
```

# Usage
```
./wu 
```

# LICENSE

See [LICENSE](./LICENSE)
