# Module 15: Authentication with JWT

## Table of Contents

<ol>
  <li><a href="#introduction-to-authentication-in-web-applications">Introduction to Authentication in Web Applications</a></li>
  <li><a href="#authentication-fundamentals">Authentication Fundamentals</a></li>
  <li><a href="#jwt-as-an-authentication-mechanism">JWT as an Authentication Mechanism</a></li>
  <li><a href="#getting-started-with-jwt-in-go">Getting Started with JWT in Go</a></li>
  <li><a href="#implementing-jwt-authentication-in-go">Implementing JWT Authentication in Go</a></li>
  <li><a href="#jwt-integration-with-gin-and-echo">JWT Integration with Gin and Echo</a></li>
  <li><a href="#advanced-jwt-techniques">Advanced JWT Techniques</a></li>
</ol>

## Introduction to Authentication in Web Applications

### What authentication is and why it is required
Authentication is the process of verifying the identity of a user or system before granting access to protected resources.

- Ensures only legitimate users access protected APIs
- Protects sensitive data and operations
- Establishes trust between client and server

### Authentication in web and API-based systems
- Web apps often use cookies and sessions
- APIs typically rely on token-based authentication
- APIs must be stateless and scalable

### Typical authentication flow in backend services
1. Client sends credentials (e.g., username/password)
2. Server validates credentials
3. Server issues an authentication token
4. Client sends token with subsequent requests
5. Server validates token on each request

---

## Authentication Fundamentals

### Identity and credentials
- **Identity**: who the user is
- **Credentials**: proof of identity (password, token, key)

### Authentication vs session management
- Authentication: verifying identity
- Session management: maintaining authenticated state

### Stateless vs stateful authentication
- **Stateful**: server stores session data
- **Stateless**: client carries authentication data

### Session-based authentication
- Server stores session in memory or database
- Session ID stored in cookies
- Not ideal for distributed systems

### Token-based authentication
- Token contains authentication information
- No server-side session storage required
- Scales well for APIs

### API keys
- Simple token tied to an application
- Not suitable for user authentication

### Why token-based authentication is common for APIs
- Stateless
- Scalable
- Works well with mobile and SPA clients

---

## JWT as an Authentication Mechanism

### What JWT is and what problem it solves
JWT (JSON Web Token) is a compact, URL-safe token format used to represent claims securely between parties.

### JWT vs traditional sessions

| Sessions        | JWT           |
|-----------------|---------------|
| Stateful        | Stateless     |
| Cookie-based    | Header-based  |
| Harder to scale | Easy to scale |

---

## Getting Started with JWT in Go

### JWT structure
A JWT consists of three parts separated by dots (.)
- Header
- Payload
- Signature

Sample format
```
header.payload.signature
```

**Header**: Typically consists of two parts: the type of the token, which is JWT, and the signing algorithm being used, such as HMAC SHA256 or RSA.
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```
**Payload**: The second part of the token is the payload, which contains the claims. Claims are statements about an entity (typically user) and additional data. 

There are three types of claims: registered, public, and private claims.

- **Registered claims**: These are a set of predefined claims which are not mandatory but recommended, to provide a set of useful, interoperable claims. Some of them are: iss (issuer), exp (expiration time), sub (subject), aud (audience), and others.
Notice that the claim names are only three characters long, as JWT is meant to be compact.

- **Public claims**: These can be defined at will by those using JWTs. But to avoid collisions, they should be defined in the [IANA JSON Web Token Registry](https://www.iana.org/assignments/jwt/jwt.xhtml) or be defined as a URI that contains a collision-resistant namespace.

- **Private claims**: These are the custom claims created to share information between parties that agree on using them and are neither registered or public claims.
```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```
**Signature**: The signature is used to verify the message wasn't changed along the way, if tokens signed with a private key, it can also verify that the sender of the JWT is correct

To create the signature part you have to take the encoded header, the encoded payload, a secret, the algorithm specified in the header, and sign that.

For example, with HMAC SHA256 algorithm, the signature will be created in the following way:
```
HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)
```

### Access tokens vs refresh tokens
Modern authentication systems typically use two different types of tokens to balance security, performance, and user experience: access tokens and refresh tokens. Although both are often implemented as JWTs, they serve very different purposes.

#### Access Token

An access token is a short-lived credential used to authenticate API requests.

- A token that represents the authenticated user
- Sent with every protected API request
- Usually implemented as a JWT
- Contains identity information (e.g., user ID)
- Short expiration time (5–15 minutes)
- Stateless (server does not store it)
- Verified by checking:
  - Signature
  - Expiration (exp)
  - Other claims

```
GET /api/v1/protected-endpoint
Authorization: Bearer <access_token>
```

#### When it expires
- Server returns 401 Unauthorized
- Client must request a new access token using a refresh token

#### Refresh Token

A refresh token is a long-lived credential used to obtain new access tokens without requiring the user to log in again.

- Used only with a refresh endpoint (e.g., /refresh)
- Never sent to protected APIs
- Must be tracked server-side
- Long expiration time (days or weeks)
- Stateful (stored in database or session store)
- Often includes a unique identifier (jti)

#### Why it exists
- Improves user experience by avoiding frequent logins
- Allows token rotation and revocation
- Enables logout in stateless systems

---

### JWT Core Components in Go

Example: generating a JWT token
```shell
go get -u github.com/golang-jwt/jwt/v5
```

```go

