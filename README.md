# timex

[![Go Reference](https://pkg.go.dev/badge/github.com/invzhi/timex.svg)](https://pkg.go.dev/github.com/invzhi/timex)
[![Go Report Card](https://goreportcard.com/badge/github.com/invzhi/timex)](https://goreportcard.com/report/github.com/invzhi/timex)
[![codecov](https://codecov.io/gh/invzhi/timex/branch/main/graph/badge.svg?token=I2M6JCGY84)](https://codecov.io/gh/invzhi/timex)

📅 A Go package that extends the standard library time with dedicated date and time-of-day types.

## Why use timex?

`timex` is a lightweight and efficiently designed Go package that provides dedicated `Date` and `TimeOfDay` types. It's built with 100% unit test coverage and has no third-party dependencies, making it a reliable choice for time-related operations.

- **Missing Standard Library Types:** Go's `time` package lacks distinct types for handling just a date or just a time of day. `timex` introduces `Date` and `TimeOfDay` to bridge this gap, offering a more intuitive and focused API for these specific use cases.
- **Simplicity and Clarity:** Our `Date` type focuses exclusively on date operations (YYYY-MM-DD), avoiding timezone-related issues inherent in `time.Time` when only a calendar date matters. Similarly, `TimeOfDay` provides a clear HH:mm:ss format for operations strictly on time, independent of any particular date.
- **Lightweight & Efficient:** `timex` is designed to be lean, with fast method implementations that often outperform `time.Time` for specific date and time-of-day calculations. It's built without any external dependencies, ensuring a minimal footprint.
- **Database & JSON Compatibility:** Work directly with `DATE` and `TIME` types in databases like MySQL or PostgreSQL, and effortlessly serialize/deserialize these types to and from JSON.
- **Reliability:** With 100% unit test coverage, you can be confident in the package's correctness and stability.

## Getting Started

```
go get github.com/invzhi/timex
```

## Reference

- [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)
- [Date Type Proposal](https://github.com/golang/go/issues/21365)
