; this is for an x64 system
.intel_syntax noprefix

.section .data
  sockaddr:             ; struct for SOCKET ADDRESS
    .word 2             ; AF_INET 
    .word 0x5000 ;80    ; port 80
    .long 0             ; INADDR_ANY
    .quad 0             ; padding
  
  client_addr:
    .space 16

  client_addr_len:
    .quad 16

  read_buf:
    .space 1024

  read_buf_len = . - read_buf

  static_buf:
    .ascii "HTTP/1.0 200 OK\r\n\r\n"

  static_buf_len = . - static_buf


.section .text
.global _start

_start:
  call create_socket
  call bind_socket
  call listen_socket
  call loop
  jmp exit

loop:
  call accept_socket
  call fork
  jmp loop

child_process:
  mov rdi, r12
  mov rax, 3
  syscall

request_loop:
  call read_static

  test rax, rax
  je child_exit          ; client closed connection

  cmp dword ptr [read_buf], 0x20544547    ; "GET "
  je handle_get

  cmp dword ptr [read_buf], 0x54534f50    ; "POST"
  je handle_post

  jmp request_loop

handle_get:
  call open_file
  call read_file
  call close_file
  call write_static
  call write_file
  call clear_buffer
  jmp request_loop 

handle_post:
  call open_file_write
  call find_body
  mov rbx, read_buf_len
  sub rbx, r15
  call write_body
  call close_file
  call write_static
  call clear_buffer
  jmp request_loop

;                                                 ===Sockets===

create_socket:
  mov rdi, 2
  mov rsi, 1
  xor rdx, rdx
  mov rax, 41
  syscall
  mov r12, rax  ; to save the fd (file descriptor)
  ret

bind_socket:
  
  ; Struct Sockaddr_in {                  
  ; sa_family_t     sin_family; AF_INET            2 bits
  ; in_port_t       sin_port; PORT 80 HTTP         2 bits written in BigEndian
  ; struct in_addr  sin_addr  IPV4 Address         4 bits
  ;
  ; BECAUSE WE WRITE TO 16 BITS AFTER THE 12 BITS BEING USED WE PAD THE REMAINING 6 BITS
  ;}
 
  mov rdi, r12
  lea rsi, [rip, + sockaddr] ; this gives us the sockaddr struct and loads it at rsi
  mov rdx, 16                ; size of struct 
  mov rax, 49
  syscall
  ret

listen_socket:
  mov rdi, r12    ; considering the order of function calls we don't really have to set this because its already set
  mov rsi, 3 
  xor rdx, rdx    ; rdx is NULL
  mov rax, 50     ; syscall listen_socket
  syscall
  ret

accept_socket:
  mov rdi, r12 
  xor rsi, rsi    ; rsi is NULL
  xor rdx, rdx    ; rdx is NULL
  mov rax, 43     ; syscall accept_socket
  syscall
  mov r13, rax    ; save client to fd
  ret

;                   ===IO Operation===

read_static:
  mov rdi, r13
  lea rsi, [rip + read_buf]
  mov rdx, read_buf_len
  mov rax, 0
  syscall
  ret

write_static:
  mov rdi, r13
  lea rsi, [rip + static_buf]
  mov rdx, static_buf_len
  mov rax, 1
  syscall
  ret

; after reading the appropriate file should open 
; right now we are taking the whole buffer but we just want the path
open_file:
  lea rsi, [rip + read_buf + 4]
  call find_space
  lea rdi, [rip + read_buf + 4] ; 4 OFFSET IS FOR GET THE OFFSET TO PARSE WHAT COMES AFTER THE 'GET' 
  mov rsi, 0 
  ;mov rdx, read_buf_len
  mov rax, 2 
  syscall 
  mov r14, rax  ; file fd (has to be saved after the syscall has been called and saved to rax)
  ret

find_space:
  cmp byte ptr [rsi], ' '
  je terminate
  inc rsi
  jmp find_space

terminate:
  mov byte ptr [rsi], 0
  ret

read_file:
  mov rdi, r14
  lea rsi, [rip + read_buf]
  mov rdx, read_buf_len
  mov rax, 0
  syscall
  mov r15, rax
  ret

write_file:
  mov rdi, r13
  lea rsi, [rip + read_buf]
  mov rdx, r15                ; bytes read from file
  mov rax, 1
  syscall
  ret

close_file:
  mov rdi, r14
  mov rax, 3
  syscall
  ret

close_socket:
  mov rdi, r13
  mov rax, 3
  syscall
  ret

open_file_write:
    lea rsi, [rip + read_buf + 5]   ; skip "POST "
    call find_space
    lea rdi, [rip + read_buf + 5]
    mov rsi, 0x241                  ; O_WRONLY|O_CREAT|O_TRUNC
    mov rdx, 0777
    mov rax, 2
    syscall
    mov r14, rax
    ret

write_body:
    mov rdi, r14        ; file fd
    mov rsi, r15        ; body ptr
    mov rdx, rbx        ; body length
    mov rax, 1
    syscall
    ret

;               ===Concurrency and Other Helpers===

find_body:
    lea rsi, [rip + read_buf]
.loop:
    cmp byte ptr [rsi], 0
    je .done
    cmp byte ptr [rsi], 0x0d
    jne .next
    cmp byte ptr [rsi+1], 0x0a
    jne .next
    cmp byte ptr [rsi+2], 0x0d
    jne .next
    cmp byte ptr [rsi+3], 0x0a
    jne .next
    add rsi, 4
    mov r15, rsi      ; body pointer
    ret
.next:
    inc rsi
    jmp .loop
.done:
    ret

fork:
  mov rax, 57
  syscall
  test rax, rax
  
  je child_process

  mov rdi, r13
  mov rax, 3
  syscall
  ret

child_exit:
    mov rdi, r13
    mov rax, 3
    syscall
    mov rdi, 0
    mov rax, 60
    syscall

clear_buffer:
  lea rdi, [rip + read_buf]
  mov rcx, read_buf_len
  xor rax, rax
  rep stosb
  ret

exit:
  mov rdi, 0
  mov rax, 60
  syscall
 

