# zlog
集成了zap、file-rotatelogs
解决了token too long的问题

#### Usage

```
// default options
DefaultFormat     = Format_CONSOLE
DefaultFileName   = "log/log"
DefaultMaxFile    = 30
DefaultCallerSkip = 1

// usage
zlog.Infof("")

// change options
zlog.SyncFile(zlog.Config{})
```

#### safe writer

```
// Resolve error: "token too long"
zlog.SafeWriter() *io.PipeWriter 
```