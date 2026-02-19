---
layout: statement
---

<h1 class="text-warmRed">Linux File Architecture</h1>

---
layout: section
---

````md magic-move {lines: true}
```bash {1}
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd /
```

```bash
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd /
mrbhanuka@sliitmozilla.org / > ls
```

```bash
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd /
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

<h1 class="text-warmRed">View Files</h1>

````md magic-move {lines: true}
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
mrbhanuka@sliitmozilla.org / > ls home
mrbhanuka/
```
````
