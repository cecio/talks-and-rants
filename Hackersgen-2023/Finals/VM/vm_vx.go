//
// go build -v  -trimpath -ldflags="-s -w -extldflags=-static" -o vm vm.go
//

package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const MAXMEM = 4096
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var killSwitch = ""
var kmemAddr = 0

var debug = false
var novirus = false

// VM holds the state of the virtual machine.
type VM struct {
	memory []byte
	stack  []byte
	ip     uint16 // Instruction pointer
	rx     byte   // X Register
	ry     byte   // Y Register
	fz     byte   // Zero Flag (zero after CMP if not Equal)
	fc     byte   // Carry Flag (zero after CMP if Less)
}

// NewVM creates a new instance of the virtual machine.
func NewVM() *VM {
	return &VM{
		memory: make([]byte, MAXMEM),
		stack:  make([]byte, 0),
		ip:     0,
		rx:     0,
		ry:     0,
		fz:     0,
		fc:     0,
	}
}

// Print log information
func printLog(vm *VM, cmd string, operand string) {
	if debug != true {
		return
	}
	if operand != "" {
		log.Printf("0x%04X    %s 0x%s --> X:%02X Y:%02X ZF:%02X ZC:%02X",
			vm.ip, cmd, operand, vm.rx, vm.ry, vm.fz, vm.fc)
	} else {
		log.Printf("0x%04X    %s --> X:%02X Y:%02X ZF:%02X ZC:%02X",
			vm.ip, cmd, vm.rx, vm.ry, vm.fz, vm.fc)
	}
}

// Check if KillSwitch is present at given address
func killed(arr []byte, startIndex int, sequence string) bool {
	// Check if the startIndex is valid
	if startIndex < 0 || startIndex+len(sequence) > len(arr) {
		return false
	}

	// Compare the sequence at the specified index
	for i := 0; i < len(sequence); i++ {
		if arr[startIndex+i] != sequence[i] {
			return false
		}
	}

	return true
}

