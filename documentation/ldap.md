## LDAP Authentication Documentation

### Overview

Our application uses the [go-guardian](https://pkg.go.dev/github.com/shaj13/go-guardian/v2) library to handle LDAP authentication. This enables integration with LDAP servers for user authentication and role-based access control.

---

### Configuration

To enable LDAP authentication, the following variables must be correctly declared in the `.env` file:

```env
LDAP_BASE_DN="dc=example,dc=com"       # LDAP domain settings
LDAP_BIND_DN="cn=read-only-admin,dc=example,dc=com"
LDAP_PORT="389"                        # LDAP server port
LDAP_HOST="ldap.forumsys.com"          # LDAP server host
LDAP_BIND_PASSWORD="password"          # LDAP bind password
LDAP_FILTER="(uid=%s)"                 # LDAP filter
```

### Admin and User Assignment Logic

#### Admin Functionality

- **Current Behavior:**  
  Administrators are hardcoded in the `isAdmin` function located in `handlers/auth.go`. For now, only the user with the username `"tesla"` is recognized as an admin.
- **Customizing Admin Roles:**
  - To use LDAP groups for admin determination, uncomment the provided group-based logic in the `isAdmin` function.
  - Update the group name `"admin"` to match the appropriate admin group in your LDAP configuration.

#### User Account Assignment

- **Standard User Mapping:**  
  Non-admin users are assigned to accounts based on their name. The username in LDAP must match the `"first_name last_name"` fields in the application's database.
  - If no match is found, the user will be logged into an empty account without access to any theses or related data.
- **Customizing User Mapping Logic:**  
  To modify how users are mapped to accounts, refer to the `getUserIdFromUsername(name string)` function in `handlers/middleware.go`.

---
