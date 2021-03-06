# sntr: all of Sentry at your fingertips

The `sntr` command-line program gives you convenient access to
[Sentry](https://sentry.io) directly from your terminal.

---

**Disclaimer: this is not an official Sentry project.**

Despite being under the `getsentry` GitHub organization, this project is only an
experiment developed during Sentry's internal Hack Week 2020.

Expect no support, expect things to be broken in unimaginable ways, use at your
own risk.

---

Features:

- [x] List organizations: `sntr get organizations`
- [x] List projects: `sntr get projects`
- [x] List project issues: `sntr get ORG_SLUG/PROJ_SLUG`
- [x] Search issues: `sntr get --query`
- [ ] Get events in JSON format
- [x] Send [test events](#test-events): `sntr send`
- [x] Create a [proxy](#proxy) between your program and the Sentry ingestion API
- [ ] Traces and transactions
- [ ] List releases
- [ ] Multi DSN -- read/write multiple projects at once
- ...

## Test events

`sntr` can send test events and time how long it takes to get it back after it
has been processed by Sentry.

```
sntr send
```

NOTE: [`sentry-cli` can send events](https://docs.sentry.io/cli/send-event/)
too and with more advanced options.

## Proxy

In proxy mode, `sntr` gives you the power to intercept all requests from a
Sentry SDK to Sentry's ingestion API.

```
sntr proxy
```

Features:

- [ ] Record any kind of outgoing data from SDKs before it goes to Sentry
- [ ] Interactively modify data before it is sent
- [ ] Programmatically modify data before it goes out

## Exec

The `exec` subcommand turns `sntr` into a wrapper that can execute and capture
errors and crashes from arbitrary processes.

```
sntr exec <cmd>
```

Features:

- [ ] If `<cmd>` exits with non-zero exit code, an error is sent to Sentry
- [ ] If `sntr-exec` receives a termination signal, it tries to send an event to
      Sentry.

## Extra

- [ ] Record latencies of every HTTP interaction
- [ ] Visualize latency distribution with "HdrHistogram-like" graph

## How to install

- [ ] `brew install sntr`
- [ ] Download [release](https://github.com/getsentry/sntr/releases) from GitHub
