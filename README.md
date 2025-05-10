# Collection Package

A thread-safe generic collection implementation in Go, similar to a Map but with extended functionality.

## Features

- **Thread-safe** operations with `sync.RWMutex`
- **Generic** support for any comparable key and any value type
- **Rich API** with over 30 methods for data manipulation
- **Immutable operations** where appropriate (methods returning new collections)
- **Functional-style** methods like Map, Filter, Reduce, etc.

## Installation

```bash
go get github.com/icehuntmen/collection
```

## Usage

### Basic Operations

```go
// Create a new collection
coll := collection.New[string, int]()

// Add items
coll.Set("a", 1)
coll.Set("b", 2)
coll.Set("c", 3)

// Get items
val, exists := coll.Get("a") // 1, true

// Check existence
hasKey := coll.Has("b") // true

// Delete items
coll.Delete("c")

// Get size
count := coll.Size() // 2
```

### Functional Operations

```go
// Filter elements
filtered := coll.Filter(func(v int, k string) bool {
    return v > 1
}) // Contains only {"b": 2}

// Map values
mapped := coll.Map(func(v int, k string) int {
    return v * 2
}) // {"a": 2, "b": 4}

// Reduce values
sum := coll.Reduce(func(acc, v int, k string) int {
    return acc + v
}, 0) // 3

// ForEach iteration
coll.ForEach(func(v int, k string) {
    fmt.Printf("%s: %d\n", k, v)
})
```

### Collection Operations

```go
// Union of two collections
union := coll1.Union(coll2)

// Intersection of two collections
intersection := coll1.Intersection(coll2)

// Difference between collections
diff := coll1.Difference(coll2)

// Symmetric difference
symDiff := coll1.SymmetricDifference(coll2)
```

### Sorting

```go
// Sort in-place
coll.Sort(func(v1, v2 int, k1, k2 string) bool {
    return v1 < v2 // ascending order
})

// Create sorted copy
sorted := coll.ToSorted(func(v1, v2 int, k1, k2 string) bool {
    return k1 > k2 // descending order by key
})
```

### Advanced Features

```go
// Ensure a value exists (lazy initialization)
val := coll.Ensure("newKey", func(k string) int {
    return 42 // Default value if key doesn't exist
})

// Merge two collections with different value types
result := coll1.Merge(
    coll2,
    func(v int, k string) collection.Keep[string] { /* whenInSelf */ },
    func(u string, k string) collection.Keep[string] { /* whenInOther */ },
    func(v int, u string, k string) collection.Keep[string] { /* whenInBoth */ },
)

// Group items by key
groups := collection.GroupBy(items, func(item Item) KeyType {
    return item.SomeField
})
```

## Performance Considerations

- All operations are thread-safe but may impact performance in high-concurrency scenarios
- Methods like `Sort`, `Values`, and `Keys` have O(n) complexity
- Underlying storage is a standard Go map, so most operations have O(1) complexity

## License

MIT License