// ExecutePseudoCommands executes the given pseudo-commands on the VM.
func (vm *VM) ExecutePseudoCommands(commands []byte) {

	// First address for Virus memory corruption
	vx_addr := 7

	var max_ip uint16 = uint16(len(commands))
	for vm.ip < max_ip {
		// FIXME: time.Sleep(1 * time.Second)
		// Process the opcode
		command := commands[vm.ip]

		switch command {
		case 0x40: // CLD
			printLog(vm, "CLD", "")
			vm.fz = 0
			vm.fc = 0
			vm.ip = vm.ip + 1
		case 0x50: // LDX #IMM
			var value byte = byte(commands[vm.ip+1])
			printLog(vm, "LDX", strconv.FormatInt(int64(value), 16))
			vm.rx = value
			vm.ip = vm.ip + 2
		case 0x51: // LDY #IMM
			var value byte = byte(commands[vm.ip+1])
			printLog(vm, "LDY", strconv.FormatInt(int64(value), 16))
			vm.ry = value
			vm.ip = vm.ip + 2
		case 0x52: // STRX #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "STRX", strconv.FormatInt(int64(addr), 16))
			vm.memory[addr] = vm.rx
			vm.ip = vm.ip + 3
		case 0x53: // STRY #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "STRY", strconv.FormatInt(int64(addr), 16))
			vm.memory[addr] = vm.ry
			vm.ip = vm.ip + 3
		case 0x54: // LDRX #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "LDRX", strconv.FormatInt(int64(addr), 16))
			vm.rx = vm.memory[addr]
			vm.ip = vm.ip + 3
		case 0x55: // LDRY #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "LDRY", strconv.FormatInt(int64(addr), 16))
			vm.ry = vm.memory[addr]
			vm.ip = vm.ip + 3
		case 0x60: // OUT
			printLog(vm, "OUT", "")
			fmt.Printf("%02X", vm.rx)
			vm.ip = vm.ip + 1
		case 0x61: // IN
			var input byte = 0
			printLog(vm, "IN", "")
			fmt.Scanf("%02X\n", &input)
			vm.rx = input
			vm.ip = vm.ip + 1
		case 0x70: // CMPX #IMM
			var value byte = byte(commands[vm.ip+1])
			printLog(vm, "CMPX", strconv.FormatInt(int64(value), 16))
			if value == vm.rx {
				vm.fz = 1
				vm.fc = 0
			} else if vm.rx > value {
				vm.fz = 0
				vm.fc = 1
			} else {
				vm.fz = 0
				vm.fc = 0
			}
			vm.ip = vm.ip + 2
		case 0x71: // CMPY #IMM
			var value byte = byte(commands[vm.ip+1])
			printLog(vm, "CMPY", strconv.FormatInt(int64(value), 16))
			if value == vm.ry {
				vm.fz = 1
				vm.fc = 0
			} else if value > vm.ry {
				vm.fz = 0
				vm.fc = 1
			} else {
				vm.fz = 0
				vm.fc = 0
			}
			vm.ip = vm.ip + 2
		case 0x72: // JE #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "JE", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 1 {
				vm.ip = addr
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x73: // JRE #RELMEM
			addr := int16(binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3]))
			printLog(vm, "JRE", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 1 {
				vm.ip = vm.ip + uint16(addr)
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x74: // JL #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "JL", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 0 && vm.fc == 0 {
				vm.ip = addr
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x75: // JRL #RELMEM
			addr := int16(binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3]))
			printLog(vm, "JRL", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 0 && vm.fc == 0 {
				vm.ip = vm.ip + uint16(addr)
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x76: // JLE #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "JLE", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 1 || vm.fc == 0 {
				vm.ip = addr
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x77: // JRLE #RELMEM
			addr := int16(binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3]))
			printLog(vm, "JRLE", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 1 || vm.fc == 0 {
				vm.ip = vm.ip + uint16(addr)
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x78: // JG #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "JG", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 0 && vm.fc == 1 {
				vm.ip = addr
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x79: // JRG #RELMEM
			addr := int16(binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3]))
			printLog(vm, "JRG", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 0 && vm.fc == 1 {
				vm.ip = vm.ip + uint16(addr)
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x7A: // JGE #MEM
			addr := (binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3])) % MAXMEM
			printLog(vm, "JGE", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 1 || vm.fc == 1 {
				vm.ip = addr
			} else {
				vm.ip = vm.ip + 3
			}
		case 0x7B: // JRGE #RELMEM
			addr := int16(binary.BigEndian.Uint16(commands[vm.ip+1 : vm.ip+3]))
			printLog(vm, "JRGE", strconv.FormatInt(int64(addr), 16))
			if vm.fz == 1 || vm.fc == 1 {
				vm.ip = vm.ip + uint16(addr)
			} else {
				vm.ip = vm.ip + 3
			}
		case 0xA0: // ADDX #IMM
			var value byte = byte(commands[vm.ip+1])
			printLog(vm, "ADDX", strconv.FormatInt(int64(value), 16))
			oldrx := vm.rx
			vm.rx = vm.rx + value
			if vm.rx < oldrx || vm.rx < value {
				vm.fc = 1
			}
			vm.ip = vm.ip + 2
		case 0xA1: // ADDXY
			printLog(vm, "ADDXY", "")
			oldrx := vm.rx
			vm.rx = vm.rx + vm.ry
			if vm.rx < oldrx || vm.rx < vm.ry {
				vm.fc = 1
			}
			vm.ip = vm.ip + 1
		case 0xA2: // DECX #IMM
			var value byte = byte(commands[vm.ip+1])
			printLog(vm, "DECX", strconv.FormatInt(int64(value), 16))
			oldrx := vm.rx
			vm.rx = vm.rx - value
			if vm.rx > oldrx {
				vm.fc = 1
			}
			vm.ip = vm.ip + 2
		case 0xA3: // DECXY
			printLog(vm, "DECXY", "")
			oldrx := vm.rx
			vm.rx = vm.rx - vm.ry
			if vm.rx > oldrx {
				vm.fc = 1
			}
			vm.ip = vm.ip + 1
		case 0xA4: // RORX
			printLog(vm, "RORX", "")
			vm.rx = (vm.rx >> 1) | (vm.rx << (8 - 1))
			vm.ip = vm.ip + 1
		case 0xA5: // ROLX
			printLog(vm, "ROLX", "")
			vm.rx = (vm.rx << 1) | (vm.rx >> (8 - 1))
			vm.ip = vm.ip + 1
		case 0xA6: // XORX
			printLog(vm, "XORX", "")
			vm.rx = (vm.rx ^ vm.ry)
			vm.ip = vm.ip + 1
		case 0xB0: // PUSHX
			printLog(vm, "PUSHX", "")
			vm.stack = append(vm.stack, vm.rx)
			vm.ip = vm.ip + 1
		case 0xB1: // POPX
			printLog(vm, "POPX", "")
			if len(vm.stack) > 0 {
				vm.rx = vm.stack[len(vm.stack)-1]
				vm.stack = vm.stack[:len(vm.stack)-1]
			}
			vm.ip = vm.ip + 1
		case 0xB2: // PUSHY
			printLog(vm, "PUSHY", "")
			vm.stack = append(vm.stack, vm.ry)
			vm.ip = vm.ip + 1
		case 0xB3: // POPY
			printLog(vm, "POPY", "")
			if len(vm.stack) > 0 {
				vm.ry = vm.stack[len(vm.stack)-1]
				vm.stack = vm.stack[:len(vm.stack)-1]
			}
			vm.ip = vm.ip + 1
		case 0xC0: // RMEMX
			printLog(vm, "RMEMX", "")
			if len(vm.stack) > 1 {
				msb := byte(vm.stack[len(vm.stack)-2])
				lsb := byte(vm.stack[len(vm.stack)-1])
				addr := uint16(msb)<<8 | uint16(lsb)
				vm.rx = vm.memory[addr]
			}
			vm.ip = vm.ip + 1
		case 0xC1: // WMEMX
			printLog(vm, "WMEMX", "")
			if len(vm.stack) > 1 {
				msb := byte(vm.stack[len(vm.stack)-2])
				lsb := byte(vm.stack[len(vm.stack)-1])
				addr := uint16(msb)<<8 | uint16(lsb)
				vm.memory[addr] = vm.rx
			}
			vm.ip = vm.ip + 1
		case 0xC2: // RMEMY
			printLog(vm, "RMEMY", "")
			if len(vm.stack) > 1 {
				msb := byte(vm.stack[len(vm.stack)-2])
				lsb := byte(vm.stack[len(vm.stack)-1])
				addr := uint16(msb)<<8 | uint16(lsb)
				vm.ry = vm.memory[addr]
			}
			vm.ip = vm.ip + 1
		case 0xC3: // WMEMY
			printLog(vm, "WMEMY", "")
			if len(vm.stack) > 1 {
				msb := byte(vm.stack[len(vm.stack)-2])
				lsb := byte(vm.stack[len(vm.stack)-1])
				addr := uint16(msb)<<8 | uint16(lsb)
				vm.memory[addr] = vm.ry
			}
			vm.ip = vm.ip + 1
		case 0x90: // NOP
			printLog(vm, "NOP", "")
			vm.ip = vm.ip + 1
		case 0x91: // RET
			printLog(vm, "RET", "")
			return
		default:
			fmt.Println("Error: Unknown opcode -", command)
			return
		}

		// Check for KillSwitch
		if killed(vm.memory, kmemAddr, killSwitch) == true {
			fmt.Println()
			fmt.Println("### Virus has been deactivated   ###")
			fmt.Println("### Congratulations! Great Work! ###")
			return
		}

		if novirus == false {
			// VirusProcessing insert RET in program memory
			vm.memory[vx_addr] = 0x91
			if vx_addr < len(commands)-1 {
				commands[vx_addr] = 0x91
			}
			vx_addr = (vx_addr + 8) % 1024
		}
	}
}

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func main() {
	// Iterate through command-line arguments
	for _, arg := range os.Args[1:] {
		if arg == "debug" {
			debug = true
		} else if arg == "novirus" {
			novirus = true
		}
	}
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("##########################################")
	fmt.Println("# Welcooome to VM-O-MATIC v2             #")
	fmt.Println("#                                        #")
	fmt.Println("# Hacking is allowed ;-)                 #")
	fmt.Println("# ...but not recommended                 #")
	fmt.Println("##########################################")
	fmt.Println()

	// Generate Killswitch and Memory address
	killSwitch = generateRandomString(6)
	kmemAddr = rand.Intn(4096-1024+1) + 1024
	fmt.Println("WARNING: System Infected!!")
	fmt.Print("Identified KillSwitch (6 chars): ")
	fmt.Println(killSwitch)
	fmt.Print("Identified Memory Address for KillSwitch (decimal): ")
	fmt.Println(kmemAddr)
	fmt.Println("Virus is writing RET instruction every 8 bytes in program space")
	fmt.Println()
	if novirus == true {
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("!NOVIRUS mode selected! Use it only for testing!")
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println()
	}
	fmt.Println("Enter bytecode in hex format and hope for the best (e.g., CAFE0101):")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		bytecodeHex := scanner.Text()
		// Convert the input bytecode from hex to a byte slice
		bytecode, err := hex.DecodeString(bytecodeHex)
		if err != nil {
			fmt.Println("Error: Invalid bytecode input")
			return
		}

		// Initialize the VM
		vm := NewVM()

		// Copy executable in memory
		copy(vm.memory, bytecode)

		// Copy the secret in memory
		// Not used in this version
		// rnd := rand.Intn(2028-1+1) + 1
		// copy(vm.memory[1024+rnd:], secret)

		// Start processing executing code
		vm.ExecutePseudoCommands(bytecode)
	}
}
