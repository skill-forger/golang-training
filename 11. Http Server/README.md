# Module 11: HTTP Server

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#what-is-a-server">What is a Server?</a></li>
    <li><a href="#http-servers">HTTP Servers</a></li>
    <li><a href="#http-server-in-go">HTTP Server in Go</a></li>
    <li><a href="#common-mistakes">Common Mistakes</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:

- Understand the Fundamental Concept of a Server
- Learn the common components of HTTP server (protocol, roles, statuses, methods, url, headers, body)
- Grasp the HTTP Request-Response Cycle
- Utilize the `net/http` in Go to host a simple HTTP server
- Implement basic use cases and handlers using `net/http` package
- Gain Practical Awareness for Backend Development

## Overview

In this module, the focus will be on the foundational principles of HTTP Servers, which are integral to nearly all
internet interactions. This module aims to equip freshers and interns with a clear conceptual understanding of how
these servers operate, gain insight into the core mechanisms of web communication, including the HTTP
protocol, the structure of requests and responses, and the server's role in processing these exchanges, thereby building
a solid foundation for backend development journey.

## What is a Server?

A server functions as a fundamental component responsible for providing information and services to other computing
devices, known as clients. Essentially, a server is a powerful computer program — or the physical hardware executing
that program — designed to remain continuously active, awaiting and responding to client requests

### Concept & Analogy

To facilitate a clearer understanding, consider the analogy of a restaurant

- **The Customer (Client)**: The customer initiates a request by placing an order for food. Similarly, a digital client,
  such as a web browser or a mobile application, sends a request for data or a specific service.
- **The Waiter (Server)**: The waiter's role is to receive the customer's order (the request), relay it to the kitchen (
  where processing occurs), and subsequently deliver the prepared meal (the response) back to the customer. This role
  mirrors that of a server in a digital context.
- **The Kitchen (Processing)**: Within the kitchen, ingredients are assembled, and the meal is prepared according to the
  customer's specifications. This represents the server's internal logic, which processes the client's request by
  executing code, retrieving data, or generating content.
- **The Food (The Response)**: The prepared food delivered to the customer is the direct fulfillment of their request.
  Digitally, this translates to the information or service sent back to the client.

### Core Function

Regardless of its specific application, every server adheres to a tripartite operational cycle:

1. **Listening for Requests**: A server maintains a constant state of readiness, monitoring designated network ports for
   incoming communication from clients. This continuous vigilance ensures that no client request goes unnoticed.
2. **Processing Requests**: Upon receiving a request, the server undertakes an analytical phase to ascertain the
   client's specific requirements. This may involve querying a database, executing complex algorithms, or dynamically
   generating content. The server's programmed logic is then applied to fulfill the request comprehensively.
3. **Sending Responses**: Following the processing phase, the server meticulously packages the resultant information or
   outcome. This could manifest as a complete web page, structured data, a confirmation of a successful operation, or an
   appropriate error message, which is then transmitted back to the initiating client.

### Common Types of Servers

While the current discussion primarily focuses on Web Servers, which are instrumental in facilitating HTTP communication
for websites and APIs, it is important to recognize the broader spectrum of server functionalities:

- **File Servers**: These servers are dedicated to the storage and management of files, enabling secure access and
  sharing among multiple users across a network.
- **Database Servers**: Specialized for the storage, management, and retrieval of vast quantities of structured data,
  database servers are crucial for applications requiring persistent data storage, such as user profiles or product
  inventories.
- **Mail Servers**: These systems are engineered to handle the complexities of email transmission, reception, and
  storage.
- **Game Servers**: In the realm of online gaming, these servers manage the real-time state of multiplayer environments,
  ensuring synchronized gameplay for all participants.

## HTTP Servers

### The Protocol

At the core of web communication lies the **HyperText Transfer Protocol (HTTP)**. This protocol serves as the standard "
language" that clients (such as web browsers or mobile applications) and servers use to exchange data. HTTP defines the
rules for how messages are formatted and transmitted, ensuring that both parties can understand each other's intentions
and responses.

### The Roles

An HTTP server is a specialized type of server program designed to process HTTP requests and deliver HTTP responses.
Its operations involve a systematic sequence of steps to fulfill client demands:

1. **Listening on Specific Ports**: An HTTP server continuously monitors designated network ports for incoming
   connections. For standard unencrypted web traffic, it typically listens on port 80 for HTTP. For secure, encrypted
   communication, it listens on port 443 for HTTPS.
