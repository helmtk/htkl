[![Go Reference](https://pkg.go.dev/badge/helmtk.dev/code/htkl.svg)](https://pkg.go.dev/helmtk.dev/code/htkl)

# htkl

htkl is a structured template language, used by the helmtk project.

## Overview

htkl provides a clean, readable syntax for generating structured data (like YAML or JSON), with support for templates, expressions, control flow, and built-in functions. It's designed to make configuration management more maintainable and less error-prone.

```go
import "helmtk.dev/code/hktl
```

## Features

- **Structured Data**: Define objects and arrays with a clean, indentation-aware syntax
- **Templates**: Reusable templates with the `define()` and `include()` functions
- **Expressions**: Arithmetic, comparison, and logical operators
- **Control Flow**: `for` loops, `if` statements, and `with` statements for scoping
- **Variables**: `let` statements for defining reusable values
- **Functions**: Built-in functions for common operations
- **String Interpolation**: Embed expressions in strings with `${expr}` syntax
- **Pipes**: Chain operations with the pipe operator
- **Spread Operator**: Merge objects and arrays easily

## Example

```helmtk
define("double") x * 2

define("makeLabel") do
    app: app
    version: "1.0"
end

let items = [1, 2, 3]

config: {
    doubled: include("double", {x: 21})
    labels: {
        include("makeLabel", {app: "myapp"})
    }
    values: [for i, x in items do x * 2 end]
}
```

## Project Structure

- `parser/` - Lexer, parser, and AST definitions
- `runtime/` - Runtime values, scopes, and comparison logic
- `eval/` - Expression evaluator and built-in functions
- `eval/testdata/` - Test files demonstrating language features

## License

MIT License - see LICENSE file for details