import (
    "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
```

Parsing and validating tokens
```go
func ParseToken(tokenStr string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
```

### Handling token validation errors
- Invalid signature → unauthorized
- Expired token → re-authentication required
- Malformed token → bad request


## Implementing JWT Authentication in Go
### Login and Token Generation
- Validate credentials
- Generate JWT
- Return token in response

```go
package main

import (
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

// Secret key for signing tokens
var jwtKey = []byte("your_secret_key")

// User model for demonstration
type User struct {
  ID       string `json:"id"`
  Username string `json:"username"`
  Password string `json:"password"` // In production, store hashed passwords only
  Role     string `json:"role"`
}

// Mock user database
var users = map[string]User{
  "johndoe": {
    ID:       "1",
    Username: "johndoe",
    Password: "password123", // Never store plain passwords in real apps
    Role:     "user",
  },
  "admin": {
    ID:       "2",
    Username: "admin",
    Password: "admin123",
    Role:     "admin",
  },
}

type Credentials struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

func GenerateToken(userID string, secret []byte) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
    "exp": time.Now().Add(15 * time.Minute).Unix(),
    "iat": time.Now().Unix(),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(secret)
}

func Login(c *gin.Context) {
  var creds Credentials

  // Bind JSON request to credentials struct
  if err := c.ShouldBindJSON(&creds); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
    return
  }

  // DEMO PURPOSE ONLY:
  // In this example, users are stored in an in-memory map and passwords
  // are compared directly for simplicity.
  //
  // In a real application:
  // - Users would be queried from a database
  // - Passwords would be stored as hashes (bcrypt/argon2)
  // - Password comparison would use a secure hash comparison
  // - No plaintext passwords would ever be stored in memory
  user, exists := users[creds.Username]
  if !exists || user.Password != creds.Password {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
    return
  }

  token, _ := GenerateToken("123", jwtKey)
  c.JSON(http.StatusOK, gin.H{
    "access_token": token,
  })
}
```

Token Verification Middleware
- Extract token from request
- Validate token
- Attach identity to request context
- Extracting tokens from HTTP headers

```go
var jwtKey = []byte("your_secret_key")
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Read Authorization header
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 2. Extract token
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		// 3. Parse and validate token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 4. Extract user identity
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["sub"].(string)

		// 5. Attach identity to Gin context
		c.Set("userID", userID)

		// Continue to handler
		c.Next()
	}
}

func ProtectedHandler(c *gin.Context) {
    userID := c.GetString("userID")

    c.JSON(http.StatusOK, gin.H{
        "message": "Hello user " + userID,
    })
}


func main(){
  r := gin.Default()
  protected := r.Group("/api")
  protected.Use(JWTMiddleware()) // Apply middleware to group
  {
    protected.GET("/protected", ProtectedHandler)
  }
}
```

## JWT Integration with Gin and Echo
- Gin JWT Middleware 
```go
func ParseToken(tokenStr string, secret []byte) (*jwt.Token, error) {
    return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Ensure the signing method is what we expect
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return secret, nil
    })
}


func JWTMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := ParseToken(tokenStr, secret)
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
```

- Echo JWT Middleware
```go
func ParseToken(tokenStr string, secret []byte) (*jwt.Token, error) {
    return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Ensure the signing method is what we expect
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return secret, nil
    })
}

