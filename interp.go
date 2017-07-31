package main

import (
	//	"io"
	"fmt"
	"math"

//	"unsafe"
)

const (
	r0 = iota
	r1
	r2
	wa
	wb
	wc
	xl
	xr
	xs
	ia
	ra
	cp
)
const (
	xt = xs
)

var regName = map[int]string{
	r0: "r0",
	r1: "r1",
	r2: "r2",
	wa: "wa",
	wb: "wb",
	wc: "wc",
	xl: "xl",
	xr: "xr",
	xs: "xs",
	ia: "ia",
	ra: "ra",
	cp: "cp",
}

/*
const (
	atn = iota
	chp
	cos
	etx
	lnf
	sin
	sqr
	tan
)
*/
/*	operation encoded in four parts:
	_w	gives width in bits
	_m	gives mask to extract value
	operand encod
*/
const (
	op_w  = 8
	dst_w = 4
	src_w = 4
	off_w = 16
	op_m  = 1<<op_w - 1
	dst_m = 1<<dst_w - 1
	src_m = 1<<src_w - 1
	off_m = 1<<off_w - 1
	op_   = 0
	dst_  = op_ + op_w
	src_  = dst_ + dst_w
	off_  = src_ + src_w
)

const (
	stackLength = 1000
)

var (
	ip       int
	mem      [100000]int
	reg      [16]int
	stackEnd int
	memLast  int // index of last allocated memory word
)

