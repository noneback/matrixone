// Code generated by command: go run avx2.go -out max/avx2.s -stubs max/avx2_stubs.go. DO NOT EDIT.

#include "textflag.h"

// func int8MaxAvx2Asm(x []int8, r []int8)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·int8MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x0000000000000080, BX
	MOVQ         BX, X0
	VPBROADCASTB X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

int8MaxBlockLoop:
	CMPQ    DX, $0x000000c0
	JL      int8MaxTailLoop
	VPMAXSB (AX), Y1, Y1
	VPMAXSB 32(AX), Y2, Y2
	VPMAXSB 64(AX), Y3, Y3
	VPMAXSB 96(AX), Y4, Y4
	VPMAXSB 128(AX), Y5, Y5
	VPMAXSB 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x000000c0, DX
	JMP     int8MaxBlockLoop

int8MaxTailLoop:
	CMPQ    DX, $0x00000004
	JL      int8MaxDone
	VPMAXSB (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000020, DX
	JMP     int8MaxTailLoop

int8MaxDone:
	VPMAXSB      Y1, Y2, Y1
	VPMAXSB      Y1, Y3, Y1
	VPMAXSB      Y1, Y4, Y1
	VPMAXSB      Y1, Y5, Y1
	VPMAXSB      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMAXSB       X0, X1
	MOVOU        X1, (CX)
	RET

// func int16MaxAvx2Asm(x []int16, r []int16)
// Requires: AVX, AVX2, SSE2
TEXT ·int16MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x0000000000008000, BX
	MOVQ         BX, X0
	VPBROADCASTW X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

int16MaxBlockLoop:
	CMPQ    DX, $0x00000060
	JL      int16MaxTailLoop
	VPMAXSW (AX), Y1, Y1
	VPMAXSW 32(AX), Y2, Y2
	VPMAXSW 64(AX), Y3, Y3
	VPMAXSW 96(AX), Y4, Y4
	VPMAXSW 128(AX), Y5, Y5
	VPMAXSW 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000060, DX
	JMP     int16MaxBlockLoop

int16MaxTailLoop:
	CMPQ    DX, $0x00000004
	JL      int16MaxDone
	VPMAXSW (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000010, DX
	JMP     int16MaxTailLoop

int16MaxDone:
	VPMAXSW      Y1, Y2, Y1
	VPMAXSW      Y1, Y3, Y1
	VPMAXSW      Y1, Y4, Y1
	VPMAXSW      Y1, Y5, Y1
	VPMAXSW      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMAXSW       X0, X1
	MOVOU        X1, (CX)
	RET

// func int32MaxAvx2Asm(x []int32, r []int32)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·int32MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x0000000080000000, BX
	MOVQ         BX, X0
	VPBROADCASTD X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

int32MaxBlockLoop:
	CMPQ    DX, $0x00000030
	JL      int32MaxTailLoop
	VPMAXSD (AX), Y1, Y1
	VPMAXSD 32(AX), Y2, Y2
	VPMAXSD 64(AX), Y3, Y3
	VPMAXSD 96(AX), Y4, Y4
	VPMAXSD 128(AX), Y5, Y5
	VPMAXSD 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000030, DX
	JMP     int32MaxBlockLoop

int32MaxTailLoop:
	CMPQ    DX, $0x00000004
	JL      int32MaxDone
	VPMAXSD (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000008, DX
	JMP     int32MaxTailLoop

int32MaxDone:
	VPMAXSD      Y1, Y2, Y1
	VPMAXSD      Y1, Y3, Y1
	VPMAXSD      Y1, Y4, Y1
	VPMAXSD      Y1, Y5, Y1
	VPMAXSD      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMAXSD       X0, X1
	MOVOU        X1, (CX)
	RET

// func uint8MaxAvx2Asm(x []uint8, r []uint8)
// Requires: AVX, AVX2, SSE2
TEXT ·uint8MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ  x_base+0(FP), AX
	MOVQ  r_base+24(FP), CX
	MOVQ  x_len+8(FP), DX
	VPXOR Y0, Y0, Y0
	VPXOR Y1, Y1, Y1
	VPXOR Y2, Y2, Y2
	VPXOR Y3, Y3, Y3
	VPXOR Y4, Y4, Y4
	VPXOR Y5, Y5, Y5

uint8MaxBlockLoop:
	CMPQ    DX, $0x000000c0
	JL      uint8MaxTailLoop
	VPMAXUB (AX), Y1, Y1
	VPMAXUB 32(AX), Y2, Y2
	VPMAXUB 64(AX), Y3, Y3
	VPMAXUB 96(AX), Y4, Y4
	VPMAXUB 128(AX), Y5, Y5
	VPMAXUB 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x000000c0, DX
	JMP     uint8MaxBlockLoop

uint8MaxTailLoop:
	CMPQ    DX, $0x00000004
	JL      uint8MaxDone
	VPMAXUB (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000020, DX
	JMP     uint8MaxTailLoop

uint8MaxDone:
	VPMAXUB      Y1, Y2, Y1
	VPMAXUB      Y1, Y3, Y1
	VPMAXUB      Y1, Y4, Y1
	VPMAXUB      Y1, Y5, Y1
	VPMAXUB      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMAXUB       X0, X1
	MOVOU        X1, (CX)
	RET

// func uint16MaxAvx2Asm(x []uint16, r []uint16)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·uint16MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ  x_base+0(FP), AX
	MOVQ  r_base+24(FP), CX
	MOVQ  x_len+8(FP), DX
	VPXOR Y0, Y0, Y0
	VPXOR Y1, Y1, Y1
	VPXOR Y2, Y2, Y2
	VPXOR Y3, Y3, Y3
	VPXOR Y4, Y4, Y4
	VPXOR Y5, Y5, Y5

uint16MaxBlockLoop:
	CMPQ    DX, $0x00000060
	JL      uint16MaxTailLoop
	VPMAXUW (AX), Y1, Y1
	VPMAXUW 32(AX), Y2, Y2
	VPMAXUW 64(AX), Y3, Y3
	VPMAXUW 96(AX), Y4, Y4
	VPMAXUW 128(AX), Y5, Y5
	VPMAXUW 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000060, DX
	JMP     uint16MaxBlockLoop

uint16MaxTailLoop:
	CMPQ    DX, $0x00000004
	JL      uint16MaxDone
	VPMAXUW (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000010, DX
	JMP     uint16MaxTailLoop

uint16MaxDone:
	VPMAXUW      Y1, Y2, Y1
	VPMAXUW      Y1, Y3, Y1
	VPMAXUW      Y1, Y4, Y1
	VPMAXUW      Y1, Y5, Y1
	VPMAXUW      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMAXUW       X0, X1
	MOVOU        X1, (CX)
	RET

// func uint32MaxAvx2Asm(x []uint32, r []uint32)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·uint32MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ  x_base+0(FP), AX
	MOVQ  r_base+24(FP), CX
	MOVQ  x_len+8(FP), DX
	VPXOR Y0, Y0, Y0
	VPXOR Y1, Y1, Y1
	VPXOR Y2, Y2, Y2
	VPXOR Y3, Y3, Y3
	VPXOR Y4, Y4, Y4
	VPXOR Y5, Y5, Y5

uint32MaxBlockLoop:
	CMPQ    DX, $0x00000030
	JL      uint32MaxTailLoop
	VPMAXUD (AX), Y1, Y1
	VPMAXUD 32(AX), Y2, Y2
	VPMAXUD 64(AX), Y3, Y3
	VPMAXUD 96(AX), Y4, Y4
	VPMAXUD 128(AX), Y5, Y5
	VPMAXUD 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000030, DX
	JMP     uint32MaxBlockLoop

uint32MaxTailLoop:
	CMPQ    DX, $0x00000004
	JL      uint32MaxDone
	VPMAXUD (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000008, DX
	JMP     uint32MaxTailLoop

uint32MaxDone:
	VPMAXUD      Y1, Y2, Y1
	VPMAXUD      Y1, Y3, Y1
	VPMAXUD      Y1, Y4, Y1
	VPMAXUD      Y1, Y5, Y1
	VPMAXUD      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMAXUD       X0, X1
	MOVOU        X1, (CX)
	RET

// func float32MaxAvx2Asm(x []float32, r []float32)
// Requires: AVX, AVX2, SSE, SSE2
TEXT ·float32MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x00000000ff7fffff, BX
	MOVQ         BX, X0
	VBROADCASTSS X0, Y0
	VMOVUPS      Y0, Y1
	VMOVUPS      Y0, Y2
	VMOVUPS      Y0, Y3
	VMOVUPS      Y0, Y4
	VMOVUPS      Y0, Y5
	VMOVUPS      Y0, Y0

float32MaxBlockLoop:
	CMPQ   DX, $0x00000030
	JL     float32MaxTailLoop
	VMAXPS (AX), Y1, Y1
	VMAXPS 32(AX), Y2, Y2
	VMAXPS 64(AX), Y3, Y3
	VMAXPS 96(AX), Y4, Y4
	VMAXPS 128(AX), Y5, Y5
	VMAXPS 160(AX), Y0, Y0
	ADDQ   $0x000000c0, AX
	SUBQ   $0x00000030, DX
	JMP    float32MaxBlockLoop

float32MaxTailLoop:
	CMPQ   DX, $0x00000004
	JL     float32MaxDone
	VMAXPS (AX), Y1, Y1
	ADDQ   $0x00000020, AX
	SUBQ   $0x00000008, DX
	JMP    float32MaxTailLoop

float32MaxDone:
	VMAXPS       Y1, Y2, Y1
	VMAXPS       Y1, Y3, Y1
	VMAXPS       Y1, Y4, Y1
	VMAXPS       Y1, Y5, Y1
	VMAXPS       Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	MAXPS        X0, X1
	MOVOU        X1, (CX)
	RET

// func float64MaxAvx2Asm(x []float64, r []float64)
// Requires: AVX, AVX2, SSE2
TEXT ·float64MaxAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0xffefffffffffffff, BX
	MOVQ         BX, X0
	VBROADCASTSD X0, Y0
	VMOVUPD      Y0, Y1
	VMOVUPD      Y0, Y2
	VMOVUPD      Y0, Y3
	VMOVUPD      Y0, Y4
	VMOVUPD      Y0, Y5
	VMOVUPD      Y0, Y0

float64MaxBlockLoop:
	CMPQ   DX, $0x00000018
	JL     float64MaxTailLoop
	VMAXPD (AX), Y1, Y1
	VMAXPD 32(AX), Y2, Y2
	VMAXPD 64(AX), Y3, Y3
	VMAXPD 96(AX), Y4, Y4
	VMAXPD 128(AX), Y5, Y5
	VMAXPD 160(AX), Y0, Y0
	ADDQ   $0x000000c0, AX
	SUBQ   $0x00000018, DX
	JMP    float64MaxBlockLoop

float64MaxTailLoop:
	CMPQ   DX, $0x00000004
	JL     float64MaxDone
	VMAXPD (AX), Y1, Y1
	ADDQ   $0x00000020, AX
	SUBQ   $0x00000004, DX
	JMP    float64MaxTailLoop

float64MaxDone:
	VMAXPD       Y1, Y2, Y1
	VMAXPD       Y1, Y3, Y1
	VMAXPD       Y1, Y4, Y1
	VMAXPD       Y1, Y5, Y1
	VMAXPD       Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	MAXPD        X0, X1
	MOVOU        X1, (CX)
	RET