# Raptor

[![tag](https://img.shields.io/github/v/tag/maddiesch/go-raptor.svg)](https://github.com/maddiesch/go-raptor/releases)
[![GoDoc](https://godoc.org/github.com/maddiesch/go-raptor?status.svg)](https://pkg.go.dev/github.com/maddiesch/go-raptor)
![Build Status](https://github.com/maddiesch/go-raptor/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/maddiesch/go-raptor/graph/badge.svg?token=K62CECQYJK)](https://codecov.io/gh/maddiesch/go-raptor)
[![License](https://img.shields.io/github/license/maddiesch/go-raptor)](./LICENSE)

Raptor is a lightweight SQLite3 wrapper for Go programming language. It provides an easy-to-use interface for working with SQLite3 databases.

## Features

- Simple and easy-to-use interface.
- Safe and secure, uses parameterized queries to prevent SQL injection attacks.
- Supports transactions with savepoints.
- Allows custom logging of SQL queries.

## Getting Started

To get started, install the package using `go get`.

```bash
go get github.com/maddiesch/go-raptor
```

Then, import the package in your code:

```go
import "github.com/maddiesch/go-raptor"
```

You can then create a new connection to a SQLite3 database using the `New` function:

```go
db, err := raptor.New("database.db")
if err != nil {
    // handle error
}
defer db.Close()
```

## Usage

You can execute SQL queries on the database using the `Exec` and `Query` methods:

```go
// Exec a query that doesn't return any rows
result, err := db.Exec(context.Background(), "INSERT INTO users (name, age) VALUES (?, ?)", "John Doe", 42)
if err != nil {
    // handle error
}

// Query rows from the database
rows, err := db.Query(context.Background(), "SELECT name, age FROM users WHERE age > ?", 30)
if err != nil {
    // handle error
}
defer rows.Close()

// Iterate over the rows
for rows.Next() {
    var name string
    var age int
    if err := rows.Scan(&name, &age); err != nil {
        // handle error
    }
    // do something with name and age
}
if err := rows.Err(); err != nil {
    // handle error
}
```

You can also execute a query that returns a single row using the `QueryRow` method:

```go
// Query a single row from the database
row := db.QueryRow(context.Background(), "SELECT name, age FROM users WHERE id = ?", 1)

// Scan the values from the row
var name string
var age int
if err := row.Scan(&name, &age); err != nil {
    // handle error
}
// do something with name and age
```

## Transactions

Raptor supports transactions with savepoints. You can execute a transaction using the `Transact` method:

```go
err := db.Transact(context.Background(), func(tx raptor.DB) error {
    // Exec some queries
    if _, err := tx.Exec(context.Background(), "INSERT INTO users (name, age) VALUES (?, ?)", "John Doe", 42); err != nil {
        return err
    }
    if _, err := tx.Exec(context.Background(), "UPDATE users SET age = ? WHERE name = ?", 43, "John Doe"); err != nil {
        return err
    }
    return nil
})
if err != nil {
    // handle error
}
```

Inside the transaction, you can use the Exec, Query, and QueryRow methods as usual.
