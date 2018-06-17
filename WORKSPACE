workspace(
    name = "rgo",
)

git_repository(
    name = "io_bazel_rules_go",
    commit = "20b34406b4741ccd5cdb600059e80b4f0f61db09",
    remote = "https://github.com/bazelbuild/rules_go.git",
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

git_repository(
    name = "com_grail_rules_r",
    commit = "cd81fb0e4f28573c11a73a6dd9950d0f2086228e",
    remote = "https://github.com/grailbio/rules_r.git",
)

load("@com_grail_rules_r//R:dependencies.bzl", "r_rules_dependencies")

r_rules_dependencies()
