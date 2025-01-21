### 0. Download and Unpack the Project

Download the provided zip file and unpack it. The project structure is already set up correctly, so no modifications are needed aside from adding the `.env` file in Step 1.

---

### 1. Preparing the Environment

1. **Install Go**  
   Follow the instructions at [https://go.dev/doc/install](https://go.dev/doc/install) to install Go.

2. **Install Node.js and npm**  
   TailwindCSS requires Node.js and npm. Install them by following the guide at [https://docs.npmjs.com/downloading-and-installing-node-js-and-npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).

3. **Install TailwindCSS**  
   Set up TailwindCSS by referring to its official installation guide: [https://tailwindcss.com/docs/installation](https://tailwindcss.com/docs/installation).

4. **Install Air**  
   Air is a live reloading tool for Go. Install it by following the instructions at [https://github.com/air-verse/air](https://github.com/air-verse/air).

5. **Install Tmux**  
   Tmux is required for managing terminal sessions. Use your package manager to install it. For example:

   ```bash
   sudo apt-get install -y tmux
   ```

6. **Create the `.env` File**  
   In the root directory of the project, create a file named `.env` with the following structure. Adjust the values based on your environment:

   ```env
   LISTEN_ADDR=":3000"                 # Application listen address
   DB_PATH="/path/to/your/database.db"  # Path to the database file
   LDAP_BASE_DN="dc=example,dc=com"    # LDAP domain settings
   LDAP_BIND_DN="cn=read-only-admin,dc=example,dc=com"
   LDAP_PORT="389"                     # LDAP server port
   LDAP_HOST="ldap.forumsys.com"       # LDAP server host
   LDAP_BIND_PASSWORD="password"       # LDAP bind password
   LDAP_FILTER="(uid=%s)"              # LDAP filter
   ```

7. **Update the `runLinux.sh` File**  
   Locate line 52 in the `runLinux.sh` script, which specifies the Go directory:
   ```bash
   "/usr/bin/go/go/bin/go build -tags=dev -o ./tmp/main ."
   ```
   Update the part before `build` to match the correct path to your Go installation.

---

### 2. Running the Project

1. **Install Dependencies**  
   In the root directory of the project, run the following command to download necessary dependencies:

   ```bash
   go mod tidy
   ```

2. **Run the Application**  
   Execute the `runLinux.sh` script to start the application:
   ```bash
   ./runLinux.sh
   ```
   - The script simplifies the process by managing terminal sessions and handling project instances.
   - If the first attempt to run the script reports an error, try running it again—it typically resolves the issue.
   - The script automatically terminates any previous instances of the project when re-run.

---
