package parser

import (
	"modernc.org/cc"
)

// exerpted from github.com/xlab/c-for-go/parser/predefined.go

var model64 = &cc.Model{
	Items: map[cc.Kind]cc.ModelItem{
		cc.Ptr:               {Size: 8, Align: 8, StructAlign: 8, More: "__TODO_PTR"},
		cc.UintPtr:           {Size: 8, Align: 8, StructAlign: 8, More: "uintptr"},
		cc.Void:              {Size: 0, Align: 1, StructAlign: 1, More: "__TODO_VOID"},
		cc.Char:              {Size: 1, Align: 1, StructAlign: 1, More: "int8"},
		cc.SChar:             {Size: 1, Align: 1, StructAlign: 1, More: "int8"},
		cc.UChar:             {Size: 1, Align: 1, StructAlign: 1, More: "byte"},
		cc.Short:             {Size: 2, Align: 2, StructAlign: 2, More: "int16"},
		cc.UShort:            {Size: 2, Align: 2, StructAlign: 2, More: "uint16"},
		cc.Int:               {Size: 4, Align: 4, StructAlign: 4, More: "int32"},
		cc.UInt:              {Size: 4, Align: 4, StructAlign: 4, More: "uint32"},
		cc.Long:              {Size: 8, Align: 8, StructAlign: 8, More: "int64"},
		cc.ULong:             {Size: 8, Align: 8, StructAlign: 8, More: "uint64"},
		cc.LongLong:          {Size: 8, Align: 8, StructAlign: 8, More: "int64"},
		cc.ULongLong:         {Size: 8, Align: 8, StructAlign: 8, More: "uint64"},
		cc.Float:             {Size: 4, Align: 4, StructAlign: 4, More: "float32"},
		cc.Double:            {Size: 8, Align: 8, StructAlign: 4, More: "float64"},
		cc.LongDouble:        {Size: 8, Align: 8, StructAlign: 4, More: "float64"},
		cc.Bool:              {Size: 1, Align: 1, StructAlign: 1, More: "bool"},
		cc.FloatComplex:      {Size: 8, Align: 8, StructAlign: 8, More: "complex64"},
		cc.DoubleComplex:     {Size: 16, Align: 16, StructAlign: 16, More: "complex128"},
		cc.LongDoubleComplex: {Size: 16, Align: 16, StructAlign: 16, More: "complex128"},
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
