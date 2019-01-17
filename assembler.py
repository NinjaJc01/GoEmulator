"""No, pylint. This isn't a module, it's a script. So here's your docstring"""
import re
INSTRUCTIONS = {
    "hlt": 0,
    "dat": 0,
    "add": 1,
    "sub": 2,
    "sta": 3,
    "str": 3,
    "sto": 3,
    "brz": 4,
    "jmz": 4,
    "brp": 5,
    "jpl": 5,
    "bra": 6,
    "jmp": 6,
    "lda": 7,
    "out": 8,
    "otc": 8,
    "inp": 9,
    "asr": 10,
    "asl": 11,
    "mul": 12,
    "cmp": 13,
    "set": 14
}


def _twos_complement(number):
    if number == "":
        number = 0
    number = int(number)
    if number < 0:
        number *= -1
        number = ("{0:<b}".format(number).rjust(16, '0'))
        number_2 = ""
        for bit in number:  # Flip the bits
            if int(bit):
                number_2 += "0"
            else:
                number_2 += "1"
        # Convert to int as otherwise I'd have to do binary addition
        int_value = int(number_2, 2)
        int_value += 1  # add one
        number = ("{0:<b}".format(int_value).rjust(
            16, '0'))  # convert back to binary
        return number
    return "{0:b}".format(number).rjust(16, '0')


SYMBOL_TABLE = {}
PROGRAM_FILE = open("src.txt", "r")
PROGRAM = PROGRAM_FILE.readlines()
PROGRAM_FILE.close()
ADDRESS_COUNTER = 0
DELABELLED = ""
for line in PROGRAM:
    # Split at spaces
    line = line.rstrip()
    line = line.lstrip()
    line = re.sub(" +", " ",line) ##Remove duplicate spaces
    line = line.split(" ")
    print(len(line),line)
    # Check if first symbol is an instruction
    if line[0].lower() in INSTRUCTIONS:
        DELABELLED += " ".join(line)+"\n"
    else:
        # Add to symbol table if declaration
        SYMBOL_TABLE[line[0]] = ADDRESS_COUNTER
        line = " ".join(line)
        key = line.split(" ")[0]
        # Replace the key with "" to remove it
        line = re.sub(r"\b"+key+r"\b", "", line)
        line = line.lstrip()
        DELABELLED += line+"\n"
    ADDRESS_COUNTER += 1  # Used for symbol table

print("Symbol Table:", SYMBOL_TABLE)
TO_ASSEMBLE = ""

for line in DELABELLED.split("\n"):
    for key, value in SYMBOL_TABLE.items():
        line = re.sub(r"\b"+key+r"\b", str(value), line)
    print(repr(line))
    if line != '':
        TO_ASSEMBLE += line+"\n"

print("Final:\n", TO_ASSEMBLE, sep="", end='')
print()
OUT_PROGRAM = ""
for line in TO_ASSEMBLE.split("\n"):
    if len(line.split(" ")) > 1:
        items = line.split(" ")
        ##print("Items: ",items)
        operator = items[0]
        operand = items[1]
        if operand.startswith("/"):
            operand = 0
    else:
        if len(line) > 2:
            operator = line
            operand = 0
            if operator.lower() == "otc":
               operand = 1 
        else:
            break
    operator = INSTRUCTIONS[operator.lower()]
    print(operator, operand)
    OUT_PROGRAM += (_twos_complement(operator)+_twos_complement(operand)+"\n")
OUT_FILE = open("program.txt", "w")
print(OUT_PROGRAM, file=OUT_FILE, end="")
OUT_FILE.close()