2. **Receiving HTTP Requests**: Upon detection of an incoming connection on a monitored port, the server receives the
   complete HTTP request sent by a client. This request encapsulates all the necessary information for the server to
   understand what the client desires.
3. **Parsing the Request**: After receiving the raw request, the server proceeds to parse it. This involves breaking
   down the request into its constituent parts:
    - _The URL (Uniform Resource Locator)_: Identifying the specific resource the client is requesting (e.g.,
      /products/123).
    - _The HTTP Method_: Indicating the action the client wishes to perform (e.g., GET to retrieve, POST to submit data,
      PUT to update).
    - _Headers_: Metadata providing additional context about the request or the client (e.g., User-Agent, Accept,
      Authorization).
    - _The Body (if present)_: The actual data payload, typically included with POST or PUT requests (e.g., JSON data
      for a new user).
4. **Processing the Request**: With the request parsed, the server's application logic takes over. This phase involves
   executing the necessary code to fulfill the request. This could mean querying a database, performing calculations,
   interacting with other services, or generating dynamic content based on the request's parameters.
5. **Constructing an HTTP Response**: Once the request has been processed, the server meticulously builds an HTTP
   response. This response package includes:
    - _A Status Code_: A three-digit number indicating the outcome of the request (e.g., 200 OK for success, 404 Not
      Found for a missing resource, 500 Internal Server Error for a server-side problem).
    - _Headers_: Metadata about the response itself or the server (e.g., Content-Type specifying the data format in the
      body, Set-Cookie to send cookies to the client).
    - _The Body (if present)_: The actual data or content being sent back to the client (e.g., HTML for a webpage, JSON
      for API data, an image).
6. **Sending the Response**: Finally, the constructed HTTP response is transmitted back to the client that initiated the
   request, completing the communication cycle.

### The Key Components

- **Request Listener**: This component is the server's "ear." Its primary function is to continuously monitor specified
  network ports for incoming client connections and HTTP requests. When a request arrives, the listener accepts the
  connection and passes the raw request data to the appropriate processing unit. It is the initial point of contact for
  all incoming web traffic.
- **Request Handler**: Once a request has been received by the listener, the request handler takes over. This component
  embodies the server's core logic. It is responsible for parsing the incoming request, determining the appropriate
  action based on the URL and HTTP method, executing the necessary application code (which might involve interacting
  with databases or other services), and preparing the data that needs to be sent back.
- **Response Sender**: After the request handler has completed its processing and prepared the necessary data, the
  response sender is tasked with formatting this data into a valid HTTP response. This involves setting the correct HTTP
  status code, adding relevant headers (such as Content-Type or Cache-Control), and encoding the response body as
  required. Finally, it transmits this complete HTTP response back to the waiting client, concluding the transaction.

### URI vs URL

#### Uniform Resource Identifier (URI)

A URI is a sequence of characters that identifies a logical or physical resource. Think of it as a "digital fingerprint"
for anything on the internet or a network. Its primary goal is to locate and retrieve a resource. It tells you where the
resource is and how to get it.

#### Uniform Resource Locator (URL)

A URL is a specific type of URI that not only identifies a resource but also provides a means of locating it and
specifies the mechanism for accessing it. Its primary goal is to identify a resource. It tells you what the resource is.

#### URI and URL Comparison

| Feature          | URI (Uniform Resource Identifier)                                            | URL (Uniform Resource Locator)                                |
|------------------|------------------------------------------------------------------------------|---------------------------------------------------------------|
| **Purpose**      | Identifies a resource (What is it?)                                          | Locates and accesses a resource (Where is it? How to get it?) |
| **Scope**        | Broader, generic identifier                                                  | Specific type of URI                                          |
| **Includes**     | Can be a name (URN), a location (URL), or both                               | Always includes location and access method                    |
| **Analogy**      | Book's ISBN (identifies it)                                                  | Book's shelf location (how to find it)                        |
| **Examples**     | `urn:isbn:1234567890`, `mailto:user@example.com`, `https://example.com/page` | `https://example.com/page`, `ftp://server/file.txt`           |
| **Relationship** | Superset (contains URLs)                                                     | Subset of URI                                                 |

### The parameters

#### Path Parameters

Path parameters are used to identify a specific resource or a specific part of a resource's hierarchy.
They are an essential part of the URL itself

- **Syntax**: They're embedded directly in the URL path, often marked with curly braces in API documentation
  (e.g., `/users/{id}`). In the actual request, the value replaces the placeholder
