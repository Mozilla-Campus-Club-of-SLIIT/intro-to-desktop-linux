---
layout: statement
---

<h1 class="text-lemonYellow">Escape the GUI</h1>

---
layout: intro
---

<h1 class="text-warmRed">List Files</h1>

````md magic-move {lines: true}
```bash {1}
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> ls
```

```bash
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> ls
assets/
layouts/
package.json
slides/
snippets/
vite.config.ts
components/
node_modules/
pnpm-lock.yaml
slides.md
uno.config.ts
```
````

---
layout: intro
---

<h1 class="text-warmRed">Change Directory</h1>

````md magic-move {lines: true}
```bash {1}
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd slides
```

```bash
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd slides
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/s/slides (main)>
```

```bash
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd slides
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/s/slides (main)> ls
escape.md
freedom.md
history.md
roast_windows+mac_users.md
welcome.md
why_switching_2_linux.md
```
````

---
layout: intro
---

<h1 class="text-warmRed">View Files</h1>

````md magic-move {lines: true}
```bash {1}
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cat roast_windows+mac_users.md
```

```bash
mrbhanuka@sliitmozilla.org /m/h/D/D/P/m/i/slides (main)> cd roast_windows+mac_users.md
---
layout: intro-image-right
image: 'assets/mac-user.png'
---

# A Typical Mac User
`@grok Why is Mac better than Windows?`

---
layout: intro-video
video: 'assets/fuck_ms_fuck.webm'
---

<div class="absolute bottom-10 left-10 text-right">
    <h2>props to <a href="https://www.imdb.com/title/tt9612516/" class="text-lemonYellow">Space Force</a></h2>
</div>
```
````

---
layout: bullets
---

* `cp`
* `mv`
* `history`
* `whoami`
* `mkdir`
* `rm`
* `touch`
* `find`
* `help` me `man`
