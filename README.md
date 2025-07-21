# go-concurrently

This repo produces the `concurrently` binary that allows for the parallel execution of several commands.

## Usage

Usage is simple - provide each command as an argument to concurrently, and they'll be run together.

```shell
concurrently "echo 'a'" "echo 'b'"
```
