def genString(string):
    print(
"""STARTLOOP LDA START
BRZ END
OUT 1
LDA STARTLOOP
ADD ONE
STA STARTLOOP
BRA STARTLOOP
ONE DAT 1
HLT
START """, end='')
    for c in string+"\n":
        print("DAT", ord(c))
    print("DAT 0")
    print("END HLT 0")


genString(input("Enter a string to be outputted: "))
