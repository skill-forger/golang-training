# Module 09: Packages and Modules

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#packages">Packages</a></li>
    <li><a href="#modules">Modules</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#common-mistakes">Common Mistakes</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:

- Understand the concept and purpose of packages and modules
- Learn how to use different packages using import syntax
- Understand the exported and unexported visibility
- Learn how to use common commands in Go module
- Know common best practices when using packages and modules
- Avoid common mistakes when working with interfaces

## Overview

Go's package and module system is a powerful mechanism for organizing, sharing, and managing code.
This module will provide a deep dive into understanding how packages and modules work in Go,
exploring their structure, usage, and best practices.

## Packages

### Introduction to Go Packages

A package in Go is a collection of source files in the same directory that are compiled together.
Packages provide a way to:

- Organize and modularize code
- Control visibility of identifiers
- Create reusable code components
- Isolate and reuse code

Every Go source file must start with a package declaration:

```go
package main  // For executable programs
package utils // For library packages
```

Code in a package can access and use all types, constants, variables and functions within that package,
even if they are declared in a different .go file.
Here is an example of the main package with a simple function to generate random number

- File: **random.go**
    ```go
    // random.go
    package main
    
    import (
        "math/rand"
    )
    
    func randomNumber() int {
        return rand.Intn(100)
    }
    ```
- File: **main.go**
    ```go
    // main.go
    package main
    
    import (
    "fmt"
    )
    
    func main() {
        fmt.Printf("Your lucky number is %d!\n", randomNumber())
    }
    ```

### The main package

In Go, `main` is actually a special package name
which indicates that the package contains the code for an executable application.
That is, it indicates that the package contains code that can be built into a binary and run.

Any package with the name main must also contain a `main()` function somewhere in the package
which acts as the entry point for the program.
If it doesn't, Go will show this error `function main is undeclared in the main package` when compiling.

It's conventional for the `main()` function to live in a file with the filename `main.go`.
Technically it doesn't have to, but following this convention makes the application entry point
easier to find for anyone reading the code in the future.

### Importing and Standard Library Packages

Each `.go` files in the package can import and use exported types, constants,
variables and functions from other packages including the packages in the Go [standard library](https://pkg.go.dev/std).

When importing a package from the standard library, it is required to use the full path
to the package in the standard library tree, not just the name of the package.
For example:

```
import (
    "fmt"
    "math/rand"         // Not "rand"
    "net/http"          // Not "http"
    "net/http/httptest" // Not "httptest"
)
```

Imported packages with aliases:

```
import (
    format "fmt"     // Alias 'fmt' to 'format'
    ."math"          // Dot import (use without package prefix)
    _ "database/sql" // Blank import (only run init())
)
```

Once imported, the package name becomes an accessor for the contents of that package.
Conveniently, all the packages in the Go standard library have a package name
which is the same as the final element of their import path.

If a package is imported but not used anywhere in the code, it will result in a compile-time error.
For example, if the `os` package is imported but not used in the code:

```
"os" imported and not used
```

Similarly, Go will throw compile-time error if a package is referenced in the code but not imported.
For example, if the `strconv` package is used without importing, Go will throw an error like this:

```
undefined: strconv
```

### Exported vs Unexported

- **Exported**
  (Capitalized identifiers. e.g., `User`, `ProcessData()`)
  are publicly visible and can be imported by other packages

- **Unexported**
  (Lowercase identifiers. e.g., `user`, `processData()`)
  are private and only visible to the code in the same package

Generally, it is a rule of thumb not to export things unless there is an actual reason to do so
(i.e. don't capitalize a name just because it looks nicer!).
Additionally, a `main` package should never normally be imported by anything,
so it probably shouldn't have any exported things in it.

Example of Visibility:

```go
package user

type User struct { // Exported, visible everywhere
	name string // Unexported, only visible in this package
}

func CreateUser() { // Exported function
	// Can be called from other packages
}

func validateUser() { // Unexported function
	// Only callable within this package
}
```

## Modules

### Introduction to Go Modules

Go module is a collection of Go packages stored in a file tree with a `go.mod` file at its root.
This `go.mod` file defines the module's path, which is its unique identifier,
and it also lists the specific versions of other modules that your project depends on.
Go modules solve dependency management challenges by:

- Track project dependencies
- Ensure consistent builds
- Manage package versioning

Example of a simple Go modules

```
myproject/
│
├── go.mod
├── main.go
│
└── utils/
    ├── math.go
    └── string.go
```

### File Structure

- `go.mod` file: This is the heart of a Go module.
  It's automatically created when you initialize a new module and contains:
    - The module path, which is typically the URL where your repository is located (e.g.,
      github.com/your-username/your-project).
    - The Go version your module is compatible with.
    - A list of required modules (your dependencies) and their specific versions.

  Example of a go.mod file
    ```
    module github.com/username/projectname
    
    go 1.24.0  // Go version
    
    require (
      // External dependencies
      github.com/somepackage/library v1.2.3
    )
    ```
- `go.sum` file: 
  This file is generated alongside go.mod and contains the cryptographic checksums of the direct and indirect dependencies. 
  This ensures the integrity and authenticity of your project's dependencies, 
  making sure you're using the exact same code every time you build your project. 
  You should commit both go.mod and go.sum to your version control system

### Common Commands

| Command         | Purpose                                     |
|-----------------|---------------------------------------------|
| `go mod init`   | Initialize a new module or module           |
| `go mod tidy`   | Add missing and remove unused modules       |
| `go get`        | Add or update dependencies                  |
| `go mod vendor` | Create a vendor directory with dependencies |

## Best Practices

### Package Design

- Keep packages small and focused
- Follow the Single Responsibility Principle
- Use clear, descriptive package names
- Minimize package dependencies

### Versioning

- Use semantic versioning for modules
- Publish packages on version control systems
- Use go modules to manage version constraints

## Common Mistakes

### Circular Dependencies

- Avoid creating circular dependencies between packages
- Restructure code to break circular references
- Use interfaces to decouple package interactions

### Dependency Management

- Regularly update dependencies
- Use `go mod tidy` to clean up unused packages
- Be cautious of transitive dependencies

## Practice Exercises

1. Create a small library package with utility functions
2. Build a project using multiple custom packages
3. Experiment with different import strategies