- **Purpose**: They are typically required and are used for resource identification.
- **Example**: If you want to get the details of a specific user with the ID 123,
  your request would look like this: `GET /users/123`

#### Query Parameters

Query parameters are used to **filter**, **sort**, and **paginate** resources.
They provide optional instructions for modifying a request without changing the resource's core identity

- **Syntax**: They are appended to the URL after a question mark (`?`).
  Each parameter is a key-value pair, separated by an ampersand (`&`) if there's more than one.
- **Purpose**: They are typically optional and are used for filtering, sorting, searching, or pagination.
- **Example**: If you want to get a list of products that are in the "electronics" category and sort them by price,
  your request would look like this: `GET /products?category=electronics&sort=price`

#### Comparison

| Feature         | Path Parameters                                | Query Parameters                                                             |
|-----------------|------------------------------------------------|------------------------------------------------------------------------------|
| **Location**    | Part of the URL path                           | Appended after the `?` in the URL                                            |
| **Purpose**     | Identify a specific resource or its hierarchy  | Filter, sort, paginate, or provide optional data                             |
| **Requirement** | Typically required for the request to be valid | Typically optional                                                           |
| **Analogy**     | A street address: "Go to this exact place."    | Search filters: "Show me this place, but only the ones with these features." |

### The Methods

The HyperText Transfer Protocol (HTTP) defines a set of standard request methods, often referred to as HTTP verbs. These
methods indicate the desired action to be performed on the resource identified by the Request-URI. Each method carries a
specific semantic meaning, guiding how clients and servers interact with web resources. Understanding these methods is
fundamental for designing robust and predictable web APIs. Common HTTP methods for client and server communication:

- **GET**: Retrieve data (e.g., getting a web page). Idempotent & Safe.
- **POST**: Submit data (e.g., submitting a form). Not idempotent.
- **PUT**: Update/replace existing data. Idempotent.
- **DELETE**: Remove data. Idempotent.
- **PATCH**: Partially update data. Not idempotent.

For the full list of HTTP methods, please refer
to [HTTP Methods (Verbs)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods)

### The Status

HTTP Status Codes constitute a critical component of every HTTP response. These three-digit integers, accompanied by a
brief, standardized phrase, serve as a concise and standardized means for the server to inform the client about the
outcome of its request. They convey whether a request has been successfully completed, requires further action, or
encountered an error. Understanding these codes is essential for debugging web applications, interpreting API responses,
and building robust client-server interactions. Common HTTP statuses for client and server communication:

- **2xx (Success)**: 200 OK, 201 Created
- **3xx (Redirection)**: 301 Moved Permanently, 302 Found
- **4xx (Client Error)**: 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found
- **5xx (Server Error)**: 500 Internal Server Error, 503 Service Unavailable

For the full list of HTTP statuses, please refer
to [HTTP Statuses](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status)

### The Headers

HTTP headers constitute a vital component of every HTTP request and response message. Positioned within the header
section of these messages, they function as fields that transmit crucial metadata about the message itself, the sender,
the recipient, or the specific resource being transferred. These headers are indispensable for enabling clients and
servers to understand each other's context, negotiate capabilities, and effectively process web communication.

#### HTTP Header Structure

Each HTTP header adheres to a straightforward **key-value pair** format. The header name (or key) and its corresponding
value are separated by a colon (:). Each distinct header occupies its own line within the message's header section.
Example: `Content-Type: application/json`

- _Content-Type_ represents the **header name** (or key), which identifies the type of information being conveyed
- _application/json_ is the **header value**, providing the specific detail associated with the header name, in this
  case, indicating that the message body contains data formatted as JSON

#### HTTP Headers Types

- **General Headers**: These headers are applicable to both request and response messages but are independent of the
  data being transmitted in the message body. Examples include Date (indicating the message's origination time) or Via (
  denoting intermediate proxies).
- **Request Headers**: These headers are sent from the client to the server and provide detailed information about the
  client, the request itself, or the client's preferred characteristics for the response. Key request headers include:
    - _User-Agent_: Identifies the client software originating the request (e.g., browser name and version).
    - _Accept_: Specifies the media types (e.g., text/html, application/json) that the client is capable of processing
      and prefers to receive.
    - _Accept-Language_: Indicates the human languages preferred by the client.
    - _Authorization_: Carries credentials (e.g., API tokens, session IDs) for authenticating the client with the
      server.
    - _Cookie_: Transmits HTTP cookies that the server previously sent to the client, facilitating session management.
    - _Host_: Specifies the domain name of the server being requested, crucial for servers hosting multiple websites (
      virtual hosting).
    - _Content-Type_: For requests with a body (e.g., POST, PUT), this header specifies the media type of the data
      contained within the request body.
    - _Content-Length_: Indicates the size of the request body in bytes.
