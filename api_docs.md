# Auth and Addresses API

This collection covers all authentication and address endpoints for the auth-service backend.
Base URL: `http://localhost:8080`

---

## Auth

### Register

**POST** `http://localhost:8080/register`

Creates a new user.

**Request Body**

```json
{
  "username": "Ana",
  "email": "Ana@example.com",
  "password": "supersecret",
  "repeatedPassword": "supersecret"
}
```

**Description**

- **username**: 3–30 chars
- **email**: valid email
- **password**: min 8 chars
- **repeatedPassword**: must match password

**Example Response** (201 Created)

```json
{
  "user": {
    "id": "8300ffff-f2c0-4e1b-b5fd-1700abd17c5c",
    "username": "Ana"
  }
}
```

---

### Login

**POST** `http://localhost:8080/login`

Validates user credentials and issues an HttpOnly JWT cookie.

**Request Body**

```json
{
  "email": "ana@example.com",
  "password": "supersecret"
}
```

**Example Response** (200 OK)

- Sets cookie:

  ```
  Set-Cookie: access_token=<jwt>; Path=/; Expires=Thu, 24 Jul 2025 11:17:15 GMT; HttpOnly; Secure; SameSite=Lax
  ```

```json
{
  "user": {
    "id": "8300ffff-f2c0-4e1b-b5fd-1700abd17c5c",
    "username": "Ana"
  }
}
```

---

### Logout

**POST** `http://localhost:8080/api/v1/logout`

Clears the JWT cookie, logging out the user.

**Headers**

```
X-CSRF-Token: pXWgYHpoJcSkVcSevqqAoKkRukWrYbjb
```

**Example Response** (204 No Content)

- Clears cookies:

  ```
  Set-Cookie: csrf_token=pXWgYHpoJcSkVcSevqqAoKkRukWrYbjb; Expires=Thu, 24 Jul 2025 11:17:38 GMT; SameSite=Strict
  Set-Cookie: access_token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure; SameSite=Lax
  ```

---

## Address

### Create Address

**POST** `http://localhost:8080/api/v1/users/address/add`

Creates a new address for the authenticated user.

**Headers**

```
X-CSRF-Token: JtOfePwnHpOUMHOrVngnLeCrPmigwYqR
```

**Request Body**

```json
{
  "addr_1": "Burgemeister str. 50",
  "addr_2": "",
  "zip": "10115",
  "city": "Berlin",
  "country": "Germany",
  "isdefault": false
}
```

---

### Get Address

**GET** `http://localhost:8080/api/v1/users/address/2`

Retrieves the address with ID 2 for the authenticated user.

---

### Delete Address

**DELETE** `http://localhost:8080/api/v1/users/address/1`

Deletes the address with ID 1 for the authenticated user.

---

### Update Address

**PATCH** `http://localhost:8080/api/v1/users/address/2`

Updates the address with ID 2 for the authenticated user.

**Headers**

```
X-CSRF-Token: JtOfePwnHpOUMHOrVngnLeCrPmigwYqR
```

**Request Body**

```json
{
  "addr_1": "Bouchestr 52",
  "addr_2": "",
  "zip": "10115",
  "city": "Berlin",
  "country": "Germany",
  "isdefault": false
}
```
