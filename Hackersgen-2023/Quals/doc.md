# Language/VM Specifications

**Stack:** Stack available

**Memory:**  4096 bytes (Big Endian)

**Registers:** X and Y (8 bit)

**Flags:** FZ, FC (Flag Zero, Flag Carry)

**Notes:** Programs are stored in memory at address 0x0000 before execution. 1024 bytes are available for the program. Execution starts at 0x0000

### Language:


| Opcode | Mnemonic | Operand [size in bytes] | Description |
| ---------- | -------- | ----------------------- | ----------- |
| 0x40      | CLD |                                  | Clear all Flags |
| 0x50 | LDX | #Immediate val [1] | Load X with byte |
| 0x51 | LDY | #Immediate val [1] | Load Y with byte |
| 0x52 | STRX | #Immediate mem addr [2] | Store X at memory addr (word) |
| 0x53 | STRY | #Immediate mem addr [2] | Store X at memory addr (word) |
| 0x54 | LDRX | #Immediate mem addr [2] | Put value at memory addr (word) in X |
| 0x55 | LDRY | #Immediate mem addr [2] | Put value at memory addr (word) in Y |
| 0x60 | OUT |                         | Print content of X to stdout |
| 0x61 | IN |                         | Read char from stdin to X |
| 0x70 | CMPX | #Immediate val [1] | Compare X with byte and set flags accordingly: FZ = 0 if equal, FC = 1 if X > #Immediate |
| 0x71 | CMPY | #Immediate val [1] | Compare Y with byte and set flags accordingly: FZ = 0 if equal, FC = 1 if Y > #Immediate |
| 0x72 | JE | #Immediate mem addr [2] | If equal Jump to #Immediate address (word). All Jxx instructions use flags to decide |
| 0x73 | JRE | #Immediate relative mem addr [2] | If equal Jump to relative #Immediate relative address |
| 0x74 | JL | #Immediate mem addr [2] | If less than compared value, jump to #Immediate address |
| 0x75 | JRL | #Immediate relative mem addr [2] | If less than compared value, jump to #Immediate relative address |
| 0x76 | JLE | #Immediate mem addr [2] | If less or equal than compared value, jump to #Immediate address |
| 0x77 | JRLE | #Immediate relative mem addr [2] |             |
| 0x78 | JG | #Immediate mem addr [2] |             |
| 0x79 | JRG | #Immediate relative mem addr [2] |             |
| 0x7A | JGE | #Immediate mem addr [2] |             |
| 0x7B | JRGE | #Immediate relative mem addr [2] |             |
| 0xA0 | ADDX | #Immediate val [1] | Add a value to X. Set FC = 1 if overflow |
| 0xA1 | ADDXY |                         | Add Y to X. Set FC = 1 if overflow |
| 0xA2 | DECX | #Immediate val [1] | Subtract a value from X. Set FC = 1 if underflow |
| 0xA3 | DECXY |                         | Subtract Y from X. Set FC = 1 if underflow |
| 0xA4 | RORX |                         | Rotate X right of 1 bit |
| 0xA5 | ROLX |                         | Rotate X left of 1 bit |
| 0xA6 | XORX |                         | Excluisve OR between X and Y |
| 0xB0 | PUSHX |                         | Push X in stack |
| 0xB1 | POPX |                         | Pop value from stack in X |
| 0xB2 | PUSHY |                         | Push Y in stack |
| 0xB3 | POPY |                         | Pop value from stack in Y |
| 0xC0 | RMEMX |                         | Get two values from stack and read memory at address putting at X. Latest byte is less significant. The values are not popped |
| 0xC1 | WMEMX |                         | Get X and put it at memory pointed by stack as for RMEMX |
| 0xC2 | RMEMY |                         |             |
| 0xC3 | WMEMY |                         |             |
| 0x90 | NOP |                         | Do nothing |
| 0x91 | RET |                         | Exit program |
|            |          |                         |             |
|            |          |                         |             |



### Examples

`501052010050006054010060`      <-- Write a byte previously set
`405400006054000160`                    <-- Write the content of the first two bytes of memory