func interp() {
//	instn := 0
	var long1, long2 int64
	//	var int1,int2 int32
	var int1, int2 int32
	var prcstack [32]int
	var inst, dst, src, off int
	var overflow bool
	var op int
//	var f1, f2 float32
	var d1 float64

	ip = start
	/*
				fmt.Printf(" startup r1 %v r2 %v wa %v wb %v wc %v xl %v xr %v xs %v cp %v ia %v\n",
					reg[r1], reg[r2], reg[wa], reg[wb], reg[wc], reg[xl], reg[xr], reg[xs],
					 reg[cp],int32(reg[ia]))
		fmt.Printf("start interp mem len %v  ip %v r0 %v\n",len(mem),ip,reg[r0])
		fmt.Printf("s_aaa %v s_yyy %v\n",s_aaa,s_yyy)
	*/
run:
	for {
		if reg[r0] != 0 {
			fmt.Println("r0 not zero", reg[r0], ip)
			panic("r0 not zero")
		}
		if ip < s_aaa || ip > s_yyy {
			fmt.Println("ip out of range ", ip)
			fmt.Println(" aaa yyy", s_aaa, s_yyy)
			return
		}
		inst = mem[ip]
/*
		instn++
		if instn > 150000 {
			fmt.Println("instruction limit exceeded", instn)
			return
		}
*/

		op = inst & op_m
		dst = inst >> dst_ & dst_m
		src = inst >> src_ & src_m
		off = inst >> off_ & off_m
		if itrace {
			fmt.Printf(" r1 %v r2 %v wa %v wb %v wc %v xl %v xr %v xs %v cp %v ia %v\n",
				reg[r1], reg[r2], reg[wa], reg[wb], reg[wc], reg[xl], reg[xr], reg[xs],
				reg[cp], int32(reg[ia]))
			fmt.Printf(" %v %v %v %v %v %v\n", ip,
				op, opName[op], regName[dst], regName[src], off)

		}
		//		fmt.Printf(" %v %v %v %v=%v %v=%v %v\n", ip,
		//		op, opName[op], regName[dst], reg[dst],
		//		regName[src], reg[src], off)
		ip++
		switch op {
		case stmt:
			if strace {
				fmt.Printf(" r1 %v r2 %v wa %v wb %v wc %v xl %v xr %v xs %v cp %v ia %v\n",
					reg[r1], reg[r2], reg[wa], reg[wb], reg[wc],
					reg[xl], reg[xr], reg[xs], reg[cp], int32(reg[ia]))
				fmt.Printf("  %v\n", stmt_text[off])
			}
		case mov:
			reg[dst] = reg[src]
		case brn:
			ip = off
		case bsw:
			if off > 0 {
				if reg[dst] >= reg[r1] {
					ip = off
					continue
				}
			}
			// when reach here, ip is pointing to first iff entry
			ip = mem[ip+reg[dst]]
		case bri:
			ip = reg[dst]
		case lei:
			reg[dst] = mem[reg[dst]-1]
		case ppm:
			ip = off
		case prc:
			prcstack[off] = mem[reg[xs]]
			reg[xs]++
		case exi:
			if itrace {
				fmt.Println("PROC:EXI  ", reg[xs], mem[reg[xs]], ip)
			}
			// dst is procedure identifier  if 'n' type procedure, 0 otherwise.
			// off is exit number
			if off >= 100 {
				ip = prcstack[off/100]
				reg[r1] = off % 100
			} else {
				// pop return address from stack
				ip = mem[reg[xs]]
				reg[xs]++
				reg[r1] = off
			}
		case err, erb:
			reg[wa] = off
			ip = error_
		case icv:
			reg[dst]++
		case dcv:
			reg[dst]--
		case add:
			reg[dst] += reg[src]
		case sub:
			reg[dst] -= reg[src]
		case ica:
			reg[dst]++
		case dca:
			reg[dst]--
		case beq:
			if reg[dst] == reg[src] {
				ip = off
			}
		case bge:
			if reg[dst] >= reg[src] {
				ip = off
			}
		case bgt:
			if reg[dst] > reg[src] {
				ip = off
			}
		case bne:
			if reg[dst] != reg[src] {
				ip = off
			}
		case ble:
			if reg[dst] <= reg[src] {
				ip = off
			}
		case blt:
			if reg[dst] < reg[src] {
				ip = off
			}
		case blo:
			if reg[dst] < reg[src] {
				ip = off
			}
		case bhi:
			if reg[dst] > reg[src] {
				ip = off
			}
		case bnz:
			if reg[dst] != 0 {
				ip = off
			}
		case bze:
			if reg[dst] == 0 {
				ip = off
			}
		case lct:
			reg[dst] = reg[src]
		case bct:
			reg[dst]--
			if reg[dst] > 0 {
				ip = off
			}
		case aov:
			if uint64(reg[dst])+uint64(reg[src]) > math.MaxUint32 {
				ip = off
			}
			reg[dst] += reg[src]
		case bev:
			if reg[dst]&1 == 0 {
				ip = off
			}
		case bod:
			if reg[dst]&1 != 0 {
				ip = off
			}
		case lcp:
			reg[cp] = reg[dst]
		case scp:
			reg[dst] = reg[cp]
		case lcw:
			reg[dst] = mem[reg[cp]]
			reg[cp]++
		case icp:
			reg[cp]++
		case ldi:
			reg[ia] = reg[dst]
		//TODO: Reminder that ia is SIGNED integer when doing arithmetic
		case adi:
			long1, long2 = int64(int32(reg[ia])), int64(int32(reg[dst]))
			long1 += long2
			if long1 > math.MaxInt32 || long1 < math.MinInt32 {
				overflow = true
			} else {
				overflow = false
				reg[ia] = int(long1)
			}
		case mli:
			long1, long2 = int64(int32(reg[ia])), int64(int32(reg[dst]))
			long1 *= long2
			if long1 > math.MaxInt32 || long1 < math.MinInt32 {
				overflow = true
			} else {
				overflow = false
				reg[ia] = int(long1)
			}
		case sbi:
			long1, long2 = int64(int32(reg[ia])), int64(int32(reg[dst]))
			long1 -= long2
			if long1 > math.MaxInt32 || long1 < math.MinInt32 {
				overflow = true
			} else {
				overflow = false
				reg[ia] = int(long1)
			}
		case dvi:
			if reg[dst] == 0 {
				overflow = true
			} else {
				overflow = false
				int1, int2 = int32(reg[ia]), int32(reg[dst])
				int1 /= int2
				reg[ia] = int(int1)
			}
		case rmi:
			if reg[dst] == 0 {
				overflow = true
			} else {
				overflow = false
				int1, int2 = int32(reg[ia]), int32(reg[dst])
				int1 %= int2
				reg[ia] = int(int1)
			}
			//		case sti:
			//			reg[dst] = reg[ia]
		case ngi:
			int1 = int32(reg[ia])
			if int1 == math.MinInt32 {
				overflow = true
			} else {
				overflow = false
				reg[ia] = int(-int1)
			}
		case ino:
			if !overflow {
				ip = off
			}
		case iov:
			if overflow {
				ip = off
			}
		case ieq:
			if int32(reg[ia]) == 0 {
				ip = off
			}
		case ige:
			if int32(reg[ia]) >= 0 {
				ip = off
			}
		case igt:
			if int32(reg[ia]) > 0 {
				ip = off
			}
		case ile:
			if int32(reg[ia]) <= 0 {
				ip = off
			}
		case ilt:
			if int32(reg[ia]) < 0 {
				ip = off
			}
		case ine:
			if reg[ia] != 0 {
				ip = off
			}
		case ldr:
			reg[ra] = reg[dst]
		case str:
			reg[dst] = reg[ra]
		case adr:
//			f1 = math.Float32frombits(reg[ra])
//			f2 = math.Float32frombits(reg[dst])
//			reg[ra] = math.Float32bits(f1 + f2)
		case sbr:
//			f1 = math.Float32frombits(reg[ra])
//			f2 = math.Float32frombits(reg[dst])
//			reg[ra] = math.Float32bits(f1 - f2)
		case mlr:
//			f1 = math.Float32frombits(reg[ra])
//			f2 = math.Float32frombits(reg[dst])
//			reg[ra] = math.Float32bits(f1 * f2)
		case dvr:
//			f1 = math.Float32frombits(reg[ra])
//			f2 = math.Float32frombits(reg[dst])
//			reg[ra] = math.Float32bits(f1 / f2)
		case rov:
//			d1 = float64(math.Float32frombits(reg[ra]))
//			if math.IsNaN(d1) || math.IsInf(d1, 0) {
//				ip = off
//			}
		case rno:
//			d1 = float64(math.Float32frombits(reg[ra]))
//			if !(math.IsNaN(d1) || math.IsInf(d1, 0)) {
//				ip = off
//			}
		case ngr:
//			f1 = math.Float32frombits(reg[ra])
//			reg[ra] = math.Float32bits(-f1)
		case req:
//			if math.Float32frombits(reg[ra]) == 0.0 {
//				ip = off
//			}
		case rge:
//			if math.Float32frombits(reg[ra]) >= 0.0 {
//				ip = off
//			}
		case rgt:
//			if math.Float32frombits(reg[ra]) < 0.0 {
//				ip = off
//			}
		case rle:
//			if math.Float32frombits(reg[ra]) <= 0.0 {
//				ip = off
//			}
		case rlt:
//			if math.Float32frombits(reg[ra]) < 0.0 {
//				ip = off
//			}
		case rne:
//			if math.Float32frombits(reg[ra]) != 0.0 {
//				ip = off
//			}
		case plc:
			reg[dst] = reg[dst] + reg[src] + 2
		case psc:
			reg[dst] = reg[dst] + reg[src] + 2
		case cne:
			if reg[dst] != reg[src] {
				ip = off
			}
		case cmc:
			s1 := mem[reg[xl]:]
			s2 := mem[reg[xr]:]
			n := int(reg[wa])
			for i := 0; i < n; i++ {
				if s1[i] < s2[i] {
					ip = reg[r1]
					reg[xl], reg[xr] = 0, 0
					break
				} else if s1[i] > s2[i] {
					ip = reg[r2]
					reg[xl], reg[xr] = 0, 0
					break
				}
			}
			reg[xl], reg[xr] = 0, 0
		case trc:
			n := int(reg[wa])
			ixl := int(reg[xl])
			ixr := int(reg[xr])
			for i := 0; i < n; i++ {
				mem[ixl+i] = mem[ixr+int(mem[ixl+i])]
			}
		case flc:
			panic("flc not implemented")
		case anb:
			reg[dst] &= reg[src]
		case orb:
			reg[dst] |= reg[src]
		case xob:
			reg[dst] ^= reg[src]
		case rsh:
			reg[dst] = reg[dst] >> uint32(off)
		case lsh:
			reg[dst] = reg[dst] << uint32(off)
		case nzb:
			if reg[dst] != 0 {
				ip = off
			}
		case zrb:
			if reg[dst] == 0 {
				ip = off
			}
		case mfi:
			if off != 0 && int32(reg[ia]) < 0 {
				ip = off
			}
		case itr:
//			reg[ia] = math.Float32bits(float32(int32(reg[ia])))
		case rti:
//			d1 = float64(math.Float32frombits(reg[ra]))
//			if math.IsNaN(d1) || math.IsInf(d1, 0) {
//				ip = off
//			}
			reg[ia] = int(int32(d1))
		case ctb, ctw:
			reg[dst] += off
		case cvm:
			long1 = int64(int32(reg[ia]))*10 - int64(reg[wb]-'0')
			if long1 > math.MaxInt32 || long1 < math.MinInt32 {
				ip = off
			}
			reg[ia] = int(long1)
		case cvd:
			int1 = int32(reg[ia])
			reg[ia] = int(int1 / 10)
			reg[wa] = int(-(int1 % 10) + '0')
		case mvc, mvw:
			n := int(reg[wa])
			for i := 0; i < n; i++ {
				mem[reg[xr]+int(i)] = mem[reg[xl]+int(i)]
			}
			reg[xl] += reg[wa]
			reg[xr] += reg[wa]
		case mcb, mwb:
			for i := 0; i < int(reg[wa]); i++ {
				mem[reg[xr]-1-int(i)] = mem[reg[xl]-1-int(i)]
			}
			reg[xl] -= reg[wa]
			reg[xr] -= reg[wa]
		case chk:
			if int(reg[xs]) < stackEnd+100 {
				ip = sec06 // branch to stack overflow section
			}
		case move:
			reg[dst] = reg[src]
		case call:
			if itrace {
				fmt.Println("PROC:CALL ", off, prc_names[off], reg[xs], mem[reg[xs]], ip)
			}
			reg[xs]--
			mem[reg[xs]] = ip
			ip = off
		case sys:
			if otrace {
				fmt.Printf("OSINT CALL %v\n", sysName[off])
			}
			reg[r1] = syscall(off)
			if otrace {
				fmt.Printf("OSINT RETN %v %v\n", sysName[off], reg[r1])
			}

			if reg[r1] == 999 {
				break run // end execution
			}
			/*
				case decv:
					int1 = int32(reg[ia])
					reg[ia] = int(reg[ia] / 10)
					int1 = int1 % 10
					reg[ia] = int(-int1 + 0x30)
			*/
		case jsrerr:
			if itrace {
				fmt.Println("PROC:JSRE ", ip, reg[r1], off)
			}
			if reg[r1] == 0 {
				ip = ip + off // skip around exi/ppm's
			} else {
				ip = ip + reg[r1] - 1
			}
		case load:
			reg[dst] = mem[reg[src]+off]
			if itrace {
				fmt.Println("  load ", regName[dst], "<-", mem[reg[src]+off], reg[src]+off)
			}
		case loadcfp:
			reg[dst] = 2147483647
		case loadi:
			reg[dst] = off
		case nop:
			// nop means 'no operation' so there is nothing to do here
		case pop:
			mem[reg[src]+off] = mem[reg[dst]]
			reg[dst]++
		case popr:
			reg[src] = mem[reg[dst]]
			reg[dst]++
		case push:
			reg[dst]--
			mem[reg[dst]+off] = mem[reg[src]+off]
		case pushi:
			reg[dst]--
			mem[reg[dst]] = off
		case pushr:
			reg[dst]--
			mem[reg[dst]] = reg[src]
		case realop:
			panic("realop not implemented")
		case store:
			mem[reg[src]+off] = reg[dst]
			if itrace {
				fmt.Println("  store ", regName[dst], reg[dst], "->", reg[src]+off)
			}
		default:
			fmt.Println("unknown opcode ", op)
			panic("unknown opcode")
		}
	}
}

