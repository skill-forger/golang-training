# Module 14: Object Relational Mapping with GORM

## Table of Contents
<ol>
	<li><a href="#objectives">Objectives</a></li>
	<li><a href="#overview">Overview</a></li>
	<li><a href="#introduction-to-gorm">Introduction to GORM</a></li>
	<li><a href="#core-concepts">Core Concepts</a></li>
	<li><a href="#installation-and-setup">Installation and Setup</a></li>
	<li><a href="#defining-models">Features</a></li>
</ol>


## Objectives

By the end of this module, you will be able to:
- Understand what an ORM is and why it is used
- Configure and connect GORM to a relational database
- Define models and map them to database tables
- Perform CRUD operations using GORM
- Work with relationships (one-to-one, one-to-many, many-to-many)
- How to use transactions, hooks, and migrations

## Overview

Relational databases are a core part of most backend systems.
While writing raw SQL provides fine-grained control, it often leads to repetitive code, tight coupling, and maintenance challenges as applications grow.

Object Relational Mapping (ORM) addresses this problem by allowing developers to interact with database records using native language constructs instead of raw SQL.

GORM is the most widely used ORM library in the Go ecosystem. It provides a developer-friendly API on top of database/sql while still allowing access to raw SQL when necessary.

## Introduction to GORM

GORM is a full-featured ORM library for Go that maps Go structs to relational database tables.

Key characteristics of GORM:
- Convention over configuration
- Struct-based model definitions
- Chainable query APIs
- Built-in support for associations and transactions
- Compatibility with major relational databases

GORM supports multiple databases, including:
- PostgreSQL
- MySQL / MariaDB
- SQLite
- SQL Server

GORM is commonly used in REST APIs, microservices, and monolithic backend applications.

## Core Concepts

Before using GORM effectively, it is important to understand its core building blocks:

`*gorm.DB`: The main database handle used to execute queries and manage transactions.

**Models**: Go structs that represent database tables.

**Conventions**: GORM automatically infers table names, column names, and relationships unless explicitly overridden.

**Sessions**: Allow configuration of query behavior (e.g. dry run, transactions, context).

**Callbacks and Hooks**: Functions that run automatically during lifecycle events (create, update, delete).

## Installation and Setup

### Install GORM and a Database Driver

Example using PostgreSQL:

```shell
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

Initialize Database Connection

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

dsn := "host=localhost user=postgres password=secret dbname=appdb port=5432 sslmode=disable"

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    panic("failed to connect to database")
}
```

The `db` object is reused across the application and should be initialized once at startup.

## Features

### Models
Models are defined as Go structs.

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Email     string `gorm:"uniqueIndex"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

GORM automatically:
- Maps struct names to table names
- Maps fields to columns
- Manages timestamps (CreatedAt, UpdatedAt)


### Create

- Create single record
```go
user := User{Name: "Bob", Age: 18, Birthday: time.Now()}
result := db.Create(&user) // pass pointer of data to Create

user.ID             // returns inserted data's primary key
result.Error        // returns error
result.RowsAffected // returns inserted records count
```

- Create multiple records
```go
users := []*User{
    {Name: "Alice", Age: 18, Birthday: time.Now()},
    {Name: "Bob", Age: 19, Birthday: time.Now()},
}
result := db.Create(users) // pass a slice to insert multiple row

