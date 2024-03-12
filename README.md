# zlog
zap & file-rotatelogs & safe writer & gorm logger

#### Usage

```
// default options
DefaultFormat     = FormatConsole
DefaultFileName   = "log/log"
DefaultMaxFile    = 30
DefaultCallerSkip = 1

// usage
zlog.Infof("")

// change options
zlog.SyncFile(zlog.Options{
    Format     Format
	FileName   string
	MaxFile    uint
	CallerSkip int
})
```

#### safe writer

```
// Resolve error: "token too long"
zlog.SafeWriter() *io.PipeWriter 
```