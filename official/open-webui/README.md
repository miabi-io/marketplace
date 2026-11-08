# Open WebUI

A self-hosted chat interface for large language models. Connect it to an Ollama
server, to OpenAI, or to any OpenAI-compatible API, and chat with models, build
custom assistants and ask questions about your own documents (RAG).

## What gets created

- One application (`open-webui`, image `ghcr.io/open-webui/open-webui:v0.11.4`)
  on port 8080
- One PostgreSQL 17 database for users, chats and settings
- One volume, `data` → `/app/backend/data`, for uploads, the vector store and
  downloaded embedding models

This template does not run a model itself. Point it at an Ollama server you run
elsewhere, or at a hosted API.

## Inputs

| Input | Notes |
|---|---|
| **Secret key** | Signs sessions. Leave blank to auto-generate. |
| **Ollama URL** | e.g. `http://ollama:11434`. Leave blank if you only use hosted APIs. |
| **OpenAI-compatible API URL** | Defaults to `https://api.openai.com/v1`. Change it to use another provider. |
| **API key** | Key for that endpoint. You can also add connections later in the admin settings. |

## After install

**The first start is slow:** Open WebUI downloads its embedding model before it
starts listening, so give it a few minutes. The template disables the
healthcheck for that reason. Watch the application logs.

Then open the app. The first account to sign up becomes the **admin**; later
sign-ups wait for an admin's approval. Manage models and connections under
**Admin Panel → Settings → Connections**.

## Data and backups

Back up the database (accounts, chats) and the `data` volume (documents and
their embeddings).
