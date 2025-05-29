# timex

[![Go Reference](https://pkg.go.dev/badge/github.com/invzhi/timex.svg)](https://pkg.go.dev/github.com/invzhi/timex)
[![Go Report Card](https://goreportcard.com/badge/github.com/invzhi/timex)](https://goreportcard.com/report/github.com/invzhi/timex)
[![codecov](https://codecov.io/gh/invzhi/timex/branch/main/graph/badge.svg?token=I2M6JCGY84)](https://codecov.io/gh/invzhi/timex)

📅 A Go package for working with date and time of day.

## Why use timex?

- Self-contained type `Date`, focusing on date-specific operations.
  - Don't rely on type `time.Time`, avoid timezone-related issues.
  - Represent date format like `YYYY-MM-DD` instead of `2006-01-02`.
  - Working with date type of MySQL or PostgreSQL directly.
  - Fast method implementations of `Date` relative to `time.Time`.
- Self-contained type `TimeOfDay`, without date.
  - Simple format like `HH:mm:ss` for clarity.
  - Working with date type of MySQL or PostgreSQL directly.
  - Efficient manipulation of hours, minutes, seconds, and nanoseconds.

## Features

- Fully-implemented type `Date`:
    - Zero value: January 1, year 1. Align with type `time.Time`.
    - Working with standard library: conversion with type `time.Time`.
    - Parsing & Formatting: conversion with formatted strings.
    - Getter: get year, quarter, month, day of year, day of month, day of week.
    - Manipulation: addition and subtraction with years, months, days.
    - Comparison: comparing dates with `Before`, `After`, `Equal`.
    - Database serialization and deserialization.
    - JSON serialization and deserialization.
- Fully-implemented type `TimeOfDay`:
    - Working with standard library: conversion with type `time.Time`.
    - Parsing & Formatting: conversion with formatted strings.
    - Getter: get hour, minute, second, nanosecond.
    - Manipulation: addition and subtraction with hours, minutes, seconds, nanoseconds.
    - Comparison: comparing time of days with `Before`, `After`, `Equal`.
    - Database serialization and deserialization.
    - JSON serialization and deserialization.

## Getting Started

```
go get github.com/invzhi/timex
```

## Reference

- [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)
- [Date Type Proposal](https://github.com/golang/go/issues/21365)
