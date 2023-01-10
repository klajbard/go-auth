# Go Auth Guard

Go service to wrap any amount of applications behind a simple auth guard

![Auth Guard](/images/auth-guard.png)

## Architecture

```
[ MongoDB ] <~ User credentials/Token validation ~> [ Backend ] <~ Token exchange ~> [ Frontend ]
```

## Backend

Written in Go

Packages used:
- [rs/cors](github.com/rs/cors)
- [julienschmidt/httprouter](github.com/julienschmidt/httprouter)
- [mgo.v2/bson](gopkg.in/mgo.v2/bson)

### Endpoints

- GET `/`
  - Serves `client/dist/index.html`
- GET `/logout`
  - Serves `client/dist/logout/index.html`
- GET `/services`
  - Serves `client/dist/services/index.html`
- GET `hello`
  - Creates/validates current session
- POST `/login`
  - Reads `username`, `password` and `rememberMe` fields from body
  - Authenticates user, set cookies if successful
- POST `/logout`
  - Clears refresh and session cookies

### Set up

Create an `.ENV` file containing:
- `MONGO_URL` secret to be able to communicate with MongoDB server
- `CLIENT_DEV_URL` secret to whitelist local url for CORS while developing

Example:
```sh
# .ENV
CLIENT_DEV_URL=http://127.0.0.1:5173
MONGO_URL=mongodb://127.0.0.1:27017/my_db
```

### Cookies
- `session_token`
  - Authenticates user
  - Valid for 2 hours
- `refresh_token`
  - Able to refresh `session_token`
  - Valid for 1 year

More about [Refresh tokens](https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/)

### Briefly about the session token validation process

- User has no session token, no refresh token
  - Get refresh token
  - Send 401

- User has no session token, has refresh token
  - Refresh token is valid
    - Create session token
    - Send 401

  - Refresh token is invalid or expired
    - Create new refresh token
    - Send 401

- User has session and refresh token
  - Refresh token is invalid
    - Create session token
    - Send 401

  - Refresh token is valid, session token is invalid or expired
    - Remove session token
    - Send 401

  - Refresh token and session token is valid
    - Send 200

## Frontend

Written in TypeScript

Packages used:
- [Vite](https://github.com/vitejs/vite)
- [Lit](https://github.com/lit/lit)

### Pages

The frontend is a multi-page app with the following pages:
- `/`
  - Root page
  - Contains the core logic
  - Calls `/hello` from the backend
  - Handles login by posting form data to `/login`
  - If login was successful redirects to the url recieved from the `?redirect=url` query parameter (url must be urlencoded!)
  - If login was successful and redirect query param is not present redirects to `/services`
- `/logout`
  - Handles logged out state
  - Redirects to home page (`/`) after 3 seconds
- `/404.html`
  - Handles not found from backend
- `/services`
  - Lists all available services configured to guard by the app

### Set up

Create an `.env` file containing `VITE_BACKEND_URL` secret to be able to communicate with the backend server. It is important to name it `VITE_*` otherwise vite won't recognize it by default.

Example:
```sh
# .env
VITE_BACKEND_URL=http://localhost:8080
```

## Integrate with guarded application

The application which needs to be guarded must call some endpoints from this auth service. For example in React app:

1. Create `AuthGuard` component

```tsx
// AuthGuard.tsx
import React, { FC } from "react";
import useAuth from "./useAuth";

const AuthGuard: FC = ({ children }) => {
  const { isLoading, isloggedIn } = useAuth();

  if (isLoading) {
    return <>Loading...</>;
  }

  if (!(isloggedIn || isLoading)) {
    window.location.assign(
      `${auth_server_url}/?redirect=${encodeURIComponent(
        window.location.href
      )}`
    );
    return null;
  }

  return <>{children}</>;
};

export default AuthGuard;
```

2. Create `useAuth` hook

```tsx
// useAuth.ts
import { useCallback, useEffect, useState } from "react";
import axios, { AxiosError } from "axios";

const useAuth = () => {
  const [isloggedIn, setIsLoggedIn] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const auth = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await axios(`${auth_server_url}/hello`, {
        withCredentials: true,
      });
      setIsLoggedIn(true);
    } catch (error) {
      const errorData = (error as AxiosError).response?.data;
      setIsLoggedIn(false);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    auth();
  }, [auth]);

  return { isLoading, isloggedIn };
};

export default useAuth;
```

3. Wrap application inside `AuthGuard`

```tsx
function App() {
  const logout = useCallback(() => {
    try {
      const response = await axios.post(`${auth_server_url}/logout`, null, {
        withCredentials: true,
      });
      window.location.assign(
        `${auth_server_url}/logout?redirect=${encodeURIComponent(
          window.location.href
        )}`
      );
    } catch (error) {
      const errorData = (error as AxiosError).response?.data;
    }
  })
  return (
    <AuthGuard>
        <button onClick={logout}>Logout</button>
    </AuthGuard>
  );
}
```