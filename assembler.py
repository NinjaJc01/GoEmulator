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
    "inp": 9,
    "asr": 10,
    "asl": 11,
    "mul": 12,
    "cmp": 13,
    "set": 14
}

symbol_table = {}
PROGRAM_FILE = open("sample program.txt")
PROGRAM = PROGRAM_FILE.readlines()
PROGRAM_FILE.close()
address_counter = 0
delabelled = ""
for line in PROGRAM:
    # Split at spaces
    line = line.rstrip()
    line = line.split(" ")
    # Check if first symbol is an instruction
    if line[0].lower() in INSTRUCTIONS:
        delabelled += " ".join(line)+"\n"
    else:
        # Add to symbol table if declaration
        symbol_table[line[0]] = address_counter
        line = " ".join(line)
        key = line.split(" ")[0]
        # Replace the key with "" to remove it
        line = re.sub(r"\b"+key+r"\b", "", line)
        line = line.lstrip()
        delabelled += line+"\n"
    address_counter += 1  # Used for symbol table
a = "".join(PROGRAM)
print(symbol_table)
to_assemble = ""
for line in delabelled.split("\n"):
    for key, value in symbol_table.items():
        line = re.sub(r"\b"+key+r"\b", str(value), line)
    print(line)
    to_assemble += line+"\n"
print("Final:\n", to_assemble, sep="")