result.Error        // returns error
result.RowsAffected // returns inserted records count
```

- Create a record and assign a value to the fields specified.
```go
db.Select("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("alice", 18, "2026-01-01 11:05:21.775")
```


- Create a record and ignore the values for fields passed to omit.
```go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2000-01-01 00:00:00.000", "2026-01-01 11:05:21.775")
```

- GORM supports create from **map[string]interface{}** and **[]map[string]interface{}{}**
```go
db.Model(&User{}).Create(map[string]interface{}{
    "Name": "jinzhu", "Age": 18,
})

// batch insert from `[]map[string]interface{}{}`
db.Model(&User{}).Create([]map[string]interface{}{
    {"Name": "jinzhu_1", "Age": 18},
    {"Name": "jinzhu_2", "Age": 20},
})
```

**NOTE**: When creating from map, hooks won’t be invoked, associations won’t be saved and primary key values won’t be backfilled

For more references and example: [GORM Create](https://gorm.io/docs/create.html)

### Query

#### Retrieving a single object
GORM provides `First`, `Take`, `Last` methods to retrieve a single object from the database, it adds `LIMIT 1` condition when querying the database, and it will return the error `ErrRecordNotFound` if no record is found.

```go
// Get the first record ordered by primary key
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&user)
result.RowsAffected // returns count of records found
result.Error        // returns error or nil

// check error ErrRecordNotFound
errors.Is(result.Error, gorm.ErrRecordNotFound)
```

#### Retrieving objects with primary key
```go
db.First(&user, 10)
// SELECT * FROM users WHERE id = 10;

db.First(&user, "10")
// SELECT * FROM users WHERE id = 10;

