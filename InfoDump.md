# GeneSyS

Generic Sync System

A single binary that can act as both client and server. Any directory becomes a sync-enabled repo once initialized. Both sides are identical in structure and capability. Syncing is metadata driven, with QUIC as the transport layer.

---

## Core Model

### Repo Structure

Every repo contains:

```
.genesys/
    genesys.db
```

Both sides use identical layout. No “local vs remote” differences. No `.syncmeta`.

### Identity

A device is “local” only because you ran the command there.
A “remote” is just another instance running `genesys serve`.

That’s the only distinction.

---

## Metadata Schema

**files** table:

* uuid (text, primary key)
* name (text)
* relative_path (text)
* file_hash (text, nullable)
* last_modified (timestamp)

### Hashing

Hashing is optional, used only when timestamps are insufficient or ambiguous. The system never tries to sync timestamps across machines.

### Timestamps are local only

Each repo treats its filesystem’s timestamps as **local ground truth** for change detection. The DB always reflects local reality before any sync.

---

## CLI

```
genesys init <folder>
genesys serve --port <n>
genesys remote add <name> <quic://ip:port>
genesys sync <remote> --mode=<merge | local-wins | remote-wins>
```

One binary does everything.

---

## Initialization

`genesys init <folder>`:

1. Create folder if missing.
2. Create `.genesys/`.
3. Create `genesys.db` with empty schema.
4. No auto-scan unless `--scan-existing` is used.

Identical behavior for all devices.

---

## Local Maintenance Model

Every time you run **any** GeneSyS command (init, sync, serve, whatever), the tool performs a local refresh:

1. Scan the repo folder.
2. For each file on disk:

   * Read OS `mtime`.
   * If size or mtime changed, recompute hash.
   * Update DB.
3. For each DB entry not found on disk, remove it.

DB always reflects the filesystem, Git-style. No timestamp synchronization across systems. No assumptions about remote timestamps ever.

---

## Remotes

`genesys remote add origin quic://192.168.1.20:9000`

Stores that endpoint for use during sync. Multiple remotes allowed.

---

## Transport Layer: QUIC

All operations (metadata, file transfer, diffs) run over QUIC.

Each file transfer uses a dedicated QUIC stream. Transfers are parallel by default.

Examples:

* Metadata fetch: one stream.
* Pull 50 files: 50 parallel streams.
* Push files: same.

---

## Sync Flow

### 1. Local refresh

Local performs maintenance scan, updates DB.

### 2. Local requests remote metadata

Via QUIC `/meta`.

Remote refreshes its own DB before replying.

### 3. Diffing rules (uuid-based)

For each uuid:

A. Exists only local
B. Exists only remote
C. Exists both, timestamps differ
D. Exists both, timestamps equal (hash optional check)
E. Exists both, name/path changed (rename)

Timestamps only apply to detect local state changes, never as authoritative between machines.

### 4. Conflict Resolution

`--mode=` controls behavior:

#### merge

Two way:

* Newer timestamp wins.
* Missing files propagate.
* Rename detected by uuid.
* Hash used only when timestamps disagree or edge cases appear.

#### local-wins

Local overrides remote.
Remote-only files can be ignored or deleted.

#### remote-wins

Remote overrides local.
Local-only files ignored or deleted.

---

## QUIC File Transfer

### Pull

Local:
`GET_FILE { uuid }`

Remote opens stream and sends raw bytes.

### Push

Local:
`PUT_FILE { uuid, size, metadata }`

Then streams bytes.
Remote writes file, updates DB.

Parallel streams default.

---

## Rename Handling

Rename is trivial because uuid is identity.

If uuid same, but name or relative path differ, it’s a rename.
DB updates accordingly.
Hash only used if needed.

---

## Serving

`genesys serve --port 9000`

Exposes minimal QUIC endpoints:

* `/meta`
* `/file/get?uuid=`
* `/file/put?uuid=`

Server does no sync on its own.
It only responds.
Any device can serve, even low powered, using wakelock if needed.

---

## Lightweight Philosophy

* Never full-file hashing unless required.
* Timestamps are local detectors, not synchronized artifacts.
* Metadata-first always.
* Deterministic. Symmetric. Zero hidden state.
* Git-style metadata tracking, rsync-style data transfer, QUIC transport.

GeneSyS is predictable, transparent, and fully two way without platform special casing.

---