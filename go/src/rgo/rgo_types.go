package main

/*
// These flags can also be provided through R's Makevars setting CGO_CCFLAGS and CGO_LDFLAGS
#cgo darwin CFLAGS: -I /Library/Frameworks/R.framework/Headers
#cgo darwin LDFLAGS: -L /Library/Frameworks/R.framework/Resources/lib -lR
#cgo linux CFLAGS: -I /usr/share/R/include
#cgo linux LDFLAGS: -L /usr/lib/R/lib -lR

#ifndef R_NO_REMAP
#define R_NO_REMAP
#endif

#ifndef STRICT_R_HEADERS
#define STRICT_R_HEADERS
#endif

#include <R.h>
#include <Rinternals.h>

void R_error(char* s) {
	Rf_error(s);
	free(s);
}
*/
import "C"

import (
	"unsafe"
)

// All Go objects must be created in Go code.

const Go_NaString = "___Go__NA__STRING__"

func NAsafeString(s string) string {
	if s == Go_NaString {
		return "NA"
	}
	return s
}

var R_NA_Int = int(C.R_NaInt)

// SEXP Type mapping
const (
	NILSXP = iota
	LGLSXP
	INTSXP
	REALSXP
	CPLXSXP
	STRSXP
	VECSXP
	RAWSXP
)

func RType(a *C.struct_SEXPREC) int {
	switch C.TYPEOF(a) {
	case C.NILSXP:
		return NILSXP
	case C.LGLSXP:
		return LGLSXP
	case C.INTSXP:
		return INTSXP
	case C.REALSXP:
		return REALSXP
	case C.CPLXSXP:
		return CPLXSXP
	case C.STRSXP:
		return STRSXP
	case C.VECSXP:
		return VECSXP
	case RAWSXP:
		return RAWSXP
	default:
		RError("Unknown type")
	}
	return -1
}

type Rboolean int

func RError(s string) {
	C.R_error(C.CString(s))
}

func RLength(a *C.struct_SEXPREC) int {
	return int(C.Rf_length(a))
}

func RGoString(s *C.struct_SEXPREC) string {
	if C.TYPEOF(s) != C.STRSXP {
		RError("Not an STRSXP type.")
		return ""
	}

	e := C.STRING_ELT(s, 0)
	if e == C.R_NaString {
		return Go_NaString
	}

	return C.GoString(C.R_CHAR(e))
}

func RGoStrSlice(v *C.struct_SEXPREC) []string {
	if C.TYPEOF(v) != C.STRSXP {
		RError("Not a STRSXP type.")
		return nil
	}

	l := int(C.Rf_length(v))
	res := make([]string, 0, l)
	for i := 0; i < l; i++ {
		e := C.STRING_ELT(v, C.R_xlen_t(i))
		if e == C.R_NaString {
			res = append(res, Go_NaString)
		} else {
			res = append(res, C.GoString(C.R_CHAR(e)))
		}
	}
	return res
}

func RGoFloatSlice(v *C.struct_SEXPREC) []float64 {
	if C.TYPEOF(v) != C.REALSXP && C.TYPEOF(v) != C.INTSXP {
		RError("Not a REALSXP or INTSXP type.")
		return nil
	}

	// TODO: Benchmark REAL_ELT and avoid having to create an unsafe slice.
	// Another option is having inline pointer arithmetic functions.
	l := int(C.Rf_length(v))
	res := make([]float64, 0, l)
	switch C.TYPEOF(v) {
	case C.REALSXP:
		{
			raw := (*[1 << 30]_Ctype_double)(unsafe.Pointer(C.REAL(v)))[:l:l]
			for i := 0; i < l; i++ {
				e := raw[i]
				res = append(res, float64(e))
			}
		}
	case C.INTSXP:
		{
			raw := (*[1 << 30]_Ctype_int)(unsafe.Pointer(C.INTEGER(v)))[:l:l]
			for i := 0; i < l; i++ {
				e := raw[i]
				if e == C.R_NaInt {
					res = append(res, float64(C.R_NaReal))
				} else {
					res = append(res, float64(e))
				}
			}
		}
	}
	return res
}

