## Practical Exercises

### Exercise 1: JWT Authentication in Go

Create a simple API with JWT authentication:

```go
// jwt_auth_server.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing tokens
var jwtSecret = []byte("your-secret-key-should-be-longer-and-secure")

// User represents application user
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Password is not included in JSON output
	Role     string `json:"role"`
}

// LoginRequest contains login form data
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// JWTClaims contains JWT claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Database simulation - in a real application, these would be in a database
var users = []User{
	{ID: 1, Username: "admin", Email: "admin@example.com", Password: "admin123", Role: "admin"},
	{ID: 2, Username: "user", Email: "user@example.com", Password: "user123", Role: "user"},
}

// findUserByCredentials returns a user if valid credentials are provided
func findUserByCredentials(username, password string) (*User, bool) {
	for _, user := range users {
		if user.Username == username && user.Password == password {
			return &user, true
		}
	}
	return nil, false
}

// generateToken creates a new JWT token for a user
func generateToken(user *User) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-auth-server",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the encoded token string
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// validateToken validates a JWT token
func validateToken(tokenString string) (*JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and return claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// JWTAuthMiddleware authenticates requests using JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Remove 'Bearer ' prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validate the token
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleAuthMiddleware restricts access based on user role
func RoleAuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by JWTAuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Check if user has required role
		userRole := role.(string)
		allowed := false
		for _, r := range roles {
			if r == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Public endpoints
	public := r.Group("/api")
	{
		public.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to the JWT Auth API",
			})
		})

		public.POST("/login", func(c *gin.Context) {
			var req LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Validate credentials
			user, found := findUserByCredentials(req.Username, req.Password)
			if !found {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// Generate token
			token, err := generateToken(user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			// Return token
			c.JSON(http.StatusOK, gin.H{
				"token":   token,
				"user_id": user.ID,
				"role":    user.Role,
				"expires": time.Now().Add(24 * time.Hour).Unix(),
			})
		})
	}

	// Protected endpoints (require authentication)
	protected := r.Group("/api/protected")
	protected.Use(JWTAuthMiddleware())
	{
		// Endpoint accessible to both users and admins
		protected.GET("/profile", func(c *gin.Context) {
			username, _ := c.Get("username")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"message":  "You have access to your profile",
				"username": username,
				"role":     role,
			})
		})

		// Endpoint accessible only to admins
		protected.GET("/admin", RoleAuthMiddleware("admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "You have admin access",
			})
		})
	}

	log.Println("Starting server on :8080...")
	r.Run(":8080")
}
```

### Exercise 2: OAuth2 Integration

Build a simple OAuth2 client that authenticates with Google:

