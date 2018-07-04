This is an example project to demonstrate how to use go code as part of an R
package's dynamically loaded libraries, using the bazel rules for R.

## Using R

To install the package and run all examples, run

```
R CMD INSTALL R/rgo
./examples.R
```

## Using bazel

To build the package, run

```
bazel build R/rgo
```

You can then run an R function with
```
R_LIBS_USER=$(pwd -P)/bazel-bin/R/rgo/lib Rscript -e 'rgo::hello()'
```

To install the package and run all examples, run

```
bazel run R/rgo:library
./examples.R
```
