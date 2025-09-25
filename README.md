<!-- markdownlint-disable-file MD010 -->

# aferox

The `aferox` packages expands on [`github.com/spf13/afero`](https://github.com/spf13/afero) by adding more `afero.Fs` implementations, as well as various `afero.Fs` utility functions.

[Go Doc](https://pkg.go.dev/github.com/unmango/aferox)

## github

The `github` package adds multiple implementations of `afero.Fs` for interacting with the GitHub API as if it were a filesystem.
In general it can turn a GitHub url into an `afero.Fs`.

```go
fs := github.NewFs(github.NewClient(nil))

file, _ := fs.Open("https://github.com/unmango")

// ["go", "thecluster", "pulumi-baremetal", ...]
file.Readdirnames(420)
```

This package lives in a separate module to avoid adding a dependendy on GitHub to `aferox`.

[Go Doc](https://pkg.go.dev/github.com/unmango/aferox/github)

```shell
go get github.com/unmango/aferox/github
```

## ignore

The `ignore` package adds a filtering `afero.Fs` that accepts a `.gitignore` file and ignores paths matched by it.

```go
base := afero.NewMemMapFs()

gitignore, _ := os.Open("path/to/my/.gitignore")

fs, _ := ignore.NewFsFromGitIgnoreReader(base, gitignore)
```

## filter

The `filter` package adds a filtering implementation of `afero.Fs` similar to `afero.RegExpFs`, but for predicates.

```go
base := afero.NewMemMapFs()

fs := filter.NewFs(base, func(path string) bool {
	return filepath.Ext(path) == ".go"
})
```

## docker

The `docker` package adds a docker `afero.Fs` implementation for operating on the filesystem of a container.

```go
client := client.NewClientWithOpts(client.FromEnv)

fs := docker.NewFs(client, "my-container-id")
```

This package lives in a separate module to avoid adding a dependency on docker to `aferox`.

[Go Doc](https://pkg.go.dev/github.com/unmango/aferox/docker)

```shell
go get github.com/unmango/aferox/docker
```

## testing

The `testing` package adds helper stubs for mocking filesystems in tests.

```go
fs := &testing.Fs{
	CreateFunc: func(name string) (afero.File, error) {
		return nil, errors.New("simulated error")
	}
}
```

## writer

The `writer` package adds a readonly `afero.Fs` implementation that dumps all file writes to the provided `io.Writer`.
Currently paths are ignored and there are no delimeters separating files.

```go
buf := &bytes.Buffer{}
fs := writer.NewFs(buf)

_ = afero.WriteFile(fs, "test.txt", []byte("testing"), os.ModePerm)
_ = afero.WriteFile(fs, "other.txt", []byte("blah"), os.ModePerm)

// "testingblah"
buf.String()
```

## context

The `context` package adds the `context.Fs` interface for filesystem implementations that accept a `context.Context` per operation.
It has a basic test suite and generally works but should be considered a 🚧 work in progress 🚧.

The `context` package re-exports various functions and types from the standard `context` packge for convenience.
Currently the creation functions focus on adapting external `context.Fs` implementations to an `afero.Fs` to be used with the other utility functions.

```go
var base context.Fs = mypkg.NewEffectfulFs()

fs := context.BackgroundFs(base)

var accessor context.AccessorFunc = func() context.Context {
	return context.Background()
}

// Equivalent to `context.BackgroundFs`
fs := context.NewFs(base, accessor)
```

The `context.AferoFs` interface is a union of `afero.Fs` and `context.Fs`, i.e. exposing both `fs.Create` and `fs.CreateContext`.
I'm not sure if this actually has any value but it exists.

The `context.Discard` function adapts an `afero.Fs` to a `context.AferoFs` by ignoring the `context.Context` argument.

```go
base := afero.NewMemMapFs()

var fs context.AferoFs = context.Discard(base)
```

## protofs

The `protofs` module contains a silly little idea I had for communicating filesystem information over gRPC.
The protobuf definitions are over at <https://github.com/unmango/protofs/>.
It... works?
[It is not well tested](./protofs/grpc/v1alpha1/e2e_test.go) and probably performs poorly.
Use at your own risk.

[Go Doc](https://pkg.go.dev/github.com/unmango/aferox/protofs)

### Server

```go
// Serve any afero.Fs implementation
var fs afero.Fs

server := grpc.NewServer()
protofsv1alpha1.RegisterFsServer(server, fs)
protofsv1alpha1.RegisterFileServer(server, fs)

lis, _ := net.Listen("unix", "/tmp/fs.sock")
server.Serve(lis)
```

### Client

```go
conn, _ := grpc.NewClient("unix:///tmp/fs.sock",
  grpc.WithTransportCredentials(insecure.NewCredentials()),
)

var fs afero.Fs = protofsv1alpha1.NewFs(conn)
```