- **Response Headers**: These headers are sent from the server to the client and provide information about the server,
  the response itself, or the resource being sent back. Important response headers include:
    - _Content-Type_: Informs the client about the media type of the data within the response body, enabling the
      client (e.g., browser) to correctly render or process the content.
    - _Content-Length_: Specifies the size of the response body in bytes.
    - _Cache-Control_: Contains directives for caching mechanisms, guiding how the response should be cached by clients
      or intermediaries.
    - _Set-Cookie_: Instructs the client to store a specific cookie, which will then be sent back by the client in
      subsequent requests for the same domain.
    - _Server_: Provides information about the web server software.
    - _Location_: Used in conjunction with redirection status codes (3xx) to inform the client of the new URI to which
      it should redirect.

### The Body

Within an HTTP message, distinct from the headers that carry metadata, resides the message body. This component
constitutes the actual data payload being transmitted between the client and the server. Its presence and content are
entirely dependent on the specific HTTP method used in the request and the nature of the information being exchanged in
the response. The message body is crucial for transferring the core information that drives web applications, whether it
involves submitting user data, retrieving documents, or exchanging API payloads.

#### In an HTTP Request

When a client wishes to send data to the server (e.g., to create a new resource or update an existing one), this data is
typically encapsulated within the request body. Common scenarios include

- Submitting form data from a web page (often as `application/x-www-form-urlencoded` or `multipart/form-data`).
- Sending JSON or XML data to a RESTful API endpoint for creating or updating records (typically `application/json`
  or `application/xml`).
- Uploading files, such as images or documents (`multipart/form-data`).

#### In an HTTP Response

When a server fulfills a client's request for a resource, the actual content of that resource is delivered in the
response body. Examples include:

- The HTML content of a webpage requested by a browser.
- An image file (e.g., `image/jpeg`, `image/png`).
- JSON or XML data representing a data structure from an API.
- Plain text, CSS stylesheets, or JavaScript files.

#### The Body is Optional

It is important to note that an HTTP message does not always include a body.

- **Requests**
    - _GET_ requests, by definition, are for retrieving data and therefore typically do not have a request body. While
      some clients might technically send a body with GET, this practice is non-standard and generally ignored by
      servers.
    - _HEAD_ requests, used to retrieve only the headers of a resource, explicitly do not have a body.
    - _DELETE_ requests, while they trigger server-side changes, usually do not require a request body as the resource
      to be deleted is identified by the URL.

- **Responses**
    - Responses to HEAD requests will never have a body.
    - Responses to GET requests for resources that do not exist (e.g., 404 Not Found) or some redirection responses (
      3xx) may not include a body, or might include a small descriptive body.
    - Certain informational responses (1xx) and specific success codes like 204 No Content explicitly state that no
      content is returned in the body.

### The Middleware

#### Definition

Middleware is a function that sits between a request and the final handler. It can execute code before or after the
handler. This concept allows you to add cross-cutting concerns to your application in a modular and reusable way.

#### Common use cases

- **Logging**: To log details of incoming requests. Gin comes with a built-in gin.Logger() middleware.
- **Authentication/Authorization**: To protect routes and ensure users are logged in or have the correct permissions.
- **Recovery**: The gin.Recovery() middleware helps prevent the server from crashing on a panic.

Example of a logging middleware
This is a classic example of middleware.
It logs information about each incoming request before it's passed to the next handler.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// The middleware function
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
		// Log the request after the handler has completed
		log.Printf("Request processed in %s | Method: %s | URL: %s", time.Since(start), r.Method, r.URL.Path)
	})
}

// Our main application handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Create an http.Handler from our function
	helloHandlerFunc := http.HandlerFunc(helloHandler)

	// Wrap our handler with the logging middleware
	loggedHandler := loggingMiddleware(helloHandlerFunc)

	// Register the wrapped handler
	http.Handle("/", loggedHandler)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

#### Multiple Middlewares Chaining

Multiple middleware functions can be chained together to create a pipeline. Each middleware function wraps the one that
comes after it, creating a nested structure.

For example, to add both a logging and an authentication middleware:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Authentication middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: Check for a valid header
		token := r.Header.Get("X-Auth-Token")
		if token != "valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return // Stop the chain
		}
		next.ServeHTTP(w, r) // Pass to the next handler
	})
}