func JWTMiddleware(secret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				return echo.ErrUnauthorized
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			token, err := ParseToken(tokenStr, secret)
			if err != nil || !token.Valid {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
```

Echo built-in JWT middleware

```go
import "github.com/labstack/echo-jwt/v5"

e.Use(echojwt.JWT([]byte("secret")))
//or with custom config
e.Use(echojwt.WithConfig(echojwt.Config{
	SigningKey:             []byte("secret"),
	// ...
}))
```

Protecting routes using middleware
```go
r.GET("/protected", JWTMiddleware(secret), handler)
```

## Advanced JWT Techniques
### Token renewal flow using refresh token
Access tokens expire quickly (e.g., 5–15 minutes) while refresh tokens allow renewal without re-login (e.g., 7–30 days)

- Client sends refresh token
- Server validates refresh token
- New access token is issued

![rf_token_flow.png](images/rf_token_flow.png)

```go
package main

import (
  "errors"
  "net/http"
  "strings"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
)

var secret = []byte("super-secret")

// --- In-memory refresh session store ---
// refreshJTI -> active
var refreshSessions = map[string]bool{}

func newAccessToken(userID string) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
    "typ": "access",
    "exp": time.Now().Add(5 * time.Minute).Unix(),
    "iat": time.Now().Unix(),
  }
  return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func newRefreshToken(userID, jti string) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
    "typ": "refresh",
    "jti": jti,
    "exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
    "iat": time.Now().Unix(),
  }
  return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
  t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
    return secret, nil
  })
  if err != nil || !t.Valid {
    return nil, errors.New("invalid token")
  }
  claims, ok := t.Claims.(jwt.MapClaims)
  if !ok {
    return nil, errors.New("invalid claims")
  }
  return claims, nil
}

func RefreshHandler(c *gin.Context) {
  var req struct {
    RefreshToken string `json:"refresh_token"`
  }
  if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
    return
  }

  claims, err := parseToken(req.RefreshToken)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
    return
  }

  // Check token type refresh token
  if claims["typ"] != "refresh" {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "not a refresh token"})
    return
  }

  userID, _ := claims["sub"].(string)
  jti, _ := claims["jti"].(string)
  if userID == "" || jti == "" {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh claims"})
    return
  }

  // Check server-side session state
  if !refreshSessions[jti] {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token revoked"})
    return
  }

  // Revoke old refresh token, issue new tokens
  refreshSessions[jti] = false
  newJTI := uuid.New()
  refreshSessions[newJTI] = true

  //New access token
  accessToken, err := newAccessToken(userID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue access token"})
    return
  }

  //New refresh token
  refreshToken, err := newRefreshToken(userID, newJTI)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue refresh token"})
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "access_token":  accessToken,
    "refresh_token": refreshToken,
  })
}
```

### Token Revocation and Logout
JWT access tokens are stateless, so logout usually targets refresh tokens (session layer)

#### Token blacklisting
A blacklist approach stores revoked token IDs and checks them on each request.
**Note**: This requires the access token to include a `jti` claim.
- The blacklist should store a TTL (e.g., until token exp) and be concurrency-safe.
- In-memory blacklist only works per instance (in real apps, use Redis).

**Flow**:
- On each request: validate JWT → extract jti → reject if blacklisted.
- On logout: blacklist the current access token jti until it expires.

```go
// For demo will use in app memory
package main

import (
  "net/http"
  "strings"
  "sync"
  "time"

  "github.com/gin-gonic/gin"
)
var secret = []byte("your-secret-key")

// Store jti -> expUnix (TTL by exp)
var (
  accessBlacklistMu sync.RWMutex
  accessBlacklist   = map[string]int64{} // jti => expUnix
)

// isBlacklisted returns true if jti exists and not expired.
// Also removes expired entries.
func isBlacklisted(jti string) bool {
  if jti == "" {
    return false
  }

  now := time.Now().Unix()

  accessBlacklistMu.RLock()
  expUnix, ok := accessBlacklist[jti]
  accessBlacklistMu.RUnlock()

  if !ok {
    return false
  }

  // Expired blacklist entry => cleanup
  if expUnix <= now {
    accessBlacklistMu.Lock()
    // re-check then delete
    if exp2, ok2 := accessBlacklist[jti]; ok2 && exp2 <= now {
      delete(accessBlacklist, jti)
    }
    accessBlacklistMu.Unlock()
    return false
  }

  return true
}