db.Find(&users, []int{1,2,3})
// SELECT * FROM users WHERE id IN (1,2,3);
```

- If the primary key is a string (for example, like a uuid), the query will be written as follows:

```go
db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```


#### String Conditions
```go
// Get first matched record
db.Where("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

// Get all matched records
db.Where("name <> ?", "jinzhu").Find(&users)
// SELECT * FROM users WHERE name <> 'jinzhu';

// IN
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;
```

#### Struct & Map Conditions
```go
// Struct
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

// Map
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

// Slice of primary keys
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```

**NOTE:** When querying with struct, GORM will only query with non-zero fields, that means if your field’s value is 0, '', false or other zero values, it won’t be used to build query conditions, for example:
```go
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu";
```
To include zero values in the query conditions, you can use a map, which will include all key-values as query conditions, for example:
```go
db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
```
For more references and example: [GORM Query](https://gorm.io/docs/query.html)



### Update

#### Update single column
When using the Model method and its value has a primary value, the primary key will be used to build the condition, for example:
```go
// Update with conditions
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// User's ID is `111`:
db.Model(&user).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update with conditions and model value
db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
```

#### Update multiple columns
Updates supports updating with struct or map[string]interface{}, when updating with struct it will only update non-zero fields by default

```go
// Update attributes with `struct`, will only update non-zero fields
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// Update attributes with `map`
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

For more references and example: [GORM Update](https://gorm.io/docs/update.html)

### Delete

```go
// Email's ID is `10`
db.Delete(&email)
// DELETE from emails where id = 10;

// Delete with additional conditions
db.Where("name = ?", "jinzhu").Delete(&email)
// DELETE from emails where id = 10 AND name = "jinzhu";
```

#### Delete with primary key
GORM allows to delete objects using primary key(s) with inline condition
```go
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);
```

#### Soft delete
If your model includes a `gorm.DeletedAt` field (which is included in `gorm.Model`), it will get soft delete ability automatically!

When calling Delete, the record won't be removed from the database, but GORM will set the DeletedAt's value to the current time, and the data is not findable with normal Query methods anymore.


```go
// user's ID is `111`
db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// Batch Delete
db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
```

If you don't want to include gorm.Model, you can enable the soft delete feature like:
```go
type User struct {
    ID      int
    Deleted gorm.DeletedAt
    Name    string
}
```
Find soft deleted records
You can find soft deleted records with Unscoped

```go
db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;
```

#### Delete permanently
You can delete matched records permanently with Unscoped
```go
db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;
```

### Querying Techniques

GORM provides a fluent, chainable API for building database queries.
Queries are composed by chaining methods on a `*gorm.DB` instance

All query methods return a new `*gorm.DB`, allowing you to chain multiple conditions.

#### Filtering with `Where`

The `Where` method adds filtering conditions.

Simple condition
```go
db.Where("email = ?", "alice@example.com").First(&user)
```

Multiple conditions
```go
db.Where("age >= ? AND active = ?", 18, true).Find(&users)
```

Struct-based filtering
```go
db.Where(&User{Active: true}).Find(&users)
```

Only non-zero fields in the struct are used as conditions.

#### Ordering Results with Order

Use `Order` to control the sorting of query results.

```go
db.Order("created_at desc").Find(&users)
```

Multiple ordering clauses:
```go
db.Order("role asc").Order("created_at desc").Find(&users)
```

#### Limiting and Offsetting Results

Restricts the number of returned records.
```go
db.Limit(10).Find(&users)
```

Skips a number of records (commonly used for pagination).

```go
db.Offset(20).Limit(10).Find(&users)
```

Pagination example
```go
page := 2
pageSize := 10

db.Offset((page-1)*pageSize).
    Limit(pageSize).
    Order("created_at desc").
    Find(&users)
```

#### Selecting Specific Columns

Use `Select` to limit returned columns.

```go
db.Select("id", "name").Find(&users)
```

This is useful for:
- Reducing payload size
- Optimizing query performance

#### Counting Records
```go
var count int64
db.Model(&User{}).Where("active = ?", true).Count(&count)
```

Common use cases:
- Pagination metadata
- Reporting

#### `IN` Queries

```go
db.Where("id IN ?", []int{1, 2, 3}).Find(&users)
```

#### LIKE Queries
```go
db.Where("name LIKE ?", "%john%").Find(&users)
```

#### Combining Conditions
```go
db.Where("age > ?", 18).
    Where("active = ?", true).
    Order("created_at desc").
    Find(&users)
```

Each Where call adds to the existing query.

#### Scopes (Reusable Query Logic)

Scopes allow reusable query fragments.

```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}
```

Usage:
```go
db.Scopes(ActiveUsers).Find(&users)
```

Scopes are useful for:
- Shared filtering logic
- Clean repository code
- Consistent query behavior

#### Raw SQL Queries

GORM allows falling back to raw SQL when needed.

```go
db.Raw("SELECT * FROM users WHERE age > ?", 30).Scan(&users)
```

This is useful for:
- Complex joins
- Database-specific features
- Performance-critical queries


##################################

### Relationships and Associations

Relational databases model data using relationships between tables.
GORM provides first-class support for defining and working with these relationships directly through Go structs.

#### One-to-One Relationship

A one-to-one relationship means one record is associated with exactly one other record.

Example: User and Profile
```go

type User struct {
    ID      uint
    Name    string
    Profile Profile
}

type Profile struct {
    ID     uint
    UserID uint
    Bio    string
}
```
UserID acts as the foreign key

GORM infers the relationship automatically

To create records:

```go
user := User{
    Name: "Alice",
    Profile: Profile{
        Bio: "Software Engineer",
    },
}
db.Create(&user)
```


#### One-to-Many Relationship

A one-to-many relationship means one record can be associated with multiple records.

Example: User and Posts
```go
type User struct {
    ID    uint
    Name  string
    Posts []Post
}

type Post struct {
    ID     uint
    Title  string
    UserID uint
}
```

- One User → many Post
- UserID is the foreign key in Post

#### Query with associations:

```go
var user User
db.Preload("Posts").First(&user, 1)
```

#### Many-to-Many Relationship

A many-to-many relationship requires a join table.

Example: Users and Roles
```go
type User struct {
    ID    uint
    Name  string
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID   uint
    Name string
}
```

GORM automatically manages the join table user_roles.

Create associations:

```go

user := User{
    Name: "Bob",
    Roles: []Role{
        {Name: "Admin"},
        {Name: "Editor"},
    },
}

db.Create(&user)
```

#### Foreign Keys

Foreign keys define how tables are connected.

By default, GORM:

- Uses <StructName>ID as the foreign key
- Uses the primary key of the parent table 
 
You can customize foreign keys explicitly:

```go
type Order struct {
    ID     uint
    UserID uint
    User   User `gorm:"foreignKey:UserID;references:ID"`
}
```

This is useful when:
- Database schema already exists
- Naming conventions differ
- Composite keys are used

#### Preloading Associations

By default, GORM does not load associations automatically. Use `Preload` to eagerly load related data.

Single association
```go
db.Preload("Profile").First(&user)
```

Multiple associations
```go
db.Preload("Posts").Preload("Roles").Find(&users)
```

Nested preloading
```go
db.Preload("Posts.Comments").Find(&users)
```

#### Conditional Preloading

You can apply conditions to preloaded associations.

```go
db.Preload("Posts", "published = ?", true).
First(&user)
```

This loads only matching related records.

#### Association Modes
Association mode allows you to manipulate relationships directly without loading the parent object fully.

Access association handler
```go
db.Model(&user).Association("Roles")
```

Append association
```go
db.Model(&user).Association("Roles").Append(&role)
```

Replace associations
```go
db.Model(&user).Association("Roles").Replace(&newRoles)
```

Delete association
```go
db.Model(&user).Association("Roles").Delete(&role)
```

Clear all associations
```go
db.Model(&user).Association("Roles").Clear()
```

Association mode is useful for:
- Updating join tables
- Managing relationships independently
- Avoiding unnecessary queries

Deleting and Associations
- By default deleting a parent record does not delete associated records

You can enable cascading behavior:
```go
type User struct {
    ID    uint
    Posts []Post `gorm:"constraint:OnDelete:CASCADE;"`
}
```

Use this carefully to avoid accidental data loss.

### Migrations and Schema Management

Database schema management is a critical part of application development.
As applications evolve, tables, columns, indexes, and constraints must be updated in a controlled and predictable way.

GORM provides built-in tools to help manage schema changes, most notably through auto migration. However, understanding when and how to use these tools is essential to avoid data loss or production issues.

This section covers:
- Auto migration
- Managing schema changes
- Index creation
- Constraints
- When auto-migrate is safe and unsafe


#### Auto Migration

Auto migration allows GORM to automatically create and update database schemas based on model definitions.

Basic Auto Migration
```go
db.AutoMigrate(&User{}, &Post{})
```

Auto migration will:
- Create tables if they do not exist
- Add missing columns
- Add missing indexes
- Add foreign key constraints (when supported)

Auto migration will NOT:
- Drop existing tables
- Remove existing columns
- Change column types
- Rename columns

This design minimizes accidental data loss.

#### Managing Schema Changes

As applications evolve, schema changes may include:
- Adding new fields
- Changing data types
- Renaming columns
- Removing unused columns

Adding a New Column
```go
type User struct {
    ID    uint
    Name  string
    Email string
    Age   int //new field
}
```
Running AutoMigrate will automatically add the Age column.

Changing or Removing Columns (Manual Process)

GORM does not automatically:
- Change column types
- Remove columns
- These changes should be handled using:
- Manual SQL migrations
- Dedicated migration tools
- Versioned migration scripts

Example (manual SQL):
```sql
ALTER TABLE users ALTER COLUMN age TYPE bigint;
```

This separation ensures schema changes are intentional and reviewable.

#### Index Creation

Indexes improve query performance and are essential for frequently queried columns.

Creating Indexes with Tags

```go
type User struct {
    ID    uint
    Email string `gorm:"index"`
}
```

Unique Index
```go
type User struct {
    Email string `gorm:"uniqueIndex"`
}
```

Composite Index
```go
type Order struct {
    UserID uint `gorm:"index:idx_user_status"`
    Status string `gorm:"index:idx_user_status"`
}
```

Indexes are automatically created during auto migration.

#### Constraints

Constraints enforce data integrity at the database level.

Foreign Key Constraints
```go
type Post struct {
    ID     uint
    UserID uint
    User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
```

Supported actions include:
- CASCADE
- SET NULL
- RESTRICT
- NO ACTION

Check Constraints
```go
type Product struct {
    Price float64 `gorm:"check:price > 0"`
}
```

Check constraints ensure values meet specific conditions.

Not Null Constraint
```go
type User struct {
    Name string `gorm:"not null"`
}
```

#### When Auto-Migrate Is Safe
Auto migration is generally safe and recommended when:
- In early development
- The schema is evolving rapidly
- You are adding new tables or columns
- You control the database entirely
- The application is internal or non-critical

Typical use cases:
- Prototyping
- Internal tools
- Side projects
- Early-stage startups

#### When Auto-Migrate Is Unsafe
Auto migration should be used with caution or avoided when:
- Running on production databases
- Schema changes involve data transformation
- Columns need to be renamed or removed
- Data types must change
- Multiple services share the same database

In these cases, prefer:
- Versioned migration tools
- Explicit SQL scripts
- Manual review and rollback plans

### Transactions and Hooks

Modern applications often require multiple database operations to succeed or fail as a single unit.
GORM provides first-class support for database transactions and model lifecycle hooks to help enforce consistency, integrity, and business rules.

#### Transactions
A transaction groups multiple operations into a single, atomic unit. It guarantees
- Atomicity – all operations succeed or none do
- Consistency – database remains in a valid state
- Isolation – intermediate states are not visible
- Durability – committed changes persist

GORM provides a high-level Transaction helper.

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&order).Error; err != nil {
        return err
    }

    if err := tx.Create(&payment).Error; err != nil {
        return err
    }

    return nil // commit
})
```

- Returning nil → commit
- Returning an error → rollback

For some scenarios, transactions can be managed manually.
```go
tx := db.Begin()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}

