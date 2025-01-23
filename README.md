For more detailed documentation regarding installation and usage navigate to documentation folder.

**Environment setup**

1. install go

- [https://go.dev/doc/install]

2. install node.js and npm (needed for tailwind)

- [https://docs.npmjs.com/downloading-and-installing-node-js-and-npm]

3. install tailwindcss

[https://tailwindcss.com/docs/installation]

4. install air

- [https://github.com/air-verse/air]

5. Create .env file and add LISTEN_ADDR=":3000", DB_PATH="/path/to/your/database.db" and variables related to ldap to it.

**Running the app for developement**

Run those commands in separate terminals or by adding `&` at the end of command to run it in the background.
All 3 commands must be running to enable live reload of the app.

```bash
make css
```

```bash
air
```

```bash
templ generate --watch --proxy=http://localhost:3000
```

Or run:

```bash
./runLinux.sh
```

if you are running it on Linux.

Now the live reloading app should be available on port 7331. To navigate to it just type http://localhost:7331 in your search bar.

**docs for technologies**

go standard library

[https://pkg.go.dev/std]

go chi

[https://go-chi.io/#/README]

go templ

[https://templ.guide/]

htmx

[https://htmx.org/docs/#introduction]

tailwindcss

[https://tailwindcss.com/docs/installation]

[https://nerdcave.com/tailwind-cheat-sheet]
