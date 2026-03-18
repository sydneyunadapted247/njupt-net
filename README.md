# njupt-net-cli

`njupt-net-cli` is the Go CLI kernel project for NJUPT network access workflows.

## Repository purpose

This repository exists to implement the **CLI kernel** for the NJUPT network stack based on the finalized reverse-engineering specification.

It is **not** the place for:
- historical reverse-engineering notes
- exploratory scripts
- browser automation experiments
- GUI code
- daemon/service orchestration
- platform-specific Wi-Fi automation

## Source of truth

The only protocol and behavior source of truth is:

- `docs/Final-SSOT.md`

Implementation task constraints are defined in:

- `docs/IMPLEMENTATION-TASK.md`

If any older note, script, experiment, or memory conflicts with `docs/Final-SSOT.md`, the SSOT wins.

## Design principle

This project should follow a **kernel + outer layer** architecture.

### Kernel
The kernel contains:
- protocol/domain primitives
- state models
- endpoint semantics
- write-operation verification rules
- error/capability classification
- SelfService and Portal protocol logic

### Outer layer
The outer layer contains:
- CLI command wiring
- config loading
- output formatting
- workflow orchestration
- safe confirmations
- optional future integrations

The kernel must remain stable even if the CLI surface evolves.

## Implementation status model

The code should preserve the certainty model from the SSOT:

- `confirmed` — behavior is established and can be implemented directly
- `guarded` — implementable, but must use conservative handling
- `blocked` — not fully proven; do not pretend success semantics are known

## Development goal

The objective is to build a complete, maintainable, production-oriented Go CLI for the confirmed and safely implementable guarded behavior in the SSOT.

## Build

Once implementation exists:

```bash
go build ./...
```

## Notes for implementers

- Do not infer protocol truth from historical scripts unless already incorporated into the SSOT.
- Do not collapse guarded or blocked behaviors into confirmed success paths.
- Do not let CLI/UI concerns leak into kernel protocol logic.
- Prefer explicit diagnostics over hidden assumptions.
