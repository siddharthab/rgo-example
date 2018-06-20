package main

import "C"

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

//export hello
func hello() {
	fmt.Println("A hello from Go!")
}

//export paste
func paste(strs, collapse, sep *C.struct_SEXPREC) *C.struct_SEXPREC {
	if RType(sep) == NILSXP {
		if RType(strs) != STRSXP {
			RError("argument should be a character vector")
		}
		v := RGoStrSlice(strs)
		for i, s := range v {
			v[i] = NAsafeString(s)
		}
		res := strings.Join(v, RGoString(collapse))
		return RString(res)
	} else {
		if RType(strs) != VECSXP {
			RError("argument should be a list of character vectors")
		}
		list := RGoList(strs)

		if len(list) == 0 {
			return RStringVec([]string{})
		}

		l := RCommonLength(list)

		// Join the elements together.
		separator := RGoString(sep)
		res := make([]string, l)
		for i, v := range list {
			for j, s := range RGoStrSlice(v) {
				s = NAsafeString(s)
				if i == 0 {
					res[j] = s
				} else {
					res[j] = res[j] + separator + s
				}
			}
		}
		return RStringVec(res)
	}
}

//export multiply
func multiply(list *C.struct_SEXPREC) *C.struct_SEXPREC {
	if RType(list) != VECSXP {
		RError("Not a list type")
	}
	v := RGoList(list)
	if len(v) == 0 {
		return RIntVec([]int{R_NA_Int})
	}

	l := RCommonLength(v)
	res := make([]float64, l)
	for i, v := range v {
		for j, n := range RGoFloatSlice(v) {
			if i == 0 {
				res[j] = n
			} else {
				res[j] = res[j] * n
			}
		}
	}
	return RFloatVec(res)
}

//export addOne
func addOne(numbers *C.struct_SEXPREC) {
	switch RType(numbers) {
	case REALSXP:
		{
			raw := RawFloatSlice(numbers)
			for i := range raw {
				raw[i] = raw[i] + 1
			}
		}
	case INTSXP:
		{
			raw := RawIntSlice(numbers)
			for i := range raw {
				raw[i] = raw[i] + 1
			}
		}
	default:
		RError("Unsupported type")
	}
}

//export busy
func busy(arg *C.struct_SEXPREC) *C.struct_SEXPREC {
	parallelism := RGoIntSlice(arg)[0]
	processorChan := make(chan int, 9)
	resultsChan := make(chan int, 9)
	var wg sync.WaitGroup
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		// Perform a random walk for a long time.
		go func() {
			defer wg.Done()
			for c := range processorChan {
				fmt.Println("Walk Number: ", c)
				randomWalk := 1
				for i := 0; i < 1e8; i++ {
					randomWalk = randomWalk + rand.Intn(8) - 4
				}
				resultsChan <- randomWalk
			}
			return
		}()
	}
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		processorChan <- i
	}
	close(processorChan)

	wg.Wait()
	close(resultsChan)

	res := make([]int, 0, 9)
	for result := range resultsChan {
		res = append(res, result)
	}
	return RIntVec(res)
}
