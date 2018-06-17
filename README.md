This is an example project to demonstrate how to use go code as part of an R
package's dynamically loaded libraries, using the bazel rules for R.

To build the package, run

```
bazel build R/rgo
```

You can then run the R function with (ignore the subsequent segfault for now please).
```
R_LIBS_USER=$(pwd -P)/bazel-bin/R/rgo/lib Rscript -e 'rgo::hello()'
```

To install the package, run

```
bazel run R/rgo:library
```
