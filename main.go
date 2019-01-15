package main

/*
Important notes:
	INP with characters only takes 1 character, ignores the rest
	If program counter is negative, the program will exit!
	dumpMem is available for you to call during the program, somehow
*/
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	memAmount = 128   //Memory amount, in 32bit words
	debug     = false //Print every instruction as it's executed or not
)

var (
	memory         [memAmount]int32
	acc            int32
	programcounter int32
	instructions   = map[int32]func(int32){
		0:  hlt,
		1:  add,
		2:  sub,
		3:  sta,
		4:  brz,
		5:  brp,
		6:  bra,
		7:  lda,
		8:  out,
		9:  inp,
		10: asr,
		11: asl,
		12: mul,
		13: cmp,
		14: set,
	}
	program [memAmount]string
)

//Higher level logic
func input(message string) string { //Python Style input function, abstracts away the readers ad gives you just the info
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	return line
}

func dumpMem(binary bool) {
	for loc := 0; loc < len(memory); loc++ {
		if binary {
			fmt.Printf("%#04x %032b ", loc, uint32(memory[loc]))
		} else {
			fmt.Printf("%#04v % 16v ", loc, memory[loc])
		}
		if loc%4 == 3 { //print 4 locs per row
			fmt.Println()
		}
	}
}

//Lower level logic
func split(x int32) (int16, int16) { //splits into operand and operator, -32768 through 32767
	highBits := x / 65536                //Takes the first 4 nibbles, basically by doing an ASR many times which loses the last 4 nibbles
	lowBits := int16(x - 65536*highBits) //Calcs what's left
	//fmt.Println(highBits,lowBits)
	//fmt.Printf("%016b %016b\n", highBits, lowBits)
	return int16(highBits), lowBits
}
func twosComplement(binary string) int16 {
	total := -32768*int(binary[0])
	toAdd, _ := strconv.ParseInt(binary[1:len(binary)],2,16)
	total += int(toAdd)
	return int16(total)
}
func power2(y int32) (a int32) {
	a = int32(1)
	for i := int32(0); i < y; i++ {
		a *= 2
	}
	return
}

func main() {
	fmt.Printf("------ Prepared %06v bytes ------\n", memAmount*4)
	program[0] = "00000000000010010000000000000000"
	program[1] = "00000000000010000000000000000000"
	program[2] = "00000000000001000000000000000110"
	program[3] = "00000000000000100000000000000111"
	program[4] = "00000000000010000000000000000000"
	program[5] = "00000000000001100000000000000010"
	program[6] = "00000000000000000000000000000000"
	program[7] = "00000000000000000000000000000001"

	for index,value := range(program){
		if value != "" {
			memory[index] = (int32(twosComplement(value[0:16]))*65536)+(int32(twosComplement(value[16:32])))
		} else {
			break
		}
	}
	dumpMem(false)
	// //Arbitrary code goes here
	// set(196621) //STA 13
	// //fmt.Print("Splitty: ")
	// //fmt.Println(split(196621))
	// sta(0)
	// set(196623) //STA 15
	// sta(1)
	// set(-5)
	// sta(12)
	// //hlt(0)
	// bra(3)
	// //dumpMem(false)
	// set(-7)
	// bra(0)
	// //fmt.Println(decode(0))
	fmt.Println("-----------------------------------")
	//fmt.Println("started",programcounter)
	for programcounter >= 0 {
		data := fetch()
		if debug {
			fmt.Println("data:", data, "pc: ", programcounter, "acc: ", acc)
		}
		operator, operand := decode(data)
		execute(operator, operand)
	}
	fmt.Println("ended", programcounter)
	//dumpMem(false)
	//fmt.Println(memory)
}

//FDE
func fetch() int32 {
	programcounter++ //technically should be last, but return exits
	//fmt.Println(programcounter)
	if programcounter >= memAmount {
		hlt(0)
		return 0
	}
	return (memory[programcounter-1]) //the -1 counters the fact you *just* incremented it
}

func decode(data int32) (func(int32), int32) {
	opIndex, operand := split(data)
	operator := instructions[int32(opIndex)]
	if operator == nil {
		fmt.Println("Unknown instruction!", opIndex)
		return hlt, 0
	}
	//fmt.Println(opIndex, " : ", operand)
	return operator, int32(operand)
}

func execute(operator func(int32), operand int32) {
	operator(operand)
}

//Instructions
func hlt(_ int32) {
	fmt.Println("Exit!")
	programcounter = -1
}

func add(index int32) { //1
	if debug {
		fmt.Println("ADD", index)
	}
	acc += memory[index]
}

func sub(index int32) { //2
	if debug {
		fmt.Println("SUB", index)
	}
	acc -= memory[index]
}

func sta(index int32) { //3
	if debug {
		fmt.Println("STA", index)
	}
	memory[index] = acc
}

func brz(index int32) { //4
	if debug {
		fmt.Println("BRZ", index)
	}
	if acc == 0 {
		programcounter = index
	}
}

func brp(index int32) { //5
	if debug {
		fmt.Println("BRP", index)
	}
	if acc >= 0 {
		programcounter = index
	}
}

func bra(index int32) { //6
	if debug {
		fmt.Println("BRA", index)
	}
	programcounter = index
}

func lda(index int32) { //7
	if debug {
		fmt.Println("LDA", index)
	}
	acc = memory[index]
}

func out(outType int32) { //8
	if debug {
		fmt.Println("OUT", outType)
	}
	if outType == 0 { //Print as number
		fmt.Println(acc)
	} else {
		fmt.Print(string(acc)) //Print as character
	}
}

func inp(inType int32) { //9
	if debug {
		fmt.Println("INP", inType)
	}
	if inType == 0 { //Inp as number
		//acc = int32([]rune(input("Program requires number input: "))[0]) - 48 //HACKY - 1 char only for numbers is bad
		a, err := strconv.ParseInt(input("Number Input: "), 10, 32)
		if err != nil {
			fmt.Println("invalid input")
		} else {
			acc = int32(a)
		}
	} else {
		acc = int32([]rune(input("Char input: "))[0]) //inp as character
	}
}

func asr(times int32) { //10
	if debug {
		fmt.Println("ASR", times)
	}
	acc = acc / power2(times)
}

func asl(times int32) { //11
	if debug {
		fmt.Println("ASL", times)
	}
	acc = acc * power2(times)
}

func mul(operand int32) { //12
	if debug {
		fmt.Println("MUL", operand)
	}
	acc *= operand
}

func cmp(compareTo int32) { //13
	if debug {
		fmt.Println("CMP", compareTo)
	}
	if acc > compareTo {
		acc = 1
	} else if acc == compareTo {
		acc = 0
	} else {
		acc = -1
	}
}

func set(value int32) { //14
	if debug {
		fmt.Println("SET", value)
	}
	acc = value
}

func and(index int32) { //
	if debug {
		fmt.Println("AND", index)
	}
	acc = acc & memory[index]
}

func lor(index int32) { //Logical OR
	if debug {
		fmt.Println("LOR", index)
	}
	acc = acc | memory[index]
}

func xor(index int32) {
	if debug {
		fmt.Println("XOR", index)
	}
	acc = acc ^ memory[index]
}
