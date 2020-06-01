"".calc STEXT nosplit size=49 args=0x20 locals=0x0
	0x0000 00000 (main.go:3)	TEXT	"".calc(SB), NOSPLIT|ABIInternal, $0-32
	0x0000 00000 (main.go:3)	PCDATA	$0, $-2
	0x0000 00000 (main.go:3)	PCDATA	$1, $-2
	0x0000 00000 (main.go:3)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:3)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:3)	FUNCDATA	$2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:3)	PCDATA	$0, $0
	0x0000 00000 (main.go:3)	PCDATA	$1, $0
	0x0000 00000 (main.go:3)	MOVQ	$0, "".~r2+24(SP)
	0x0009 00009 (main.go:3)	MOVQ	$0, "".~r3+32(SP)
	0x0012 00018 (main.go:4)	MOVQ	"".a+8(SP), AX
	0x0017 00023 (main.go:4)	ADDQ	"".b+16(SP), AX
	0x001c 00028 (main.go:4)	MOVQ	AX, "".~r2+24(SP)
	0x0021 00033 (main.go:4)	MOVQ	"".a+8(SP), AX
	0x0026 00038 (main.go:4)	SUBQ	"".b+16(SP), AX
	0x002b 00043 (main.go:4)	MOVQ	AX, "".~r3+32(SP)
	0x0030 00048 (main.go:4)	RET
	0x0000 48 c7 44 24 18 00 00 00 00 48 c7 44 24 20 00 00  H.D$.....H.D$ ..
	0x0010 00 00 48 8b 44 24 08 48 03 44 24 10 48 89 44 24  ..H.D$.H.D$.H.D$
	0x0020 18 48 8b 44 24 08 48 2b 44 24 10 48 89 44 24 20  .H.D$.H+D$.H.D$ 
	0x0030 c3                                               .
"".main STEXT size=68 args=0x0 locals=0x28
	0x0000 00000 (main.go:7)	TEXT	"".main(SB), ABIInternal, $40-0
	0x0000 00000 (main.go:7)	MOVQ	(TLS), CX
	0x0009 00009 (main.go:7)	CMPQ	SP, 16(CX)
	0x000d 00013 (main.go:7)	PCDATA	$0, $-2
	0x000d 00013 (main.go:7)	JLS	61
	0x000f 00015 (main.go:7)	PCDATA	$0, $-1
	0x000f 00015 (main.go:7)	SUBQ	$40, SP
	0x0013 00019 (main.go:7)	MOVQ	BP, 32(SP)
	0x0018 00024 (main.go:7)	LEAQ	32(SP), BP
	0x001d 00029 (main.go:7)	PCDATA	$0, $-2
	0x001d 00029 (main.go:7)	PCDATA	$1, $-2
	0x001d 00029 (main.go:7)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001d 00029 (main.go:7)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001d 00029 (main.go:7)	FUNCDATA	$2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001d 00029 (main.go:8)	PCDATA	$0, $0
	0x001d 00029 (main.go:8)	PCDATA	$1, $0
	0x001d 00029 (main.go:8)	MOVQ	$66, (SP)
	0x0025 00037 (main.go:8)	MOVQ	$77, 8(SP)
	0x002e 00046 (main.go:8)	CALL	"".calc(SB)
	0x0033 00051 (main.go:9)	MOVQ	32(SP), BP
	0x0038 00056 (main.go:9)	ADDQ	$40, SP
	0x003c 00060 (main.go:9)	RET
	0x003d 00061 (main.go:9)	NOP
	0x003d 00061 (main.go:7)	PCDATA	$1, $-1
	0x003d 00061 (main.go:7)	PCDATA	$0, $-2
	0x003d 00061 (main.go:7)	CALL	runtime.morestack_noctxt(SB)
	0x0042 00066 (main.go:7)	PCDATA	$0, $-1
	0x0042 00066 (main.go:7)	JMP	0
	0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 2e 48  eH..%....H;a.v.H
	0x0010 83 ec 28 48 89 6c 24 20 48 8d 6c 24 20 48 c7 04  ..(H.l$ H.l$ H..
	0x0020 24 42 00 00 00 48 c7 44 24 08 4d 00 00 00 e8 00  $B...H.D$.M.....
	0x0030 00 00 00 48 8b 6c 24 20 48 83 c4 28 c3 e8 00 00  ...H.l$ H..(....
	0x0040 00 00 eb bc                                      ....
	rel 5+4 t=17 TLS+0
	rel 47+4 t=8 "".calc+0
	rel 62+4 t=8 runtime.morestack_noctxt+0
go.cuinfo.packagename. SDWARFINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
go.loc."".calc SDWARFLOC size=0
go.info."".calc SDWARFINFO size=84
	0x0000 03 22 22 2e 63 61 6c 63 00 00 00 00 00 00 00 00  ."".calc........
	0x0010 00 00 00 00 00 00 00 00 00 01 9c 00 00 00 00 01  ................
	0x0020 0f 61 00 00 03 00 00 00 00 01 9c 0f 62 00 00 03  .a..........b...
	0x0030 00 00 00 00 02 91 08 0f 7e 72 32 00 01 03 00 00  ........~r2.....
	0x0040 00 00 02 91 10 0f 7e 72 33 00 01 03 00 00 00 00  ......~r3.......
	0x0050 02 91 18 00                                      ....
	rel 0+0 t=24 type.int+0
	rel 9+8 t=1 "".calc+0
	rel 17+8 t=1 "".calc+49
	rel 27+4 t=30 gofile../Users/zhanghui/gomod/ironzhang/projects/runtime/func/main.go+0
	rel 37+4 t=29 go.info.int+0
	rel 48+4 t=29 go.info.int+0
	rel 62+4 t=29 go.info.int+0
	rel 76+4 t=29 go.info.int+0
go.range."".calc SDWARFRANGE size=0
go.debuglines."".calc SDWARFMISC size=15
	0x0000 04 02 11 06 69 06 6a 06 41 04 01 03 7d 06 01     ....i.j.A...}..
go.loc."".main SDWARFLOC size=0
go.info."".main SDWARFINFO size=33
	0x0000 03 22 22 2e 6d 61 69 6e 00 00 00 00 00 00 00 00  ."".main........
	0x0010 00 00 00 00 00 00 00 00 00 01 9c 00 00 00 00 01  ................
	0x0020 00                                               .
	rel 9+8 t=1 "".main+0
	rel 17+8 t=1 "".main+68
	rel 27+4 t=30 gofile../Users/zhanghui/gomod/ironzhang/projects/runtime/func/main.go+0
go.range."".main SDWARFRANGE size=0
go.debuglines."".main SDWARFMISC size=18
	0x0000 04 02 03 01 14 0a a5 9c 06 5f 06 9c 71 04 01 03  ........._..q...
	0x0010 7a 01                                            z.
""..inittask SNOPTRDATA size=24
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 00 00 00                          ........
gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
