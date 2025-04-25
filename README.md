# timex

[![Go Reference](https://pkg.go.dev/badge/github.com/invzhi/timex.svg)](https://pkg.go.dev/github.com/invzhi/timex)
[![Go Report Card](https://goreportcard.com/badge/github.com/invzhi/timex)](https://goreportcard.com/report/github.com/invzhi/timex)
[![codecov](https://codecov.io/gh/invzhi/timex/branch/main/graph/badge.svg?token=I2M6JCGY84)](https://codecov.io/gh/invzhi/timex)

📅 A Go package for working with date.

## Why use timex?

- Self-contained type `Date`, focusing on date-specific operations.
- Don't rely on type `time.Time`, avoid timezone-related issues.
- Represent date format like `YYYY-MM-DD` instead of `2006-01-02`.
- Working with date type of MySQL or PostgreSQL directly.
- Fast method implementations of `Date` relative to `time.Time`.

| Method        | package timex | package time |
|:--------------|--------------:|-------------:|
| OrdinalDate   |    4.07 ns/op |   6.28 ns/op |
| Date          |    6.81 ns/op |   8.04 ns/op |
| WeekDay       |    0.31 ns/op |   2.80 ns/op |
| ISOWeek       |    7.52 ns/op |   9.80 ns/op |
| Add           |   12.23 ns/op |  24.18 ns/op |
| AddDays       |    0.31 ns/op |   3.12 ns/op |
| Sub           |    0.31 ns/op |   6.53 ns/op |
| Parse         |   34.72 ns/op |  55.47 ns/op |
| Format        |   28.05 ns/op |  59.46 ns/op |
| MarshalJSON   |   27.53 ns/op |  41.42 ns/op |
| UnmarshalJSON |   10.89 ns/op |  52.11 ns/op |

## Features

- Fully-implemented type `Date`.
    - Zero value: January 1, year 1. Align with type `time.Time`.
    - Working with standard library: conversion with type `time.Time`.
    - Parsing & Formatting: conversion with formatted strings.
    - Getter: get year, quarter, month, day of year, day of month, day of week.
    - Manipulation: addition and subtraction with years, months, days.
    - Comparison: comparing dates with `Before`, `After`, `Equal`.
    - Serialization: JSON and database.

## Getting Started

```
go get github.com/invzhi/timex
```

## Reference

- [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)
- [Date Type Proposal](https://github.com/golang/go/issues/21365)
