workspace(
    name = "rgo",
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f70c35a8c779bb92f7521ecb5a1c6604e9c3edd431e50b6376d7497abc8ad3c1",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.11.0/rules_go-0.11.0.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

git_repository(
    name = "com_grail_rules_r",
    commit = "ce78de8ecb1e337cebbd44106e476ba989b07664",
    remote = "https://github.com/grailbio/rules_r.git",
)

load("@com_grail_rules_r//R:dependencies.bzl", "r_rules_dependencies")

r_rules_dependencies()