func blacklistAccessJTI(jti string, expUnix int64) {
  if jti == "" || expUnix == 0 {
    return
  }
  accessBlacklistMu.Lock()
  accessBlacklist[jti] = expUnix
  accessBlacklistMu.Unlock()
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
  t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
    // IMPORTANT: enforce expected algorithm
    if t.Method != jwt.SigningMethodHS256 {
      return nil, errors.New("unexpected signing method")
    }
    return secret, nil
  })
  if err != nil || !t.Valid {
    return nil, errors.New("invalid token")
  }

  claims, ok := t.Claims.(jwt.MapClaims)
  if !ok {
    return nil, errors.New("invalid claims")
  }
  return claims, nil
}

func JWTMiddlewareWithBlacklist() gin.HandlerFunc {
  return func(c *gin.Context) {
    auth := c.GetHeader("Authorization")
    if !strings.HasPrefix(auth, "Bearer ") {
      c.AbortWithStatus(http.StatusUnauthorized)
      return
    }

    tokenStr := strings.TrimPrefix(auth, "Bearer ")
    claims, err := parseToken(tokenStr)
    if err != nil {
      c.AbortWithStatus(http.StatusUnauthorized)
      return
    }

    if claims["typ"] != "access" {
      c.AbortWithStatus(http.StatusUnauthorized)
      return
    }

    jti, _ := claims["jti"].(string)
    if isBlacklisted(jti) {
      c.AbortWithStatus(http.StatusUnauthorized)
      return
    }

    userID, _ := claims["sub"].(string)
    if userID == "" {
      c.AbortWithStatus(http.StatusUnauthorized)
      return
    }
    c.Set("userID", userID)
    c.Next()
  }
}

// LogoutBlacklistAccessToken blacklists current access token
func LogoutBlacklistAccessToken(c *gin.Context) {
  auth := c.GetHeader("Authorization")
  if !strings.HasPrefix(auth, "Bearer ") {
    c.JSON(http.StatusBadRequest, gin.H{"error": "missing access token"})
    return
  }

  tokenStr := strings.TrimPrefix(auth, "Bearer ")
  claims, err := parseToken(tokenStr)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
    return
  }

  if claims["typ"] != "access" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
    return
  }

  jti, _ := claims["jti"].(string)
  if jti == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "access token missing jti"})
    return
  }

  expFloat, _ := claims["exp"].(float64) // jwt.MapClaims decodes numbers as float64
  expUnix := int64(expFloat)
  if expUnix <= time.Now().Unix() {
    // already expired, nothing to blacklist
    c.Status(http.StatusNoContent)
    return
  }

  blacklistAccessJTI(jti, expUnix)
  c.Status(http.StatusNoContent)
}

```

#### Stateful revocation
Refresh tokens should be stateful via a session store (DB/Redis). Each refresh token contains a `jti` (session id).
**Flow:**
- Login creates a refresh session and issues token pair.
- Refresh validates refresh token, checks session active, then rotates:
  - revoke old session
  - create new session (new jti)
  - issue new token pair
- Logout revokes a specific refresh session.

**Notes:**
- If you only revoke refresh/session, access token remains valid until expired
- For instant logout, revoke refresh/session and blacklist the current access token jti.

```go
package main

import (
  "errors"
  "net/http"
  "sync"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
)

var secret = []byte("your-secret-key")

// DEMO purpose: in-memory session table
// In real apps: store in Postgres/Redis
var (
  sessionMu    sync.RWMutex
  sessionTable = map[string]RefreshSession{} // refresh_jti => session
)

type RefreshSession struct {
  JTI       string
  UserID    string
  RevokedAt *time.Time
  ExpiresAt time.Time
}

func isActive(session RefreshSession) bool {
  if session.RevokedAt != nil {
    return false
  }
  if time.Now().After(session.ExpiresAt) {
    return false
  }
  return true
}

func issueAccessToken(userID string) (string, error) {
  accessJTI := uuid.NewString()
  claims := jwt.MapClaims{
    "sub": userID,
    "typ": "access",
    "jti": accessJTI,
    "exp": time.Now().Add(5 * time.Minute).Unix(),
    "iat": time.Now().Unix(),
  }
  return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func issueRefreshToken(userID string, jti string) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
    "typ": "refresh",
    "jti": jti,
    "exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
    "iat": time.Now().Unix(),
  }
  return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
  t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
    if t.Method != jwt.SigningMethodHS256 {
      return nil, errors.New("unexpected signing method")
    }
    return secret, nil
  })
  if err != nil || !t.Valid {
    return nil, errors.New("invalid token")
  }

  claims, ok := t.Claims.(jwt.MapClaims)
  if !ok {
    return nil, errors.New("invalid claims")
  }
  return claims, nil
}

