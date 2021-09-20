# Struct to Orderedmap

A sample that converts any `struct` into an ordered map.

## When is it useful

- To ignore the JSON `omitempty` tag that already defined in `struct`.
- To add some additional fields.

## Warning

- Only `struct` or the pointer of `struct` is allowed as an argument.
- If the field contains a map, the key must be string.
- Performance is not high as it uses the `reflect` package. Please use it in tools.
