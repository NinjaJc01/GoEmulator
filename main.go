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
	"log"
	"os"
	"strconv"

	"./fpu"
)

const (
	memAmount = 52   //Memory amount, in 32bit words
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
		256: fad,
		257: fsu,
		258: fmu,
		259: fdv,
		260: fsq,
		261: fcv,
		262: fci,
	}
	// instructionStrings = map[int32]string{
	// 	0:   "hlt",
	// 	1:   "add",
	// 	2:   "sub",
	// 	3:   "sta",
	// 	4:   "brz",
	// 	5:   "brp",
	// 	6:   "bra",
	// 	7:   "lda",
	// 	8:   "out",
	// 	9:   "inp",
	// 	10:  "asr",
	// 	11:  "asl",
	// 	12:  "mul",
	// 	13:  "cmp",
	// 	14:  "set",
	// 	256: "fad",
	// 	257: "fsu",
	// 	258: "fmu",
	// 	259: "fdv",
	// 	260: "fsq",
	// 	261: "fcv",
	// }
	program []string
)

//Higher level logic
func input(message string) string { //Python Style input function, abstracts away the readers ad gives you just the info
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	return line
}
func readProgram() {
	file, err := os.Open("program.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		program = append(program, scanner.Text())
	}
	program = append(program, "")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func dumpMem(binary bool) {
	for loc := 0; loc < len(memory); loc++ {
		if binary {
			fmt.Printf(" %#04x %032b ", loc, uint32(memory[loc]))
		} else {
			fmt.Printf(" %#06v % 32v ", loc, memory[loc])
		}
		if loc%4 == 3 { //print 2 locs per row
			fmt.Println()
		}
	}
	fmt.Println()
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
	total := -32768 * int(binary[0])
	toAdd, _ := strconv.ParseInt(binary[1:len(binary)], 2, 16)
	total += int(toAdd)
	return int16(total)
}

func power2(y int32) (a int32) {
	a = int32(1)
	// for i := int32(0); i < y; i++ {
	// 	a *= 2
	// }
	a = a << uint32(y)
	return
}

func main() {
	fmt.Printf("----------------------------- Prepared %06v bytes ------------------------------\n", memAmount*4)
	readProgram() //Read in the program from program.txt, into slice of strings
	//Load program into memory
	for index, value := range program {
		if value != "" { //Should be \n terminated, therefore last line will be ""
			memory[index] = (int32(twosComplement(value[0:16])) * 65536) + (int32(twosComplement(value[16:32])))
		} else {
			break
		}
	}

	dumpMem(false)
	fmt.Println("----------------------------------------------------------------------------------")
	for programcounter >= 0 {
		data := fetch()
		// if debug { //Irritating
		// 	fmt.Println("data:", data, "pc: ", programcounter, "acc: ", acc)
		// }
		operator, operand := decode(data)
		execute(operator, operand)
	}
	fmt.Println("----------------------------------------------------------------------------------")
	dumpMem(false)
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
// func decodeString(data int32) string {
// 	opIndex, operand := split(data)
// 	operator := instructionStrings[int32(opIndex)]
// 	//fmt.Println(opIndex, " : ", operand)
// 	return fmt.Sprintln(operator, int32(operand))
// }

func execute(operator func(int32), operand int32) {
	operator(operand)
}

//Instructions
func hlt(_ int32) {
	fmt.Println("\nExit!")
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
	switch outType {
	case 0:
		fmt.Print(acc)
	case 1:
		fmt.Print(string(acc))
	case 2:
		fmt.Print(fpu.FormatFP(acc))
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
	acc *= memory[operand]
}

func cmp(compareTo int32) { //13
	compareTo = memory[compareTo]
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

func fad(location int32) {
	if debug {
		fmt.Println("FAD", location)
	}
	acc = fpu.FloatingPointAdd(acc, memory[location])
}

func fsu(location int32) {
	if debug {
		fmt.Println("FSU", location)
	}
	acc = fpu.FloatingPointSub(acc, memory[location])
}

func fmu(location int32) {
	if debug {
		fmt.Println("FMU", location)
	}
	acc = fpu.FloatingPointMultiply(acc, memory[location])
}

func fdv(location int32) {
	if debug {
		fmt.Println("FDV", location)
	}
	acc = fpu.FloatingPointDivide(acc, memory[location])
}

func fsq(_ int32) {
	if debug {
		fmt.Println("FSQ")
	}
	acc = fpu.FloatingPointMultiply(acc, acc)
}

func fcv(_ int32) {
	if debug {
		fmt.Println("FCV")
	}
	acc = fpu.FloatingPointConvert(acc)
}
func fci(_ int32) {
	if debug {
		fmt.Println("FCI")
	}
	acc = fpu.FloatToInt(acc)
}