// Create refresh session (stateful) and return token pair.
func login(userID string) (string, string, error) {
  accessToken, err := issueAccessToken(userID)
  if err != nil {
    return "", "", err
  }

  refreshJTI := uuid.NewString()
  expiresAt := time.Now().Add(7 * 24 * time.Hour)

  sessionMu.Lock()
  sessionTable[refreshJTI] = RefreshSession{
    JTI:       refreshJTI,
    UserID:    userID,
    RevokedAt: nil,
    ExpiresAt: expiresAt,
  }
  sessionMu.Unlock()

  refreshToken, err := issueRefreshToken(userID, refreshJTI)
  if err != nil {
    return "", "", err
  }
  return accessToken, refreshToken, nil
}

// Refresh token flow
// - validate refresh token
// - check sessionTable[jti] active AND belongs to same user
// - revoke old session
// - create new session with new jti
// - return new token pair
func refresh(oldRefreshToken string) (string, string, error) {
  claims, err := parseToken(oldRefreshToken)
  if err != nil {
    return "", "", err
  }

  if claims["typ"] != "refresh" {
    return "", "", errors.New("invalid refresh token")
  }

  userID, _ := claims["sub"].(string)
  jti, _ := claims["jti"].(string)
  if userID == "" || jti == "" {
    return "", "", errors.New("missing refresh claims")
  }

  sessionMu.RLock()
  sess, ok := sessionTable[jti]
  sessionMu.RUnlock()

  if !ok || !isActive(sess) {
    return "", "", errors.New("refresh session revoked/expired")
  }

  // ensure session belongs to this user
  if sess.UserID != userID {
    return "", "", errors.New("refresh session mismatch")
  }

  // revoke old session
  now := time.Now()
  sessionMu.Lock()
  oldSession := sessionTable[jti]
  // re-check in case changed
  if !isActive(oldSession) {
    sessionMu.Unlock()
    return "", "", errors.New("refresh session revoked/expired")
  }
  oldSession.RevokedAt = &now
  sessionTable[jti] = oldSession

  // create new session
  newJTI := uuid.NewString()
  sessionTable[newJTI] = RefreshSession{
    JTI:       newJTI,
    UserID:    userID,
    RevokedAt: nil,
    ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
  }
  sessionMu.Unlock()

  newAccessToken, err := issueAccessToken(userID)
  if err != nil {
    return "", "", err
  }
  newRefreshToken, err := issueRefreshToken(userID, newJTI)
  if err != nil {
    return "", "", err
  }
  return newAccessToken, newRefreshToken, nil
}

// logout revoke a specific refresh session
func logout(refreshToken string) error {
  claims, err := parseToken(refreshToken)
  if err != nil {
    return err
  }
  if claims["typ"] != "refresh" {
    return errors.New("invalid refresh token")
  }
  jti, _ := claims["jti"].(string)
  if jti == "" {
    return errors.New("missing token jti")
  }

  sessionMu.Lock()
  defer sessionMu.Unlock()

  sess, ok := sessionTable[jti]
  if !ok {
    return nil
  }
  now := time.Now()
  sess.RevokedAt = &now
  sessionTable[jti] = sess
  return nil
}

func LoginHandler(c *gin.Context) {
  var req struct {
    UserID string `json:"user_id"`
  }
  if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
    return
  }

  accessToken, refreshToken, err := login(req.UserID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "access_token":  accessToken,
    "refresh_token": refreshToken,
  })
}

func RefreshHandler(c *gin.Context) {
  var req struct {
    RefreshToken string `json:"refresh_token"`
  }
  if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
    return
  }

  newAccessToken, newRefreshToken, err := refresh(req.RefreshToken)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "access_token":  newAccessToken,
    "refresh_token": newRefreshToken,
  })
}

func LogoutHandler(c *gin.Context) {
  var req struct {
    RefreshToken string `json:"refresh_token"`
  }
  if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
    return
  }

  if err := logout(req.RefreshToken); err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    return
  }
  c.Status(http.StatusNoContent)
}

func main() {
  r := gin.Default()

  r.POST("/login", LoginHandler)
  r.POST("/refresh", RefreshHandler)
  r.POST("/logout", LogoutHandler)

  r.Run(":8080")
}
```