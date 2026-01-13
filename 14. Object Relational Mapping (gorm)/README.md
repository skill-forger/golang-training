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

### Defining Models
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

When calling Delete, the record wont’t be removed from the database, but GORM will set the DeletedAt‘s value to the current time, and the data is not findable with normal Query methods anymore.


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

If you don’t want to include gorm.Model, you can enable the soft delete feature like:
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





Topics:

Struct-to-table mapping

Field tags (column, type, index, unique)

Embedded structs

Custom table names

JSON vs DB field separation

### Database Operations (CRUD)

Outline:

Creating records

Reading records

Updating records

Deleting records (hard vs soft delete)

Batch operations

### Querying Techniques

Topics:

Where, First, Find

Ordering and limiting

Selecting specific columns

Raw SQL with GORM

Pagination patterns

Scopes

### Relationships and Associations

Cover:

One-to-One

One-to-Many

Many-to-Many

Foreign keys

Preloading associations

Association modes

### Migrations and Schema Management

Topics:

Auto migration

Managing schema changes

Index creation

Constraints

When auto-migrate is safe / unsafe

### Transactions and Hooks

Outline:

Database transactions

Commit / rollback patterns

Model hooks (BeforeCreate, AfterUpdate, etc.)

Use cases and pitfalls