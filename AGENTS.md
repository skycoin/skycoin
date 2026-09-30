# AGENTS.md

Guidance for AI coding agents working in this repository.

These rules are combined from the agent and AI policies of the upstream
projects this one contributes to and depends on.

- tinygo-org/tinygo and tinygo-org/net, `AGENTS.md`
- mvdan/sh, `AGENTS.md`
- u-root/u-root, `AI_POLICY.md`, itself adapted from Delve and Ghostty

`DEVELOPMENT.md` is the human contributor guide and stays authoritative for
anything the two documents both cover.

## There are humans here

Every issue, pull request and review comment is read by a person who is
donating the time to read it. Work that has not been checked moves the cost of
checking it onto them.

## Communication style

Use simple direct language everywhere outside the code. That includes commit
messages, pull request titles and descriptions, issues, and review comments.

Lead with the change and the reason for it. Two or three short paragraphs is
usually the whole description. Include a measurement, a table or a log excerpt
only where it is the evidence for a claim, not one for every claim.

Avoid extra colons, semicolons and dashes. Write plain sentences instead. AI
writing is verbose by default and adds noise that hides the point, so trim it
before it leaves the machine.

## Comments

Omit redundant comments. Where a comment is needed, keep it to two lines.

Where a comment records a specific requirement, cite the source, including the
section number, page number or URL.

## No coding tool attributions

Never add `Co-Authored-By`, a "Generated with" footer, or any other coding tool
attribution to a commit, pull request, issue or comment.

Some projects require AI assistance to be disclosed. They ask for free-form
prose in the pull request body, which is what to write. Still never a trailer.

## Verify before you submit

Run the change. Do not submit code that is only plausibly correct.

Never write code for a platform or environment you cannot test on. Where
something was not verified, say so plainly rather than implying it was.

This repository is a cryptocurrency node. Consensus, cipher and wallet code
decides whether coins move correctly, so a change there is verified against a
running node and the integration tests, not against reasoning alone.

## Commits

The subject line is `package: lowercase summary` with no trailing period. Use a
comma separated list for several packages, and `all:` for a repo wide change.

Any behavior change gets a body, wrapped at about 72 columns, giving the
symptom, the cause and the fix. `Fixes #N` closes an issue and `For #N`
references one without closing it.

Organize a change into small commits that each build and pass the tests. A
commit whose only purpose is to fix up an earlier one in the same branch does
not belong in the final history.

## Pull requests

Avoid force pushing while a review is in progress, because it breaks the
conversation. Push fixup commits instead, then rebase once and force push a
single time when the review is done.

Load an issue's comments before acting on it, not just its description.

Branch from `develop`, which is the default branch. `master` carries releases.

## Working in this repository

`make check` runs the linters, the tests under both `amd64` and `386`, the
stable integration tests and the newcoin template check. `make test` alone runs
the unit tests.

`make format` formats with goimports, and `make lint` needs
`make install-linters` first.

Dependencies live in the `vendor` directory and the go command defaults to
`-mod=vendor` here, so a dependency change means running `go mod vendor` and
committing the result. `DEVELOPMENT.md` has the rules about which packages may
take a dependency at all, since `cipher` deliberately keeps none.

`make newcoin COIN=foo` regenerates `cmd/$COIN/$COIN.go` from the template. Do
not hand edit a generated file, change the template instead.