// startup
// xs = one past stack base (subtract 1 to get first stack entry)
// xr = address first word data area
// xl = address last word data area
// wa=initial stack pointer
// wb=wc=ia=ra=cp=0

func startup() int {
	//	var memMinimum int = len(program) + 3 * (maxreclen + 2) + stackLength
	for i := 0; i < len(program); i++ {
		mem[i] = program[i]
		//		fmt.Printf("mem %v %x\n",i,mem[i])
	}
	memLast = int(len(program)) + 10
	// scblk is just four words, sufficient to hold null string
	scblk0 = memLast
	scblk1 = memLast + 4
	scblk2 = scblk1 + maxreclen + 1
	memLast += maxreclen + 1
	stackEnd = int(memLast)
	memLast += stackLength
	stackStart := int(memLast)
	memLast += 10
	reg[xl] = memLast // start data area
	// allocate 10000 words for initial data area
	memLast += 10000
	reg[xr] = memLast // end data area
	reg[xs] = stackStart
	reg[wa] = reg[xs] - 1
	interp()
	return 0
}

var opName = map[int]string{
	add:     "add",
	adi:     "adi",
	adr:     "adr",
	anb:     "anb",
	aov:     "aov",
	bct:     "bct",
	beq:     "beq",
	bev:     "bev",
	bge:     "bge",
	bgt:     "bgt",
	bhi:     "bhi",
	ble:     "ble",
	blo:     "blo",
	blt:     "blt",
	bne:     "bne",
	bnz:     "bnz",
	bod:     "bod",
	bri:     "bri",
	brn:     "brn",
	bsw:     "bsw",
	bze:     "bze",
	call:    "call",
	chk:     "chk",
	cmc:     "cmc",
	cne:     "cne",
	cvd:     "cvd",
	cvm:     "cvm",
	dca:     "dca",
	dcv:     "dcv",
	dvi:     "dvi",
	dvr:     "dvr",
	erb:     "erb",
	err:     "err",
	exi:     "exi",
	flc:     "flc",
	ica:     "ica",
	icp:     "icp",
	icv:     "icv",
	ieq:     "ieq",
	ige:     "ige",
	igt:     "igt",
	ile:     "ile",
	ilt:     "ilt",
	ine:     "ine",
	ino:     "ino",
	iov:     "iov",
	itr:     "itr",
	jsrerr:  "jsrerr",
	lcp:     "lcp",
	lcw:     "lcw",
	ldi:     "ldi",
	ldr:     "ldr",
	lei:     "lei",
	load:    "load",
	loadcfp: "loadcfp",
	loadi:   "loadi",
	lsh:     "lsh",
	mfi:     "mfi",
	mli:     "mli",
	mlr:     "mlr",
	mov:     "mov",
	move:    "move",
	mvc:     "mvc",
	mvw:     "mvw",
	mwb:     "mwb",
	ngi:     "ngi",
	ngr:     "ngr",
	nzb:     "nzb",
	orb:     "orb",
	plc:     "plc",
	pop:     "pop",
	popr:    "popr",
	ppm:     "ppm",
	prc:     "prc",
	psc:     "psc",
	push:    "push",
	pushi:   "pushi",
	pushr:   "pushr",
	realop:  "realop",
	req:     "req",
	rge:     "rge",
	rgt:     "rgt",
	rle:     "rle",
	rlt:     "rlt",
	rmi:     "rmi",
	rne:     "rne",
	rno:     "rno",
	rov:     "rov",
	rsh:     "rsh",
	rti:     "rti",
	sbi:     "sbi",
	sbr:     "sbr",
	scp:     "scp",
	store:   "store",
	sub:     "sub",
	sys:     "sys",
	trc:     "trc",
	xob:     "xob",
	zrb:     "zrb",
}
