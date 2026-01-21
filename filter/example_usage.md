# Filter Package - Operation-Based Filtering

## Overview

The filter package now uses a type-safe Operation-based API that allows you to inspect and filter filesystem operations based on their specific type and parameters.

## Basic Usage

### Using PathPredicate (Backward Compatible)

```go
// Filter based on path only
fs := filter.NewFs(baseFs, filter.PathPredicate(func(path string) bool {
    return strings.HasSuffix(path, ".txt")
}))
```

### Using Operation Type Switch (New API)

```go
// Filter based on operation type and parameters
fs := filter.NewFs(baseFs, func(op filter.Operation) bool {
    switch v := op.(type) {
    case filter.OpenOp:
        // Allow opening only .txt files
        return strings.HasSuffix(v.Name, ".txt")
    
    case filter.OpenFileOp:
        // Block write operations
        return v.Flag & os.O_WRONLY == 0
    
    case filter.RemoveOp:
        // Block removal of important files
        return !strings.Contains(v.Name, "important")
    
    case filter.RenameOp:
        // Log rename operations
        log.Printf("Renaming %s to %s", v.Oldname, v.Newname)
        return true
    
    default:
        // Allow all other operations
        return true
    }
})
```

## Available Operation Types

- `ChmodOp` - Change file permissions
- `ChownOp` - Change file ownership
- `ChtimesOp` - Change file times
- `CreateOp` - Create new file
- `MkdirOp` - Create directory
- `MkdirAllOp` - Create directory tree
- `OpenOp` - Open file (read-only)
- `OpenFileOp` - Open file with flags
- `RemoveOp` - Remove file
- `RemoveAllOp` - Remove recursively
- `RenameOp` - Rename/move file
- `StatOp` - Get file info
- `ReaddirOp` - Read directory contents
- `ReaddirnamesOp` - Read directory names

## Migration Guide

Old API:
```go
fs := filter.NewFs(baseFs, func(path string) bool {
    return allowedPath(path)
})
```

New API (equivalent):
```go
fs := filter.NewFs(baseFs, filter.PathPredicate(func(path string) bool {
    return allowedPath(path)
}))
```

New API (with operation inspection):
```go
fs := filter.NewFs(baseFs, func(op filter.Operation) bool {
    return allowedOperation(op)
})
```
