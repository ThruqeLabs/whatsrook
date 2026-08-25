# AGENTS Architecture & Guidelines

Welcome to WhatsRook.

## Core Architecture

- **Root Library (`whatsrook`)**: High-level Go abstraction cleanly wrapping `wa-core` (`Client`, session initialization, phone number pairing, QR channel streaming, session wipe helpers, and client lifecycle management).
- **Package Layout**:
  - `whatsrook`: Core library entrypoint and high-level client abstraction.
  - `whatsrook/logger`: Universal, high-performance Zap-based logger (`package Logger`) with multi-destination routing, per-level log files, and zero-allocation structured fields.
  - `whatsrook/wa-core`: WhatsApp protocol core library and VoIP call engine (`wacaller`), binary XML node encoding/decoding, socket engine, SRTP/SRTCP media pipeline, STUN NAT traversal, audio/video playout controllers, and database stores (`wa-core/store/sqlstore`).
  - `whatsrook/utils`: Core protocol & messaging abstractions over `wa-core`, with modularized domain subpackages:
    - `whatsrook/utils/media`: Audio/video transcoding engine (Opus OGG with waveform calculation, JPEG, MP3, Annex-B H.264 stream splitting) and WebP sticker EXIF metadata manipulation.
    - `whatsrook/utils/system`: Runtime OS, CPU, RAM metrics, uptime calculation, and diagnostic helpers.
    - `whatsrook/utils/formatting`: WhatsApp text formatting, sanitization, emoji removal, and language validation.
    - `whatsrook/utils/httputil`: HTTP client utility helpers with timeouts and header management.
    - `whatsrook/utils/cache`: Multi-backend caching abstractions (in-memory, Redis).
    - `whatsrook/utils/qr`: QR rendering and streaming HTTP server.
  - `whatsrook/cli`: WhatsApp bot CLI application, plugin commands (`cli/plugins`), dedicated TUI package (`cli/tui`: interactive Bubbletea setup wizard and live agentic dashboard), and modular domain subpackages under `whatsrook/cli/utils`:
    - `whatsrook/cli/utils/font`: Unicode font styling, styling tables, and text normalization.
    - `whatsrook/cli/utils/timezone`: Windows & IANA timezone mapping database and alias resolver.
    - `whatsrook/cli/utils/markets`: Real-time market quotes and Bitcoin network / halving predictors.
    - `whatsrook/cli/utils/news`: AP News & WABetaInfo scrapers, article parsers, and article selection.
    - `whatsrook/cli/utils/meta`: Meta AI client, protocol response streamer, prompt builder, and personality traits.
    - `whatsrook/cli/utils/games`: Game engines (Tic-Tac-Toe, Word Chain Game, Unscramble word game).
    - `whatsrook/cli/utils/fun`: Content and entertainment utilities (facts, quotes, jokes, rizz, math solver, why.com search).
    - `whatsrook/cli/utils/system`: Process diagnostics, CPU / RAM / uptime metric helpers.
    - `whatsrook/cli/utils/media`: Message history and quoted context tracking.

## Development Management

Utilize the [Taskfile](./Taskfile.yml)

## Relevant Documentation

- [Docs](./docs)
- [Security](./SECURITY.md)
