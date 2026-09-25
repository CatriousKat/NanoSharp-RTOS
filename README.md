# NanoSharp-RTOS
A port of NanoSharp to microcontrollers (eg. ESP32). <br>
NOTE: gui.* and fs.open doesn't work. <br>
<h1>Installation</h1>
Run install.ps1 or install.sh with your microcontroller plugged in. <br>
If it worked, it should print something like this: <br>
<br>
Connected to COM3. Press Ctrl-C to exit. <br>
ets Jul 29 2019 12:21:46 <br>
<br>
rst:0x1 (POWERON_RESET),boot:0x13 (SPI_FAST_FLASH_BOOT) <br>
configsip: 0, SPIWP:0xee <br>
clk_drv:0x00,q_drv:0x00,d_drv:0x00,cs0_drv:0x00,hd_drv:0x00,wp_drv:0x00 <br>
mode:DIO, clock div:1 <br>
load:0x3ffaf000,len:21572 <br>
load:0x3ffb4444,len:1408 <br>
ho 0 tail 12 room 4 <br>
load:0x40080000,len:44772 <br>
1150 mmu set 00010000, pos 00010000 <br>
entry 0x4008000c <br>
Starting NanoSharp on Microcontroller... <br>
Test <br>
Test contents 123 !@# <br>
Execution finished. <br>
<br>
(the script used print, and simple FS) <br>