return tx.Commit().Error
```

GORM supports nested transactions, you can rollback a subset of operations performed within the scope of a larger transaction
```go
db.Transaction(func(tx *gorm.DB) error {
  tx.Create(&user1)
  tx.Transaction(func(tx2 *gorm.DB) error {
    tx2.Create(&user2)
    return errors.New("rollback user2") // Rollback user2
  })
  tx.Transaction(func(tx3 *gorm.DB) error {
    tx3.Create(&user3)
    return nil
  })
  return nil
})
// Commit user1, user3
```
GORM provides `SavePoint`, `RollbackTo` to save points and roll back to a savepoint, for example:

```go
tx := db.Begin()
tx.Create(&user1)

tx.SavePoint("sp1")
tx.Create(&user2)
tx.RollbackTo("sp1") // Rollback user2

tx.Commit() // Commit user1
```

#### Hooks
Hooks allow logic to run automatically at specific lifecycle stages.

Common hooks
- BeforeCreate
- AfterCreate
- BeforeUpdate
- AfterUpdate
- BeforeSave
- AfterSave
- BeforeDelete
- AfterDelete
- AfterFind

**Example**
```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.Email = strings.ToLower(u.Email)
    return nil
}
```
- Runs automatically before inserting a record
- Can modify data
- Return an error to stop the operation

**Example**
```go
func (o *Order) AfterUpdate(tx *gorm.DB) error {
    log.Println("Order updated:", o.ID)
    return nil
}
```
Hooks run inside the transaction when applicable.

#### Common Use Cases

Transactions are commonly used for
- Multi-table writes
- Data migrations
- Audit logging

Hooks are commonly used for
- Data normalization
- Automatic timestamps
- Audit trails
- Validation logic

For more references and examples
- [GORM Hooks](https://gorm.io/docs/hooks.html)
- [GORM Transaction](https://gorm.io/docs/transactions.html)

### Recommended Resources for ORM

1. **Official Documentation & Tutorials**
    - [GORM Documentation](https://gorm.io/docs/)
    - [SQL Performance Best Practices](https://use-the-index-luke.com/)

2. **Books & Courses**
    - "SQL Performance Explained" by Markus Winand.
    - "GORM for Go Developers" Udemy course.

3. **Open Source Examples**
    - [RealWorld Example App with GORM](https://github.com/gothinkster/golang-gin-realworld-example-app)