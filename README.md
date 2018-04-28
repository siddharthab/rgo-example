This is an example project to demonstrate how to use go code as part of an R
package's dynamically loaded libraries, using the bazel rules for R.

### go build

The first attempt is to generate a static archive of the go code and cgo
exported functions using `go build`.

`go env` for my builds was:
```
GOARCH="amd64"
GOBIN=""
GOCACHE="/home/ubuntu/.cache/go-build"
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="/home/ubuntu/Workspace/rgo-example/go/src"
GORACE=""
GOROOT="/home/ubuntu/.gimme/versions/go1.10.linux.amd64"
GOTMPDIR=""
GOTOOLDIR="/home/ubuntu/.gimme/versions/go1.10.linux.amd64/pkg/tool/linux_amd64"
GCCGO="gccgo"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build078063192=/tmp/go-build -gno-record-gcc-switches"
```

The R package was then built using:
```
bazel --bazelrc=/dev/null clean
bazel --bazelrc=/dev/null build R/rgo
```

You can then run the R function with (ignore the subsequent segfault for now please).
```
R_LIBS_USER=$(pwd -P)/bazel-bin/R/rgo/lib Rscript -e 'rgo::hello()'
```

## rules_go with linkmode=c-archive

We need to edit the BUILD file to use `bazel_build_archive` as the cc_deps attribute.
And then we can build with:
```
# Edit cc_deps attribute in R/rgo/BUILD.bazel.
# Then:
bazel --bazelrc=/dev/null clean
bazel --bazelrc=/dev/null build --copt=-fPIC R/rgo
```

This however results in the following error on linux:
```
/usr/bin/ld: /home/ubuntu/.cache/bazel/_bazel_ubuntu/5e9bbe5bc47cb3b6197ce9cc5d76e117/bazel-sandbox/6635601934392181222/execroot/rgo/bazel-out/k8-fastbuild/bin/go/src/hello/linux_amd64_stripped_c-archive/hello.a(go.o): relocation R_X86_64_TPOFF32 against `runtime.tlsg' can not be used when making a shared object; recompile with -fPIC
```
