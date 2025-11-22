# GeneSyS

Generic Sync System

GeneSyS is a lightweight, file-level sync tool that keeps two folders identical. It works over QUIC for fast, stable transfers and uses a small `.genesys` directory to track file metadata. Both sides run the same binary. One side acts as the **local** executor, the other as the **remote** server.

GeneSyS keeps decisions simple: add, overwrite, delete, or skip. No cloud, no accounts, no hidden magic.

---
## Features

GeneSyS provides a predictable sync model:

**Single binary**
One executable handles both modes, `serve` and `sync`.

**QUIC streaming**
Fast, low-latency transfers with modern congestion control.

**Minimal metadata**
A tiny SQLite DB inside `.genesys` stores UUIDs, names, sizes, timestamps, and per-file hashes.

**Deterministic diffing**
Each run builds a local snapshot, compares it with the remote snapshot, and produces a clear action set.

**Two-way aware, one-way controlled**
Sync direction is explicit. Only the initiating side performs the final decision and reconciliation.

**Platform-neutral**
Runs anywhere Windows or Termux can handle a filesystem and SQLite.

---

## How It Works

1. **Each folder has a `.genesys` directory**
   This stores the internal database and metadata.

2. **Local side runs a scan**
   It updates metadata for new, removed, or modified files before any sync occurs.

3. **Remote exposes its snapshot**
   Running `genesys serve` makes the folder's metadata and files available over QUIC.

4. **Local pulls the snapshot**
   It compares UUID sets and metadata to determine exact actions.

5. **Local sends back an action plan**
   The remote applies deletes, renames, creates, and overwrites exactly as instructed.

6. **Remote recalculates metadata**
   After changes, the remote updates its local snapshot to match reality.

7. **Both sides end fully aligned**
   UUIDs, names, sizes, timestamps, and content hashes match.

---

## Sync Logic

### When comparing snapshots:

* UUID in both, metadata identical:
  skip

* UUID in both, metadata differs:
  overwrite

* UUID only in local:
  create on remote

* UUID only on remote:
  delete on remote

* Filename changed, UUID matches:
  rename

This keeps behavior stable, predictable, and safe.

---

## Commands (conceptual)

```
genesys init <folder>
    Initializes .genesys and builds the metadata DB.

genesys scan <folder>
    Refreshes metadata without contacting any remote.

genesys serve <folder> --port 4433
    Exposes the folder metadata and file streams over QUIC.

genesys sync <folder> <remote-address>
    Pulls remote snapshot, computes diff, sends actions, applies results.
```

---

## Ideal Use Cases

Personal two-device sync.
Phone to laptop.
Laptop to server.
Photos. Documents. Music. Code.
Anything where you want deterministic, offline-friendly folder mirroring without cloud dependencies.

---

## Status

Under active development.
Design is stable.
Implementation work is iterative.

---