---
layout: intro-image-right
image: 'assets/file-struct.webp'
---

<h1 class="text-warmRed">Linux File Architecture</h1>
````md magic-move {lines: true}
```bash 
mrbhanuka@sliitmozilla.org E/Melon_Eusk > 
```

```bash 
mrbhanuka@sliitmozilla.org E/Melon_Eusk > cd /
```

```bash
mrbhanuka@sliitmozilla.org E/Melon_Eusk > cd /
mrbhanuka@sliitmozilla.org / > 
```

```bash
mrbhanuka@sliitmozilla.org E/Melon_Eusk > cd /
mrbhanuka@sliitmozilla.org / > ls
```

```bash
mrbhanuka@sliitmozilla.org E/Melon_Eusk > cd /
mrbhanuka@sliitmozilla.org / > ls
bin@  boot/  dev/  etc/  home/  lib@  lib64@
lost+found/  mnt/  opt/  proc/  root/  run/
sbin@  srv/  sys/  tmp/  usr/  var/
```

````
---
layout: intro
---

<h1 class="text-warmRed">Binaries</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{2}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash
mrbhanuka@sliitmozilla.org / > ls /bin | grep -E '^(ls|mkdir|rm|mv)$'
ls
mkdir
mv
rm
```

```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{18}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

````
---
layout: intro
---

<h1 class="text-warmRed">Shared Libraries</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{14}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

```bash{5}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
````

---
layout: intro
---

<h1 class="text-warmRed">
  <span v-if="$slidev.nav.clicks < 5">
    <span v-mark.strike.red="4">User System Resources</span>
  </span>
</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
````

---
layout: intro
---

<h1 class="text-warmRed">
    UNIX System Resources
</h1>

````md magic-move {lines: true}
```bash{20}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash
mrbhanuka@sliitmozilla.org / > ls /usr/
```
```bash
mrbhanuka@sliitmozilla.org / > ls /usr/
bin/  include/  lib/  lib32/  lib64@  local/  sbin@  share/  src/
```
````


---
layout: statement
---

<h1 class="text-lemonYellow">$PATH</h1>


---
layout: intro
---

<h1 class="text-warmRed">Editable Text Configuration</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{13}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /etc | grep -E '^(machine-id|hostname|hosts)$'
```

```bash
mrbhanuka@sliitmozilla.org / > ls /etc | grep -E '^(machine-id|hostname|hosts)$'
hostname
hosts
machine-id
```
```bash
mrbhanuka@sliitmozilla.org / > cat /etc/hostname
```

```bash
mrbhanuka@sliitmozilla.org / > cat /etc/hostname
sliitmozilla.org
```

```bash
mrbhanuka@sliitmozilla.org / > cat /etc/hosts
```

```bash
mrbhanuka@sliitmozilla.org / > cat /etc/hosts
127.0.0.1   sliitmozilla.org
::1         sliitmozilla.org
```
````

---
layout: intro
---

<h1 class="text-warmRed">HOME 🏠</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{4}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /home
```

```bash
mrbhanuka@sliitmozilla.org / > ls /home
mrbhanuka/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /home/mrbhanuka
```

```bash
mrbhanuka@sliitmozilla.org / > mrbhanuka@aspire3 /> ls /home/mrbhanuka/
'2026-02-15 14-10-41.mp4'  '2026-02-16 00-00-26.mp4'   Desktop@     etc@        Public/
'2026-02-15 14-18-09.mp4'  '2026-02-16 00-01-03.mp4'   Documents@   go/         Templates/
'2026-02-15 22-31-38.mp4'  '2026-02-16 00-02-03.mp4'   dotfiles/    Music@      Videos@
'2026-02-15 22-35-11.mp4'  '2026-02-16 00-03-16.mp4'   Downloads@   Pictures@
```
````

---
layout: intro
---

<h1 class="text-warmRed">boot 👢</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{12}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /boot
```

```bash
mrbhanuka@sliitmozilla.org / > ls /boot
EFI/  grub/  initramfs-linux.img*  intel-ucode.img*  vmlinuz-linux*
```
````

---
layout: intro
---

<h1 class="text-warmRed">Variables 👩‍🦱</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{11}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /var
```

```bash
mrbhanuka@sliitmozilla.org / > ls /var
cache/  db/  empty/  games/  lib/  local/  lock@  log/  mail@  opt/  run@  spool/  tmp/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /var/log
```

```bash
mrbhanuka@sliitmozilla.org / > ls /var/log
audit/  dnscrypt-proxy@  lastlog   old/        private/  swtpm/
btmp    journal/         libvirt/  pacman.log  README@   wtmp
```
````

---
layout: intro
---

<h1 class="text-warmRed">Devices</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > ls
```
```bash
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```
```bash{3}
mrbhanuka@sliitmozilla.org / > ls
bin@
dev/
home/
lib64@
mnt/
proc/
run/
srv/
tmp/
var/
boot/
etc/
lib@
lost+found/
opt/
root/
sbin@
sys/
usr/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /dev | grep -E '^(cpu|stdin|stdout|stdout|)$|nvme'
```

```bash
mrbhanuka@sliitmozilla.org / > ls /dev | grep -E '^(cpu|stdin|stdout|stdout|)$|nvme'
cpu
nvme0
nvme0n1
nvme0n1p1
nvme0n1p2
stdin
stdout
```

```bash
mrbhanuka@sliitmozilla.org / > ls /dev/cpu/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /dev/cpu/
0/  1/  2/  3/  4/  5/  6/  7/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /dev/cpu/0/
```

```bash
mrbhanuka@sliitmozilla.org / > ls /dev/cpu/0/
cpuid  msr
```

```bash
mrbhanuka@sliitmozilla.org / > cat /dev/cpu/0/cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > cat /dev/cpu/0/cpuid
cat: /dev/cpu/0/cpuid: Permission denied
```
````

---
layout: intro-image-right
image: 'assets/sudo.webp'
---

<h1 class="text-lemonYellow">The god mode 🛡️</h1>
<h2 class="text-warmRed">SUDO:Superuser Do</h2>

---
layout: section
---

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > sudo cat /dev/cpu/0/cpuid
```
```bash
mrbhanuka@sliitmozilla.org / > sudo cat /dev/cpu/0/cpuid
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d
$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d$d^C⏎             

mrbhanuka@sliitmozilla.org /  [SIGINT]>
```
```bash
mrbhanuka@sliitmozilla.org / > cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > cpuid
fish: Unknown command: cpuid
```
````
