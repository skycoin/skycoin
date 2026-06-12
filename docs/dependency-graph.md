# Dependency Graph

The Go package dependency graph for the Skycoin node, made with
[goda](https://github.com/loov/goda).

Regenerate with:

```bash
go run github.com/loov/goda@latest graph github.com/skycoin/skycoin/... | dot -Tsvg -o docs/skycoin-goda-graph.svg
```

![Skycoin dependency graph](skycoin-goda-graph.svg){ loading=lazy }
