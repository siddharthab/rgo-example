#include <stdio.h>

#include "go/src/hello/hello.h"

void C_hello() {
  hello();
  printf("And a hello from C!\n");
  return;
}
