# Adaptive Card Builder for Go

A simple, strongly-typed Go package to construct **Microsoft Teams Adaptive Cards**.  
This package allows you to build Adaptive Cards programmatically using Go structs instead of manually crafting JSON.

---

## Features

- Build cards with:
  - `TextBlock`
  - `Container`
  - `FactSet` and `Fact`
  - `Action` buttons (`OpenUrl`, etc.)
- Support for nested elements (`Container` inside `Container`)
- Fluent API with **receiver methods** (`AddBody`, `AddItem`, `AddAction`)
- JSON output ready to post to Teams via Power Automate or webhook
- Optional Teams mentions support (`MSTeamsInfo` / `Entities`)
- Strongly typed — reduces errors compared to raw JSON strings

---

## Installation

```bash
go get github.com/luisdibdin/adaptivecard
```