```go
// oauth2_client.go
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo represents user information returned from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// OAuth2Config holds OAuth2 configuration
var OAuth2Config = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),     // Set these environment variables
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), // or replace with actual values
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// generateStateOauthCookie creates a state token and stores it in a cookie
func generateStateOauthCookie(c *gin.Context) string {
	// Generate a random state
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in session
	session := sessions.Default(c)
	session.Set("oauth_state", state)
	session.Save()

	return state
}

// getUserInfoFromGoogle fetches the user's info from Google API
func getUserInfoFromGoogle(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	// Create HTTP client using the provided token
	client := OAuth2Config.Client(ctx, token)

	// Make a request to Google API
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func main() {
	// Check if Google OAuth credentials are set
	if OAuth2Config.ClientID == "" || OAuth2Config.ClientSecret == "" {
		log.Println("Warning: Google OAuth credentials not set. Using dummy values for demo.")
		OAuth2Config.ClientID = "dummy_client_id"
		OAuth2Config.ClientSecret = "dummy_client_secret"
	}

	r := gin.Default()

	// Set up session middleware
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("oauth-session", store))

	// Serve static files
	r.Static("/static", "./static")

	// Home page
	r.GET("/", func(c *gin.Context) {
		// Check if user is logged in
		session := sessions.Default(c)
		userInfo := session.Get("user_info")

		if userInfo != nil {
			// User is logged in, serve dashboard
			c.HTML(http.StatusOK, "dashboard.html", gin.H{
				"user": userInfo,
			})
		} else {
			// User is not logged in, serve login page
			html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>OAuth2 Demo</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
					h1 { color: #333; }
					.login-box { 
						margin: 20px 0; 
						padding: 20px; 
						border: 1px solid #ddd; 
						border-radius: 5px;
						text-align: center;
					}
					.btn {
						display: inline-block;
						padding: 10px 20px;
						background-color: #4285F4;
						color: white;
						text-decoration: none;
						border-radius: 5px;
						font-weight: bold;
					}
				</style>
			</head>
			<body>
				<h1>Welcome to OAuth2 Demo</h1>
				
				<div class="login-box">
					<h2>Please Sign In</h2>
					<a href="/auth/google/login" class="btn">Sign in with Google</a>
				</div>
			</body>
			</html>
			`
			c.Header("Content-Type", "text/html")
			c.String(http.StatusOK, html)
		}
	})

	// Login endpoint - redirects to Google
	r.GET("/auth/google/login", func(c *gin.Context) {
		// Generate and store state
		state := generateStateOauthCookie(c)

		// Redirect to Google's OAuth 2.0 server
		url := OAuth2Config.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	// Callback endpoint - handles the response from Google
	r.GET("/auth/google/callback", func(c *gin.Context) {
		// Get state from session
		session := sessions.Default(c)
		sessionState := session.Get("oauth_state")

		// Compare state parameter
		if c.Query("state") != sessionState {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
			return
		}

		// Exchange code for token
		token, err := OAuth2Config.Exchange(context.Background(), c.Query("code"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token: " + err.Error()})
			return
		}

		// Get user info
		userInfo, err := getUserInfoFromGoogle(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info: " + err.Error()})
			return
		}

		// Store user info in session
		session.Set("user_info", userInfo)
		session.Set("access_token", token.AccessToken)
		session.Set("logged_in_at", time.Now().Format(time.RFC3339))
		session.Save()

		// Redirect to home page
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	// Dashboard page - only accessible when logged in
	r.GET("/dashboard", func(c *gin.Context) {
		// Check if user is logged in
		session := sessions.Default(c)
		userInfo := session.Get("user_info")

		if userInfo == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Serve dashboard HTML with user info
		dashboardHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Dashboard - OAuth2 Demo</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
				h1 { color: #333; }
				.profile {
					display: flex;
					align-items: center;
					margin-bottom: 20px;
				}
				.profile img {
					width: 100px;
					height: 100px;
					border-radius: 50%%;
					margin-right: 20px;
				}
				.card { 
					margin: 20px 0; 
					padding: 20px; 
					border: 1px solid #ddd; 
					border-radius: 5px;
				}
				.btn {
					display: inline-block;
					padding: 10px 20px;
					background-color: #f44336;
					color: white;
					text-decoration: none;
					border-radius: 5px;
					font-weight: bold;
				}
			</style>
		</head>
		<body>
			<h1>Dashboard</h1>
			
			<div class="profile">
				<img src="%s" alt="Profile Picture">
				<div>
					<h2>Welcome, %s!</h2>
					<p>Email: %s</p>
					<p>Logged in at: %s</p>
				</div>
			</div>
			
			<div class="card">
				<h3>Protected Content</h3>
				<p>This page is only accessible to authenticated users.</p>
			</div>
			
			<a href="/logout" class="btn">Logout</a>
		</body>
		</html>
		`, 
		userInfo.(GoogleUserInfo).Picture,
		userInfo.(GoogleUserInfo).Name,
		userInfo.(GoogleUserInfo).Email,
		session.Get("logged_in_at").(string),
		)

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, dashboardHTML)
	})

	// Logout endpoint
	r.GET("/logout", func(c *gin.Context) {
		// Clear session
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		// Redirect to home page
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	log.Println("Starting server on :8080...")
	r.Run(":8080")
}
```

### Exercise 3: RBAC (Role-Based Access Control) System

Implement a role-based access control system:

```go
// rbac_system.go
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing tokens
var jwtSecret = []byte("your-secret-key-for-rbac-system")

// User represents application user
type User struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"-"` // Not included in JSON
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// Resource represents a protected resource
type Resource struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// Role defines a set of permissions
type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// Permission represents an action that can be performed
type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// JWTClaims contains JWT claims with roles
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// RBAC system configuration
type RBACSystem struct {
	Users       []User
	Roles       map[string]Role
	Permissions map[string]Permission
	Resources   map[string]Resource
}

// NewRBACSystem initializes the RBAC system with sample data
func NewRBACSystem() *RBACSystem {
	// Define permissions
	permissions := map[string]Permission{
		"read":       {Name: "read", Description: "Read access to a resource"},
		"write":      {Name: "write", Description: "Write access to a resource"},
		"delete":     {Name: "delete", Description: "Delete access to a resource"},
		"admin":      {Name: "admin", Description: "Administrative access"},
		"create_user": {Name: "create_user", Description: "Can create new users"},
		"manage_users": {Name: "manage_users", Description: "Can manage existing users"},
	}

	// Define roles with permissions
	roles := map[string]Role{
		"user": {
			Name:        "user",
			Description: "Regular user with limited access",
			Permissions: []string{"read"},
		},
		"editor": {
			Name:        "editor",
			Description: "Can edit content",
			Permissions: []string{"read", "write"},
		},
		"manager": {
			Name:        "manager",
			Description: "Can manage content and users",
			Permissions: []string{"read", "write", "delete", "create_user"},
		},
		"admin": {
			Name:        "admin",
			Description: "Full administrative access",
			Permissions: []string{"read", "write", "delete", "admin", "create_user", "manage_users"},
		},
	}

	// Define protected resources
	resources := map[string]Resource{
		"articles": {
			ID:          "articles",
			Name:        "Articles",
			Description: "Blog articles and posts",
			Permissions: []string{"read", "write", "delete"},
		},
		"users": {
			ID:          "users",
			Name:        "Users",
			Description: "User management",
			Permissions: []string{"read", "create_user", "manage_users"},
		},
		"settings": {
			ID:          "settings",
			Name:        "System Settings",
			Description: "Application configuration",
			Permissions: []string{"read", "write", "admin"},
		},
	}

	// Define users
	users := []User{
		{
			ID:       1,
			Username: "regular_user",
			Password: "password123",
			Email:    "user@example.com",
			Roles:    []string{"user"},
		},
		{
			ID:       2,
			Username: "content_editor",
			Password: "password123",
			Email:    "editor@example.com",
			Roles:    []string{"user", "editor"},
		},
		{
			ID:       3,
			Username: "manager",
			Password: "password123",
			Email:    "manager@example.com",
			Roles:    []string{"user", "editor", "manager"},
		},
		{
			ID:       4,
			Username: "admin",
			Password: "admin123",
			Email:    "admin@example.com",
			Roles:    []string{"admin"},
		},
	}

	return &RBACSystem{
		Users:       users,
		Roles:       roles,
		Permissions: permissions,
		Resources:   resources,
	}
}

// FindUserByCredentials returns a user if valid credentials are provided
func (rbac *RBACSystem) FindUserByCredentials(username, password string) (*User, bool) {
	for _, user := range rbac.Users {
		if user.Username == username && user.Password == password {
			return &user, true
		}
	}
	return nil, false
}

// GetUserPermissions returns all permissions for a user based on their roles
func (rbac *RBACSystem) GetUserPermissions(user *User) []string {
	permissionSet := make(map[string]bool)

	// Add permissions from each role
	for _, roleName := range user.Roles {
		if role, exists := rbac.Roles[roleName]; exists {
			for _, permission := range role.Permissions {
				permissionSet[permission] = true
			}
		}
	}

	// Convert map to slice
	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	return permissions
}

// CheckAccess checks if a user has the required permission for a resource
func (rbac *RBACSystem) CheckAccess(user *User, resourceID, permission string) bool {
	// Get the resource
	resource, exists := rbac.Resources[resourceID]
	if !exists {
		return false
	}

	// Check if the resource supports this permission
	permissionAllowed := false
	for _, p := range resource.Permissions {
		if p == permission {
			permissionAllowed = true
			break
		}
	}

	if !permissionAllowed {
		return false
	}

	// Check if the user has this permission through any of their roles
	userPermissions := rbac.GetUserPermissions(user)
	for _, p := range userPermissions {
		if p == permission || p == "admin" { // Admin permission grants access to everything
			return true
		}
	}

	return false
}

// GenerateToken creates a JWT token for a user
func (rbac *RBACSystem) GenerateToken(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rbac-system",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (rbac *RBACSystem) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AuthMiddleware authenticates requests using JWT
func (rbac *RBACSystem) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Remove 'Bearer ' prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validate token
		claims, err := rbac.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Get user by ID (in a real app, this would be a database lookup)
		var user *User
		for i, u := range rbac.Users {
			if u.ID == claims.UserID {
				user = &rbac.Users[i]
				break
			}
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Store user and roles in context
		c.Set("user", user)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// RequirePermission creates middleware to check if user has permission for a resource
func (rbac *RBACSystem) RequirePermission(resourceID, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userValue, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		user := userValue.(*User)

		// Check if user has required permission
		if !rbac.CheckAccess(user, resourceID, permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Permission denied: requires '%s' permission for '%s'", 
					permission, resourceID),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	// Initialize RBAC system
	rbac := NewRBACSystem()

	// Create Gin router
	r := gin.Default()

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the RBAC System",
			"version": "1.0",
		})
	})

	// Login endpoint
	r.POST("/login", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user
		user, found := rbac.FindUserByCredentials(creds.Username, creds.Password)
		if !found {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate token
		token, err := rbac.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Return user info and token
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"roles":    user.Roles,
			},
			"token":      token,
			"expires_at": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		})
	})

	// Protected routes
	api := r.Group("/api")
	api.Use(rbac.AuthMiddleware())
	{
		// User info endpoint - any authenticated user can access
		api.GET("/me", func(c *gin.Context) {
			user := c.MustGet("user").(*User)
			permissions := rbac.GetUserPermissions(user)

			c.JSON(http.StatusOK, gin.H{
				"user": gin.H{
					"id":          user.ID,
					"username":    user.Username,
					"email":       user.Email,
					"roles":       user.Roles,
					"permissions": permissions,
				},
			})
		})

		// Articles endpoints
		articles := api.Group("/articles")
		{
			// List articles - requires 'read' permission on 'articles'
			articles.GET("", rbac.RequirePermission("articles", "read"), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"articles": []gin.H{
						{"id": 1, "title": "Introduction to RBAC", "author": "John Doe"},
						{"id": 2, "title": "Advanced Authentication", "author": "Jane Smith"},
						{"id": 3, "title": "Security Best Practices", "author": "Bob Johnson"},
					},
				})
			})

			// Create article - requires 'write' permission on 'articles'
			articles.POST("", rbac.RequirePermission("articles", "write"), func(c *gin.Context) {
				var article struct {
					Title   string `json:"title" binding:"required"`
					Content string `json:"content" binding:"required"`
				}

				if err := c.ShouldBindJSON(&article); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"id":      4, // In a real app, this would be generated
					"title":   article.Title,
					"content": article.Content,
					"author":  c.MustGet("user").(*User).Username,
				})
			})

			// Delete article - requires 'delete' permission on 'articles'
			articles.DELETE("/:id", rbac.RequirePermission("articles", "delete"), func(c *gin.Context) {
				id := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Article %s deleted successfully", id),
				})
			})
		}

		// User management endpoints
		users := api.Group("/users")
		{
			// List users - requires 'read' permission on 'users'
			users.GET("", rbac.RequirePermission("users", "read"), func(c *gin.Context) {
				// In a real app, you'd probably want to exclude passwords
				userList := make([]gin.H, len(rbac.Users))
				for i, user := range rbac.Users {
					userList[i] = gin.H{
						"id":       user.ID,
						"username": user.Username,
						"email":    user.Email,
						"roles":    user.Roles,
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"users": userList,
				})
			})

			// Create user - requires 'create_user' permission on 'users'
			users.POST("", rbac.RequirePermission("users", "create_user"), func(c *gin.Context) {
				var newUser struct {
					Username string   `json:"username" binding:"required"`
					Email    string   `json:"email" binding:"required"`
					Password string   `json:"password" binding:"required"`
					Roles    []string `json:"roles" binding:"required"`
				}

				if err := c.ShouldBindJSON(&newUser); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// In a real app, this would add to the database
				c.JSON(http.StatusCreated, gin.H{
					"message": "User created successfully",
					"user": gin.H{
						"id":       len(rbac.Users) + 1,
						"username": newUser.Username,
						"email":    newUser.Email,
						"roles":    newUser.Roles,
					},
				})
			})

			// Manage user - requires 'manage_users' permission on 'users'
			users.PUT("/:id", rbac.RequirePermission("users", "manage_users"), func(c *gin.Context) {
				id := c.Param("id")
				var updateUser struct {
					Roles []string `json:"roles" binding:"required"`
				}

				if err := c.ShouldBindJSON(&updateUser); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("User %s updated successfully", id),
					"roles":   updateUser.Roles,
				})
			})
		}

		// Settings endpoint - requires 'admin' permission on 'settings'
		api.GET("/settings", rbac.RequirePermission("settings", "admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"settings": gin.H{
					"app_name":     "RBAC Demo",
					"version":      "1.0.0",
					"debug_mode":   false,
					"max_users":    1000,
					"allowed_ips":  []string{"127.0.0.1", "::1"},
					"config_path":  "/etc/rbac-demo/config.json",
					"log_level":    "info",
					"backup_cron":  "0 0 * * *", // Daily at midnight
					"admin_email":  "admin@example.com",
					"system_theme": "light",
				},
			})
		})
	}

	// Start server
	log.Println("Starting RBAC server on :8080...")
	r.Run(":8080")
}