func RGoIntSlice(v *C.struct_SEXPREC) []int {
	if C.TYPEOF(v) != C.INTSXP {
		RError("Not an INTSXP type.")
		return nil
	}

	l := int(C.Rf_length(v))
	raw := (*[1 << 30]_Ctype_int)(unsafe.Pointer(C.INTEGER(v)))[:l:l]
	res := make([]int, 0, l)
	for i := 0; i < l; i++ {
		e := raw[i]
		res = append(res, int(e))
	}
	return res
}

func RGoCplxSlice(v *C.struct_SEXPREC) []complex128 {
	if C.TYPEOF(v) != C.CPLXSXP {
		RError("Not a CPLXSXP type.")
		return nil
	}

	l := int(C.Rf_length(v))
	raw := (*[1 << 30]C.Rcomplex)(unsafe.Pointer(C.COMPLEX(v)))[:l:l]
	res := make([]complex128, 0, l)
	for i := 0; i < l; i++ {
		e := raw[i]
		res = append(res, complex(e.r, e.i))
	}
	return res
}

func RGoLglSlice(v *C.struct_SEXPREC) []Rboolean {
	if C.TYPEOF(v) != C.LGLSXP {
		RError("Not a LGLSXP type.")
		return nil
	}

	l := int(C.Rf_length(v))
	raw := (*[1 << 30]C.Rboolean)(unsafe.Pointer(C.LOGICAL(v)))[:l:l]
	res := make([]Rboolean, 0, l)
	for i := 0; i < l; i++ {
		e := raw[i]
		res = append(res, Rboolean(e))
	}
	return res
}

func RGoRawSlice(v *C.struct_SEXPREC) []byte {
	if C.TYPEOF(v) != C.RAWSXP {
		RError("Not a RAWSXP type.")
		return nil
	}

	l := int(C.Rf_length(v))
	raw := (*[1 << 30]C.Rbyte)(unsafe.Pointer(C.RAW(v)))[:l:l]
	res := make([]byte, 0, l)
	for i := 0; i < l; i++ {
		e := raw[i]
		res = append(res, byte(e))
	}
	return res
}

func RGoList(v *C.struct_SEXPREC) []*C.struct_SEXPREC {
	if C.TYPEOF(v) != C.VECSXP {
		RError("Not a VECSXP type.")
		return nil
	}

	l := int(C.Rf_length(v))
	res := make([]*C.struct_SEXPREC, 0, l)
	for i := 0; i < l; i++ {
		e := C.VECTOR_ELT(v, C.R_xlen_t(i))
		res = append(res, e)
	}
	return res
}

func RGoNames(v *C.struct_SEXPREC) []string {
	return RGoStrSlice(C.Rf_getAttrib(v, C.R_NamesSymbol))
}

func RSetNames(v *C.struct_SEXPREC, names []string) *C.struct_SEXPREC {
	return C.Rf_setAttrib(v, C.R_NamesSymbol, RStringVec(names))
}

func RGoDim(v *C.struct_SEXPREC) []int {
	return RGoIntSlice(C.Rf_getAttrib(v, C.R_DimSymbol))
}

func RSetDim(v *C.struct_SEXPREC, dims []int) *C.struct_SEXPREC {
	return C.Rf_setAttrib(v, C.R_DimSymbol, RIntVec(dims))
}

func RGoClass(v *C.struct_SEXPREC) string {
	return RGoString(C.Rf_getAttrib(v, C.R_ClassSymbol))
}

func RSetClass(v *C.struct_SEXPREC, name string) *C.struct_SEXPREC {
	return C.Rf_setAttrib(v, C.R_ClassSymbol, RString(name))
}

