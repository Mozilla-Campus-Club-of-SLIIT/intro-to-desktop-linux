---
layout: intro-image-right
image: 'assets/pacman.webp'
---


<h1 class="text-lemonYellow">Play/App Store</h1>

## Let's download a spyware called "Google Chrome"

---
layout: statement
---

<h1 class="text-warmRed">
  <a href="https://music.youtube.com/watch?v=UytDoxG5kFA&si=XzSP9ft9MNgQqzp6">Windows(aka redStartOs for US) vs Linux</a>
</h1>

---
layout: image-right
image: 'assets/p-pacman.webp'
---

<h1 class="text-warmRed">Pacman 👻</h1>

* dnf
* apt
* snap 🤮
* flatpak

---
layout: intro
---

<h1 class="text-warmRed">Let's install `cpuid`</h1>

````md magic-move {lines: true}
```bash
mrbhanuka@sliitmozilla.org / > pacman -Ss cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > pacman -Ss cpuid
extra/kcpuid 6.19-2 (linux-tools)
    Kernel tool for various cpu debug outputs
extra/libcpuid 0.8.1-1
    A small C library for x86 CPU detection and feature extraction
```

```bash
mrbhanuka@sliitmozilla.org / > paru -Ss cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > paru -Ss cpuid
extra/kcpuid 6.19-2 [27.58 KiB 115.74 KiB] (linux-tools)
    Kernel tool for various cpu debug outputs
extra/libcpuid 0.8.1-1 [115.18 KiB 425.93 KiB]
    A small C library for x86 CPU detection and feature extraction
aur/cpuid 20250513-1 [+79 ~0.09]
    Linux tool to dump x86 CPUID information about the CPU(s)
aur/libcpuid-git 2:0.8.0.r0.g5bb7c32-1 [+68 ~0.01]
    A small C library for x86 CPU detection and feature extraction
aur/libicuid 1.4.1-1 [+9 ~0.00] [Out-of-date: 2024-01-12] [Orphaned]
    I C U ID is a modern library that provides a C interface for the CPUID
    opcode
aur/cpuid2cpuflags 17-1 [+7 ~0.00]
    Tool to get the instruction sets supported by the local CPU
aur/cpuid.py-git r27.fcffb0a-2 [+1 ~0.00]
    Pure Python interface to the CPUID instruction
aur/ddcpuid 0.14.1-1 [+1 ~0.00] [Out-of-date: 2024-01-08]
    dd's x86 CPU Identification tool
aur/cpuid2cpuflags-git 15.r0.g4c6aedf-1 [+0 ~0.00]
    Tool to get the instruction sets supported by the local CPU (git version)
```

```bash
mrbhanuka@sliitmozilla.org / > paru -S cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > paru -S cpuid
:: Resolving dependencies...
:: Calculating conflicts...
:: Calculating inner conflicts...

Aur (1) cpuid-20250513-1

:: Proceed to review? [Y/n]: Y

:: Downloading PKGBUILDs...
 (1/1) downloading: cpuid-20250513-1
:: cpuid:
...
```

```bash
mrbhanuka@sliitmozilla.org / > cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > cpuid
CPU 0:
   vendor_id = "GenuineIntel"
   version information (1/eax):
      processor type  = primary processor (0)
      family          = 0x6 (6)
      model           = 0xe (14)
      stepping id     = 0xc (12)
      extended family = 0x0 (0)
      extended model  = 0x8 (8)
      (family synth)  = 0x6 (6)
      (model synth)   = 0x8e (142)
...
```

```bash
mrbhanuka@sliitmozilla.org / > paru -R cpuid
```

```bash
mrbhanuka@sliitmozilla.org / > paru -S cpuid
[sudo] password for mrbhanuka:
checking dependencies...

Packages (1) cpuid-20250513-1

Total Removed Size:  0.41 MiB

:: Do you want to remove these packages? [Y/n] Y
:: Processing package changes...
(1/1) removing cpuid                               [####################] 100%
:: Running post-transaction hooks...
(1/1) Arming ConditionNeedsUpdate...
```
````
