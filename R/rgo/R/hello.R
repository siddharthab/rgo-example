hello <- function() {
  .Call("C_hello")
  cat("And a hello from R!\n")
  invisible(NULL)
}