func RString(s string) *C.struct_SEXPREC {
	if s == Go_NaString {
		return C.R_NaString
	}

	res := C.Rf_protect(C.Rf_allocVector(C.STRSXP, C.R_xlen_t(1)))
	C.SET_STRING_ELT(res, 0, C.Rf_mkCharLen(C._GoStringPtr(s), C.int(C._GoStringLen(s))))
	C.Rf_unprotect(1)
	return res
}

func RStringVec(v []string) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.STRSXP, C.R_xlen_t(len(v))))
	for i, s := range v {
		if s == Go_NaString {
			C.SET_STRING_ELT(res, C.long(i), C.R_NaString)
		} else {
			C.SET_STRING_ELT(res, C.long(i), C.Rf_mkCharLen(C._GoStringPtr(s), C.int(C._GoStringLen(s))))
		}
	}
	C.Rf_unprotect(1)
	return res
}

func RFloatVec(v []float64) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.REALSXP, C.R_xlen_t(len(v))))
	raw := (*[1 << 30]_Ctype_double)(unsafe.Pointer(C.REAL(res)))[:len(v):len(v)]
	for i, n := range v {
		raw[i] = _Ctype_double(n)
	}
	C.Rf_unprotect(1)
	return res
}

func RIntVec(v []int) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.INTSXP, C.R_xlen_t(len(v))))
	raw := (*[1 << 30]_Ctype_int)(unsafe.Pointer(C.INTEGER(res)))[:len(v):len(v)]
	for i, n := range v {
		raw[i] = _Ctype_int(n)
	}
	C.Rf_unprotect(1)
	return res
}

func RComplexVec(v []complex128) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.CPLXSXP, C.R_xlen_t(len(v))))
	raw := (*[1 << 30]C.Rcomplex)(unsafe.Pointer(C.COMPLEX(res)))[:len(v):len(v)]
	for i, n := range v {
		raw[i] = C.Rcomplex{_Ctype_double(real(n)), _Ctype_double(imag(n))}
	}
	C.Rf_unprotect(1)
	return res
}

func RLogicalVec(v []Rboolean) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.LGLSXP, C.R_xlen_t(len(v))))
	raw := (*[1 << 30]C.Rboolean)(unsafe.Pointer(C.LOGICAL(res)))[:len(v):len(v)]
	for i, n := range v {
		raw[i] = C.Rboolean(n)
	}
	C.Rf_unprotect(1)
	return res
}

func RRawVec(v []byte) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.RAWSXP, C.R_xlen_t(len(v))))
	raw := (*[1 << 30]C.Rbyte)(unsafe.Pointer(C.RAW(res)))[:len(v):len(v)]
	for i, n := range v {
		raw[i] = C.Rbyte(n)
	}
	C.Rf_unprotect(1)
	return res
}

func RList(v []*C.struct_SEXPREC) *C.struct_SEXPREC {
	res := C.Rf_protect(C.Rf_allocVector(C.VECSXP, C.R_xlen_t(len(v))))
	for i, e := range v {
		C.SET_VECTOR_ELT(res, C.long(i), e)
	}
	C.Rf_unprotect(1)
	return res
}

// RCommonLength returns the common length of all the objects. Returns an error
// if the lengths differ.
func RCommonLength(v []*C.struct_SEXPREC) int {
	l := -1
	for _, s := range v {
		if l == -1 {
			l = RLength(s)
		} else if l != RLength(s) {
			RError("lengths of all character vectors in the list must be the same")
		}
	}
	return l
}

func RawIntSlice(v *C.struct_SEXPREC) []_Ctype_int {
	l := int(C.Rf_length(v))
	return (*[1 << 30]_Ctype_int)(unsafe.Pointer(C.INTEGER(v)))[:l:l]
}

func RawFloatSlice(v *C.struct_SEXPREC) []_Ctype_double {
	l := int(C.Rf_length(v))
	return (*[1 << 30]_Ctype_double)(unsafe.Pointer(C.REAL(v)))[:l:l]
}
