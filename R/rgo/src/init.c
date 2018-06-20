#ifndef R_NO_REMAP
#define R_NO_REMAP
#endif

#ifndef STRICT_R_HEADERS
#define STRICT_R_HEADERS
#endif

#include <R_ext/Rdynload.h>
#include <Rinternals.h>

#include "go/src/rgo/rgo.h"

static const R_CallMethodDef callMethods[] = {
    {"hello", (DL_FUNC)&hello, 0},       {"paste", (DL_FUNC)&paste, 3},
    {"multiply", (DL_FUNC)&multiply, 1}, {"addOne", (DL_FUNC)&addOne, 1},
    {"busy", (DL_FUNC)&busy, 1},         {NULL, NULL, 0}};

// Initialize the shared library in the package.
// The name of this function should always match the package name.
void R_init_rgo(DllInfo* dll) {
  R_registerRoutines(dll, NULL, callMethods, NULL, NULL);
  R_useDynamicSymbols(dll, FALSE);
  R_forceSymbols(dll, TRUE);
}

// Prepare to unload the shared library in the package.
// The name of this function should always match the package name.
void R_unload_rgo(DllInfo* dll) {}
