# Qdrant

A vector database for similarity search: store embeddings with payloads, then
query by vector with filtering. A common backend for semantic search and
retrieval-augmented generation (RAG).

## What gets created

- One application (`qdrant`, image `qdrant/qdrant:v1.20.0`). The REST API and
  web dashboard are on port 6333; gRPC stays on 6334 inside the workspace
  network.
- One volume, `storage` → `/qdrant/storage`, for collections, indexes and
  snapshots

## Inputs

| Input | Notes |
|---|---|
| **API key** | Clients must send it in the `api-key` header. Leave blank to auto-generate. It is stored as a secret env var. |

## After install

Other apps in the workspace reach Qdrant at its in-network alias on port 6333
(REST) or 6334 (gRPC), with the API key. If you attach a domain, the built-in
dashboard is at `/dashboard`; it asks for the same key.

## Data and backups

Everything lives in the `storage` volume. For a consistent copy of a running
instance, create a snapshot through the API and back up the volume afterwards.

## Upgrades

Qdrant supports upgrading **one minor version at a time** (1.19 → 1.20). Step
through intermediate versions instead of skipping them.
