hello <- function() {
  cat("\n=====\n")
  .Call(Go_hello)
  cat("And a hello from R!\n")

  invisible(NULL)
}

gpaste <- function() {
  fruits <- c("apples", "oranges", "pineapples")
  collapsed <- .Call(Go_paste, fruits, " and ", NULL)
  cat("\n=====\n")
  print("Collapsed:")
  print(collapsed)

  quantities <- c("five", "two", "three")
  pasted <- .Call(Go_paste, list(quantities, fruits), NULL, " ")
  cat("\n=====\n")
  print("Pasted:")
  print(pasted)

  #pasted <- .Call(Go_paste, data.frame(quantities, fruits), NULL, " ")
  #cat("\n=====\n")
  #print("Pasted (w/ data.frame)")
  #print(pasted)

  quantities[1] <- NA
  pasted <- .Call(Go_paste, list(quantities, fruits), NULL, " ")
  cat("\n=====\n")
  print("Pasted (w/ unknown # apples):")
  print(pasted)

  invisible(NULL)
}

numbers <- function() {
  cat("\n=====\n")
  print("Multiplication (1:3 * 2):")
  print(.Call(Go_multiply, list(1:3, rep(2, 3))))

  cat("\n=====\n")
  print("Multiplication (w/ NA):")
  print(.Call(Go_multiply, list(c(1, NA, 3), rep(2, 3))))

  a <- runif(1e5)
  b <- runif(1e5)

  expected <- a * b
  actual <- .Call(Go_multiply, list(a, b))

  cat("\n=====\n")
  print("Multiplication (w/ large arrays):")
  print(all.equal(expected, actual))

  cat("\n=====\n")
  print("Multiplication (w/ empty list):")
  print(.Call(Go_multiply, list()))

  invisible(NULL)
}

inplace <- function() {
  mat <- matrix(1:9, nrow = 3, dimnames = list(c("a", "b", "c"), c("i", "j", "k")))
  cat("\n=====\n")
  print("Inplace (add 1):")
  print("Original = ")
  print(mat)
  .Call(Go_addOne, mat)
  print("+ 1 = ")
  print(mat)

  invisible(NULL)
}

busy <- function() {
  cat("\n=====\n")
  print("Random walk:")
  print(.Call(Go_busy, as.integer(parallel::detectCores() - 1)))

  invisible(NULL)
}

# Expected errors below

gpaste_error1 <- function() {
  .Call(Go_paste, list(c("one", "two"), c("mississipi")), NULL, " ")
}

gpaste_error2 <- function() {
  .Call(Go_paste, list(c(1, 2), c("one", "two")), NULL, " ")
}