// The middleware function
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
		// Log the request after the handler has completed
		log.Printf("Request processed in %s | Method: %s | URL: %s", time.Since(start), r.Method, r.URL.Path)
	})
}

// Our main application handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Final handler
	helloHandlerFunc := http.HandlerFunc(helloHandler)

	// Chain the middleware
	chain := loggingMiddleware(authMiddleware(helloHandlerFunc))

	http.Handle("/", chain)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

```

## HTTP Server in Go

Go provides powerful, yet straightforward tools for building web applications through its standard library.
The `net/http` package offers everything needed to create robust HTTP servers without requiring external dependencies,
embodying Go's philosophy of simplicity and efficiency.

### Serving Basic HTTP Server

The standard library's `net/http` package makes it incredibly easy to create web servers with just a few lines of code:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// Basic HTTP server example
func main() {
	// Define a handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s", r.URL.Path)
	})

	// Start the server on port 8080
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

```

### Handling HTTP Request

Go's HTTP server revolves around the `http.Handler` interface, which defines how to process incoming requests:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// CustomHandler implementation
type CustomHandler struct{}

func (c CustomHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintf(writer, "Custom handler serving: %s", request.URL.Path)
}

func main() {
	handler := new(CustomHandler)

	// Register our custom handler for all routes
	http.Handle("/", handler)

	// Start the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

```

#### Working with Form Data

Form data from request body can be extracted directly using `http.Request`

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	// Access form values
	name := r.FormValue("name")

	// Write response
	fmt.Fprintf(w, "Form submission successful. Name = %s", name)
}

func main() {
	http.HandleFunc("/form", formHandler)
	// Start the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

```

### Working with JSON

Working with JSON is common in web applications and APIs:

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create sample user
	var user = User{
		ID:       1,
		Username: "gopher",
		Email:    "gopher@example.com",
	} // Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send JSON response
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func apiCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the user (in a real app, save to database, etc.)

	// Return created status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created successfully", user.Username)
}

```

### Routing with ServeMux

The `http.ServeMux` provides a flexible routing mechanism for your web applications:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Create a new router
	mux := http.NewServeMux()

	// Register route handlers
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/api/", apiHandler)

	// Start server with custom router
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About Us")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	fmt.Fprintf(w, "API request for: %s", path)
}

```

### Serving Static Files

Go makes it easy to serve static assets like images, CSS, and JavaScript:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create file server handler
	fs := http.FileServer(http.Dir("./static"))

	// Register the handler
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register other routes
	http.HandleFunc("/", indexHandler)

	// Start server with custom router
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serve an HTML page that references static files
	html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Go Web Server</title>
            <link rel="stylesheet" href="/static/style.css">
        </head>
        <body>
            <h1>Welcome to Go Web Development</h1>
            <img src="/static/gopher.png" alt="Go Gopher">
            <script src="/static/script.js"></script>
        </body>
    </html>
    `
	fmt.Fprint(w, html)
}

```

### Applying Middleware

Middleware functions in Go allow you to wrap handlers with common functionality:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Post-processing logic
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(startTime))
	})
}

