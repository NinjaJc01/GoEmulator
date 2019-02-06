PROG_FILE = open("program.txt", "r")
PROGRAM = PROG_FILE.readlines()
PROGRAM_EDITED = []
INSTRUCTIONS = {
    0:  "hlt/dat",
    1:  "add",
    2:  "sub",
    3:  "sta",
    4:  "brz",
    5:  "brp",
    6:  "bra",
    7:  "lda",
    8:  "out",
    9:  "inp",
    10: "asr",
    11: "asl",
    12: "mul",
    13: "cmp",
    14: "set",
    256: "fad",
	257: "fsu",
	258: "fmu",
	259: "fdv",
	260: "fsq",
    261: "fsr",
	262: "fcv",
	263: "fci"
}


def twos_comp(val, bits):
    """compute the 2's complement of int value val"""
    if (val & (1 << (bits - 1))) != 0:  # if sign bit is set e.g., 8bit: 128-255
        val = val - (1 << bits)        # compute negative value
    return val                         # return positive value as is


for line in PROGRAM:
    PROGRAM_EDITED.append(line.rstrip())
##print(PROGRAM_EDITED)

for line in PROGRAM_EDITED:
    operator = twos_comp(int(line[:16], 2), 16)
    operand = twos_comp(int(line[16:], 2), 16)
    print(INSTRUCTIONS[operator], operand)
