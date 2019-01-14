package main

import "fmt"

var (
	memory         [128]uint32
	acc            uint32
	programcounter uint32
)

func main() {
	fmt.Println("------")
}

func add(index uint32) {
	acc += memory[index]
}

func sub(index uint32) {
	acc -= memory[index]
}

func sta(index uint32) {
	memory[index] = acc
}

func lda(index uint32) {
	acc = memory[index]
}

func bra(index uint32) {
	programcounter = index - 1
}

func brz(index uint32) {
	if acc == 0 {
		programcounter = index - 1
	}
}

func brp(index uint32) {
	if acc >= 0 {
		programcounter = index - 1
	}
}

func out(outType uint32) {
	if outType == 0 { //Print as number
		fmt.Println(acc)
	} else {
		fmt.Println(acc) //Print as character
	}
}
