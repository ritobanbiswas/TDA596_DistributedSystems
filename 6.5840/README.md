## How to compile
### With race test
```
go build -race -buildmode=plugin -gcflags="all=-N -l" ../mrapps/wc.go

go run -race -gcflags="all=-N -l" mrcoordinator.go pg-*.txt

```
> On different session
```
go run -race -gcflags="all=-N -l" mrworker.go wc.so
```

### if no race test

```
go build  -buildmode=plugin -gcflags="all=-N -l" ../mrapps/wc.go

go run -gcflags="all=-N -l" mrcoordinator.go pg-*.txt
```

Again in different terminal.
```
go run -gcflags="all=-N -l" mrworker.go wc.so
```
