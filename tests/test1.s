.org 0

ldlit $i0, 0x40
sind $i0

ldlit $r1b, 0x01
soff $r1b

ldind $r2s, 0x01

jmp $r2

.fillto 0x40, 0
.short 0x0
.short 0xD0

.fillto 0xD0, 0
hlt