package parser

import (
	"modernc.org/cc"
)

// exerpted from github.com/xlab/c-for-go/parser/predefined.go

var model64 = &cc.Model{
	Items: map[cc.Kind]cc.ModelItem{
		cc.Ptr:               {8, 8, 8, "__TODO_PTR"},
		cc.UintPtr:           {8, 8, 8, "uintptr"},
		cc.Void:              {0, 1, 1, "__TODO_VOID"},
		cc.Char:              {1, 1, 1, "int8"},
		cc.SChar:             {1, 1, 1, "int8"},
		cc.UChar:             {1, 1, 1, "byte"},
		cc.Short:             {2, 2, 2, "int16"},
		cc.UShort:            {2, 2, 2, "uint16"},
		cc.Int:               {4, 4, 4, "int32"},
		cc.UInt:              {4, 4, 4, "uint32"},
		cc.Long:              {8, 8, 8, "int64"},
		cc.ULong:             {8, 8, 8, "uint64"},
		cc.LongLong:          {8, 8, 8, "int64"},
		cc.ULongLong:         {8, 8, 8, "uint64"},
		cc.Float:             {4, 4, 4, "float32"},
		cc.Double:            {8, 8, 4, "float64"},
		cc.LongDouble:        {8, 8, 4, "float64"},
		cc.Bool:              {1, 1, 1, "bool"},
		cc.FloatComplex:      {8, 8, 8, "complex64"},
		cc.DoubleComplex:     {16, 16, 16, "complex128"},
		cc.LongDoubleComplex: {16, 16, 16, "complex128"},
	},
}

// builtins not supported by modernc.org/cc
var builtinBase = `
#define __builtin_va_list void *
#define __asm(x)
#define __inline
#define __inline__
#define __signed
#define __signed__
#define __const const
#define __extension__
#define __attribute__(x)
#define __attribute(x)
#define __restrict
#define __volatile__

#define __builtin_inff() (0)
#define __builtin_infl() (0)
#define __builtin_inf() (0)
#define __builtin_fabsf(x) (0)
#define __builtin_fabsl(x) (0)
#define __builtin_fabs(x) (0)

#define __INTRINSIC_PROLOG(name)
`

// var builtinBaseUndef = `
// #undef __llvm__
// #undef __BLOCKS__
// `

// var basePredefines = `
// #define __STDC_HOSTED__ 1
// #define __STDC_VERSION__ 199901L
// #define __STDC__ 1
// #define __GNUC__ 4
// #define __GNUC_PREREQ(maj,min) 0
// #define __POSIX_C_DEPRECATED(ver)

// #define __FLT_MIN__ 0
// #define __DBL_MIN__ 0
// #define __LDBL_MIN__ 0

// void __GO__(char*, ...);
// `