func main() {
	// Create a handler
	handler := http.HandlerFunc(homeHandler)

	// Wrap it with middleware
	wrappedHandler := loggingMiddleware(handler)

	// Register the wrapped handler
	http.Handle("/", wrappedHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

```

### Securing HTTP Server

Securing your web application with TLS/SSL:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Configure TLS server
	server := &http.Server{
		Addr:         ":8443",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Register handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Secure HTTPS server")
	})

	// Start HTTPS server
	log.Println("Starting HTTPS server on :8443")
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

### Setting Server Timeout

Properly configuring your HTTP server for production:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
}

func main() {
	// Create a custom server with configurations
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Register handlers
	http.HandleFunc("/", homeHandler)

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Create shutdown context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the server
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// Start server
	log.Printf("Starting HTTP server on %s", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
	log.Println("Server gracefully stopped")
}

```

## Common Mistakes

1. **Ignoring HTTP Methods (Verbs) Semantics**: Using `GET` for operations that change server state (e.g.,
   /deleteUser?id=123) or using `POST` when `PUT` or `PATCH` would be more semantically appropriate for updates.
   This often leads to non-idempotent operations being executed repeatedly if a client retries a `GET` request,
   or inefficient API design
2. **Lack of Input Validation and Sanitization**:  Trusting all data received in request bodies or URL parameters
   directly from the client. Junior developers might process user input without checking for expected formats,
   lengths, or malicious content
3. **Inadequate Error Handling and Status Codes**: Consistently returning `200 OK` for all responses, even when an error
   occurred, or using generic `500 Internal Server Error` for all client-side issues. Conversely, sometimes
   returning `404 Not Found` when a resource exists but the user is unauthorized
4. **Misunderstanding Statelessness and Session Management**: Expecting the server to automatically remember previous
   client interactions without explicit mechanisms. This often leads to issues with maintaining user sessions or
   carrying context between requests.
5. **Ignoring HTTP Headers**: Not setting appropriate `Content-Type` headers in responses, causing clients (especially
   browsers) to misinterpret data. Also, neglecting `Cache-Control` headers, leading to inefficient caching or serving
   stale content.
6. **Hard coding Configurations and Sensitive Information**: Embedding database credentials, API keys, or server port
   numbers directly within the code rather than using environment variables or configuration files.
7. **Poor Logging and Monitoring**: Not implementing sufficient logging for server-side operations, or logging too much
   irrelevant information. Also, not understanding how to use monitoring tools to observe server health and performance.
8. **Security Overlooks (e.g., CORS, HTTPS)**: Not properly configuring Cross-Origin Resource Sharing (CORS) policies,
   leading to client-side fetching issues. Also, deploying servers without HTTPS, leaving data vulnerable during
   transit.
9. **Over-engineering or Under-engineering**:  Trying to implement complex architectural patterns (e.g., microservices)
   for simple applications, or conversely, writing monolithic code for complex systems without proper modularization

## Best Practices

1. **Adhere to HTTP Method Semantics (RESTfulness)**: Always use the correct HTTP
   method (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`) that semantically aligns with the desired operation on the resource.
2. **Implement Robust Input Validation and Sanitization**: Never trust client-side input. All data received from
   requests ( URL parameters, query strings, headers, and especially request bodies)
   must be rigorously validated against expected formats, types, and constraints.
   Additionally, sanitize inputs to neutralize potentially malicious content
   (e.g., stripping HTML tags, escaping special characters).
3. **Utilize Appropriate HTTP Status Codes Consistently**:  Return precise HTTP status codes to clearly communicate the
   outcome of a request. Use `2xx` for success, `4xx` for client-side errors, and `5xx` for server-side errors. Provide
   meaningful error messages in the response body for 4xx and 5xx codes.
4. **Manage State Explicitly (Session Management)**: Since HTTP is inherently stateless, implement explicit mechanisms
   for session management when persistent state is required.
   Common approaches include using cookies to store session IDs or JSON Web Tokens (JWTs) for stateless authentication.
5. **Leverage HTTP Headers Effectively**: Set appropriate `Content-Type` headers in responses to ensure clients
   correctly interpret the data. Utilize `Cache-Control` headers to optimize caching and reduce server load. Implement
   security-related headers (e.g., `Strict-Transport-Security`, `X-Content-Type-Options`).
6. **Prioritize Asynchronous (Non-Blocking) I/O**: In environments that support it (e.g., Node.js, asynchronous
   Python/Java), ensure that I/O operations (database calls, file system access, external API requests) are
   non-blocking. This allows the server to process other requests concurrently without waiting for slow operations to
   complete.
7. **Externalize Configuration and Sensitive Data**: Never hardcode sensitive information (e.g., database credentials,
   API keys, encryption secrets) or environment-specific settings (e.g., port numbers, external service URLs)
   directly into the codebase. Instead, use environment variables, configuration files, or secret management services.
8. **Implement Comprehensive Logging and Monitoring**: Establish clear logging policies. Log meaningful events (e.g.,
   request received, errors, critical business operations) with sufficient context (timestamps, request IDs, user IDs).
   Integrate with monitoring tools to track server health, performance metrics (latency, throughput, error rates), and
   resource utilization.
9. **Enforce Security Best Practices (HTTPS, CORS, Authentication/Authorization)**: Always use HTTPS for all production
   traffic to encrypt data in transit. Properly configure Cross-Origin Resource Sharing (CORS) policies to control which
   origins can access your API. Implement robust authentication (verifying user identity) and authorization (checking
   user permissions)

## Practice Exercises

### Exercise 1: Simple REST-ful API

Build a simple REST-ful API to manage a collection of books:

### Exercise 2: File Server with Custom Handler

Create a file server with a custom middleware for logging:

### Exercise 3: HTTP Client and Concurrent Requests

Build an HTTP client that makes concurrent requests to different APIs:

