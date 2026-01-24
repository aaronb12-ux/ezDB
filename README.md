# Embedded Database in Go

This project is a simple embedded database implemented from scratch in Go. It focuses on low-level storage, disk I/O, and indexing using a B+ tree.

## Overview

The database uses a page-based file manager to persist data on disk and a B+ tree for efficient indexing. All data is stored on disk and loaded into memory only when needed.

## Architecture

### File Manager
- Handles disk I/O using fixed-size blocks (pages)
- Each page is a raw byte buffer loaded from a disk block
- Pages are read from and written back to disk at specific block offsets
- All disk data flows through pages

### Pages and Serialization
- Pages store raw bytes only
- Data is deserialized from a page into structured in-memory node objects
- Nodes are modified in memory and then serialized back into pages
- Pages are flushed back to disk to persist changes

### B+ Tree Index
- Uses a disk-backed B+ tree for indexing
- Supports `get`, `put`, and `delete` operations
- Tree traversal locates the correct disk block for each operation
- Leaf nodes store key-value pairs; internal nodes store keys and child pointers
- Operations run in `O(log n)` time
