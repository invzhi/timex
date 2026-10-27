# timex

[![Go Reference](https://pkg.go.dev/badge/github.com/invzhi/timex.svg)](https://pkg.go.dev/github.com/invzhi/timex)
[![codecov](https://codecov.io/gh/invzhi/timex/branch/main/graph/badge.svg?token=I2M6JCGY84)](https://codecov.io/gh/invzhi/timex)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#date-and-time)

📅 A Go package that extends the standard library time with dedicated date and time-of-day types.

`timex` gives you two types missing from Go's standard `time` package: `Date` and `TimeOfDay`. Use them for birthdays, due dates, opening hours, and other values that don't identify a moment on a timeline.

## Why use timex?

An invoice is due on **2026-10-26**. A shop opens at **09:30**. Representing either with `time.Time` means supplying extra information: a clock time for the invoice, or a date for the opening hour, plus a location. Using strings leaves validation, comparison, and arithmetic to your application.

`timex` lets your types express exactly what your application means:

```go
type Invoice struct {
    DueDate timex.Date `json:"due_date"`
}

type Shop struct {
    OpensAt timex.TimeOfDay `json:"opens_at"`
}
```

- **Keep calendar values independent of time zones.** Neither type stores a location or UTC offset, so there is no implicit time zone conversion to shift a date or clock time.
- **Validate values when you create them.** Constructors reject invalid dates and clock times, including February 29 in a non-leap year.
- **Do arithmetic in the units you mean.** Add calendar days, count days between dates, compare opening hours, and get the day rollover when clock arithmetic crosses midnight.
- **Use the same types in Go, JSON, and SQL.** Built-in JSON support and `database/sql` interfaces reduce conversion code at application boundaries.
- **Keep runtime dependencies small.** The library imports only the Go standard library; third-party dependencies are used by its tests.

| What your value represents | Use |
| --- | --- |
| A calendar day: birthday, booking date, invoice due date | `timex.Date` |
| A clock time: opening hour, daily cutoff, recurring appointment time | `timex.TimeOfDay` |
| A specific instant: event timestamp, audit log entry, expiration instant | `time.Time` |
| An elapsed amount of time | `time.Duration` |

## Install

Requires Go 1.20 or later.

```sh
go get github.com/invzhi/timex
```

## Quick start

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/invzhi/timex"
)

func main() {
    // Calendar arithmetic: no midnight or time zone needed.
    issued := timex.MustNewDate(2026, 10, 26)
    due := issued.AddDays(30)
    fmt.Println(due)             // 2026-11-25
    fmt.Println(due.Sub(issued)) // 30 (days)

    // Clock arithmetic reports how many days you cross.
    starts := timex.MustNewTimeOfDay(23, 30, 0, 0)
    days, ends := starts.AddDuration(2 * time.Hour)
    fmt.Println(days, ends) // 1 01:30:00

    data, err := json.Marshal(struct {
        DueDate timex.Date      `json:"due_date"`
        OpensAt timex.TimeOfDay `json:"opens_at"`
    }{
        DueDate: due,
        OpensAt: timex.MustNewTimeOfDay(9, 30, 0, 0),
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
    // {"due_date":"2026-11-25","opens_at":"09:30:00"}
}
```

Use `NewDate` and `NewTimeOfDay` when handling input that may be invalid; they return errors. The `Must…` variants panic on invalid input and are convenient for known constants.

## Parse and format

Layouts use tokens such as `YYYY-MM-DD` and `HH:mm:ss`, rather than Go's reference-time layouts.

```go
date, err := timex.ParseDate(timex.RFC3339Date, "2026-10-26")
if err != nil {
    return err
}
fmt.Println(date.Format("MMM D, YYYY")) // Oct 26, 2026

opens, err := timex.ParseTimeOfDay("h:mm A", "9:30 AM")
if err != nil {
    return err
}
fmt.Println(opens) // 09:30:00
```

`TimeOfDay` supports nanosecond precision, including fractional seconds in parsing, formatting, and JSON. See the [API reference](https://pkg.go.dev/github.com/invzhi/timex) for all layout tokens and methods, including comparisons, calendar arithmetic, weekdays, and ISO week numbers.

## JSON and databases

Both types marshal to JSON strings and unmarshal from the same formats:

```json
{"due_date":"2026-10-26","opens_at":"09:30:00"}
```

`Date` and `TimeOfDay` implement `sql.Scanner` and `driver.Valuer`, so you can pass them as query arguments and scan directly into them:

```go
var due timex.Date
err := db.QueryRowContext(ctx,
    "SELECT due_date FROM invoices WHERE id = $1", id,
).Scan(&due)
```

The example uses PostgreSQL's placeholder syntax. SQL scanning accepts `string`, `[]byte`, or `time.Time` values from your driver; values are written as date or clock-time strings. Use `NullDate` and `NullTimeOfDay` for nullable SQL columns. `TimeOfDay` represents a time within a day (00:00 to before 24:00), so use a duration type for SQL values representing elapsed time or more than 24 hours.

## Working with `time.Time`

When extracting a date or clock time from an instant, first put that instant in the location your application means:

```go
local := instant.In(location)
date := timex.DateFromTime(local)
clock := timex.TimeOfDayFromTime(local)

midnight := date.Time(location) // time.Time at the start of this date
```

`Today(location)` and `TimeOfDayNow(location)` provide the current date and clock time in a chosen location. These types store calendar values, not time zone rules: scheduling an actual occurrence still requires a date, a location, and a policy for daylight saving transitions.

The zero value of `Date` is `0001-01-01`; the zero value of `TimeOfDay` is midnight (`00:00:00`). Neither zero value means “missing.”

## Learn more

- [Go proposal for a date type](https://github.com/golang/go/issues/21365)
- [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)
