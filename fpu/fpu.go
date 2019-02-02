package fpu
import (
	"fmt"
	"math"
	"strconv"
)
func intToTwos(an int, bitcount int) string {
	representation := ""
	if an < 0 {
		an *= -1
		temp := fmt.Sprintf("%0"+fmt.Sprintf("%v", bitcount)+"v", strconv.FormatInt(int64(an), 2))
		for _, bit := range temp { //Flip bits
			if bit == '0' {
				representation += "1"
			} else {
				representation += "0"
			}
		}
		//Add one
		toAdd, _ := strconv.ParseInt(representation, 2, 64)
		toAdd++
		representation = strconv.FormatInt(toAdd, 2)
		if len(representation) != bitcount {
			fmt.Println("Oh fuck")
		}
		fmt.Println("repr:",representation, "Bitcount:", bitcount, "Len: ", len(representation))
		return representation
	}
	return fmt.Sprintf("%0"+fmt.Sprintf("%v", bitcount)+"v", strconv.FormatInt(int64(an), 2))
}

func twosComplementNew(binary string) int {
	//fmt.Println(binary)
	if binary[0] == '1' {
		//fmt.Println("Negative")
		val, _ := strconv.ParseInt(binary, 2, 64)
		fmt.Println(power2int(len(binary)))
		val -= power2int(len(binary))
		//fmt.Println("Value: ",val)
		return int(val)
	}
	val, _ := strconv.ParseInt(binary, 2, 64)
	return int(val)
}
func power2int(y int) (a int64) {
	a = int64(1)
	a = a << uint64(y)
	return
}

func normalise(mantissa string, exponent string) string {
	leadingString := string(mantissa[0])
	if leadingString == "1" { //Establish if it should be 10 or 01 at the beginning
		leadingString += "0"
	} else {
		leadingString += "1"
	}
	counter := 0
	for i := 0; i < len(mantissa)-1; i++ {
		if leadingString == mantissa[i:i+2] {
			counter = i
			break
		}
	}
	return rPadZero(mantissa[counter:], 24) + intToTwos(twosComplementNew(exponent)-counter, 8)
}
func fpToBaseTen(mantissa string, exponent string) float64 {
	return (float64(twosComplementNew(mantissa)) / math.Pow(2, 23)) * math.Pow(2, float64(twosComplementNew(exponent)))
}
func splitFP(binary string) (string, string) {
	return binary[0:24], binary[24:]
}
func bTenToFP(a float64) string { //Credit To Stav/MrMercury20!
	count := 0
	mantissa := a
	for (mantissa < float64(int(power2int(23)-1)) && mantissa > float64(-1*int(power2int(23)))) {
		mantissa = mantissa*2
		count = count+1
		if count > int(power2int(8)) {
			break
		}
	}

	for mantissa > float64(int(1 << 23)-1) || mantissa < float64(-1*int(power2int(23))) {
		mantissa = mantissa/2
		count = count-1
		if count < -1*(int(1 << 8)-1) {
			break
		}
	}
	count -= 23
	return intToTwos(int(mantissa), 24)+intToTwos(-count, 8)
}
func rPadZero(s string, size int) string { //recursively rightpad with zeroes
	if len(s) == size {
		return s
	}
	return rPadZero(s+"0", size)
}
//FloatingPointMultiply is a function for multiplying two int32s as floats and giving an answer as a converted int32
func FloatingPointMultiply(acc int32, multiplier int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) * fpToBaseTen(splitFP(intToTwos(int(multiplier),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
	//fmt.Println("Result:",a)
}
//FloatingPointDivide is a function for dividing two int32s as floats and giving an answer as a converted int32
func FloatingPointDivide(acc int32, divisor int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) / fpToBaseTen(splitFP(intToTwos(int(divisor),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
	//fmt.Println("Result:",a)
}
//FloatingPointAdd is a function for multiplying two int32s as floats and giving an answer as a converted int32
func FloatingPointAdd(acc int32, plus int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) + fpToBaseTen(splitFP(intToTwos(int(plus),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
	//fmt.Println("Result:",a)
}
//FloatingPointSub is a function for dividing two int32s as floats and giving an answer as a converted int32
func FloatingPointSub(acc int32, subtractor int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) - fpToBaseTen(splitFP(intToTwos(int(subtractor),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
	//fmt.Println("Result:",a)
}
//FormatFP is a wrapper for fpToBaseTen and splitFP to convert an int32 into a float for printing
func FormatFP(number int32) string {
	return fmt.Sprint(fpToBaseTen(splitFP(intToTwos(int(number), 32))))
}
//FloatingPointConvert is a function to convert an input value from and int to the equivilent in floating point
func FloatingPointConvert(acc int32) int32{
	return int32(twosComplementNew(bTenToFP(float64(acc))))
}