package fpu
import (
	"fmt"
	"math"
	"strconv"
)

//Convert an integer into the 2s complement binary representation
func intToTwos(an int, bitcount int) string {
	representation := ""
	if an < 0 { //If negative
		an *= -1 //Make it positive
		temp := fmt.Sprintf("%0"+fmt.Sprintf("%v", bitcount)+"v", strconv.FormatInt(int64(an), 2))
		for _, bit := range temp { //Then flip bits
			if bit == '0' {
				representation += "1"
			} else {
				representation += "0"
			}
		}
		//then add one
		toAdd, _ := strconv.ParseInt(representation, 2, 64)
		toAdd++
		representation = strconv.FormatInt(toAdd, 2)
		return representation
	}
	//Otherwise, if positive, it's just a binary conversion.
	return fmt.Sprintf("%0"+fmt.Sprintf("%v", bitcount)+"v", strconv.FormatInt(int64(an), 2))
}
//Convert a binary string into an integer
func twosComplementNew(binary string) int {
	if binary[0] == '1' {
		val, _ := strconv.ParseInt(binary, 2, 64)
		val -= power2int(len(binary))
		return int(val)
	}
	val, _ := strconv.ParseInt(binary, 2, 64)
	return int(val)
}
//Return 2^y
func power2int(y int) (a int64) {
	a = int64(1)
	a = a << uint64(y)
	return
}

//Normalise a binary float
func normalise(mantissa string, exponent string) string {
	//This does actually work with my implementation of NaN and Inf! 
	//Counter is 0 as 01 is never found in mantissa in that case
	//And exponent doesn't get modified
	leadingString := string(mantissa[0])
	if leadingString == "1" { //Establish if it should be 10 or 01 at the beginning
		leadingString += "0"
	} else {
		leadingString += "1"
	}
	counter := 0
	for i := 0; i < len(mantissa)-1; i++ { //Find where you have to adjust the binary point to
		if leadingString == mantissa[i:i+2] { //Once you find it
			counter = i
			break
		}
	}
	return rPadZero(mantissa[counter:], 24) + intToTwos(twosComplementNew(exponent)-counter, 8)
}

//Floating point binary to base10 float
func fpToBaseTen(mantissa string, exponent string) float64 {
	if mantissa == "000000000000000000000000" { //if special case then
		if exponent == "10000000" { //Represents NaN
			return math.NaN()
		}
		if exponent == "01111111" { //Positive Inf
			return math.Inf(1)
		}
		if exponent == "11111111" { //Negative Inf
			return math.Inf(-1)
		}
	}
	//Otherwise it's just a conversion
	return (float64(twosComplementNew(mantissa)) / math.Pow(2, 23)) * math.Pow(2, float64(twosComplementNew(exponent)))
}
//Split a binary float into mantissa and exponent
func splitFP(binary string) (string, string) {
	return binary[0:24], binary[24:]
}

//Base 10 float to binary floating point
func bTenToFP(a float64) string { //Credit To Stav/MrMercury20!
	count := 0
	mantissa := a
	if math.IsNaN(a) { //If NaN
		return "000000000000000000000000"+"10000000"
	}
	if math.IsInf(a, -1) {// If Negative inf
		return "000000000000000000000000"+"11111111"
	}
	if math.IsInf(a, 1) {// If Positive inf
		return "000000000000000000000000"+"01111111"
	}
	//This is kinda magic, just pokes the float around until it reaches the exponent limit or
	//mantissa precision limit
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
//Used for putting zeroes in front of numbers when numbers don't use all the bits available.
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
}
//FloatingPointDivide is a function for dividing two int32s as floats and giving an answer as a converted int32
func FloatingPointDivide(acc int32, divisor int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) / fpToBaseTen(splitFP(intToTwos(int(divisor),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
}
//FloatingPointAdd is a function for multiplying two int32s as floats and giving an answer as a converted int32
func FloatingPointAdd(acc int32, plus int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) + fpToBaseTen(splitFP(intToTwos(int(plus),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
}
//FloatingPointSub is a function for dividing two int32s as floats and giving an answer as a converted int32
func FloatingPointSub(acc int32, subtractor int32) int32{
	a := fpToBaseTen(splitFP(intToTwos(int(acc),32))) - fpToBaseTen(splitFP(intToTwos(int(subtractor),32)))
	return int32(twosComplementNew(normalise(splitFP(bTenToFP(a)))))
}
//FormatFP is a wrapper for fpToBaseTen and splitFP to convert an int32 into a float for printing
func FormatFP(number int32) string {
	return fmt.Sprint(fpToBaseTen(splitFP(intToTwos(int(number), 32))))
}
//FloatingPointConvert is a function to convert an input value from and int to the equivilent in floating point
func FloatingPointConvert(acc int32) int32{
	return int32(twosComplementNew(bTenToFP(float64(acc))))
}
//FloatToInt converts the accumulator from a float to an int. Takes the floor.
func FloatToInt(acc int32) int32 {
	return int32(math.Floor(fpToBaseTen(splitFP(intToTwos(int(acc), 32)))))
}