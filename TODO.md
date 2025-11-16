## Project: String Visualiser with Web UI (Go)

### Stage 1: Core logic module
- [x] Create Go module `internal/visualiser` and define `Result` struct
    - [x] Fields: character, codePointHex, codePointDec, utf8BytesHex, utf8BytesDec, utf8BytesBinary, htmlEntityDecimal, htmlEntityHex
- [x] Implement `AnalyseString(input string) ([]Result, error)`
- [x] Write unit tests covering:
    - [x] ASCII input
    - [x] Multi-byte input (e.g., “ई”, emoji)
    - [x] Empty input / error case
- [x] Verify output on console for example strings
- [x] Refactor logic for readability, add comments/documentation

### Stage 2: CLI tool interface
- [x] Create `cmd/visualiser/main.go`
- [x] Accept input string (via argument or prompt)
- [x] Call `AnalyseString` and print table to terminal
- [x] Format table with headings and aligned columns
- [x] Add `-reverse` flag to allow input of bytes or code points
    - [x] Implement reverse logic: parse bytes/code points → character(s) → `Result`
- [x] Test CLI tool manually and via simple script

### Stage 3: HTTP API backend
- [ ] Create `internal/web/handlers.go` and setup HTTP server
- [ ] `GET /` handler: serve placeholder page (or simple home)
- [ ] `POST /api/visualise` handler: accept JSON request → call `AnalyseString` (or reverse) → respond JSON
- [ ] Define JSON request/response schema
- [ ] Write integration test (e.g., using `httptest`): send request → verify response fields and values
- [ ] Run HTTP server locally (e.g., port 8080) and test with curl/Postman

### Stage 4: Front-end web UI
- [ ] Create HTML template `internal/web/templates/index.html`
    - [ ] Input field(s): for string and/or bytes/code points
    - [ ] Submit button
    - [ ] Section for displaying results table
- [ ] Add JS logic:
    - [ ] Capture form submission
    - [ ] Call `POST /api/visualise`
    - [ ] On response, build and insert results table dynamically
- [ ] Style UI with CSS (e.g., Bootstrap or Tailwind)
- [ ] Add validation/error-handling in the UI (e.g., empty input, invalid bytes)
- [ ] Test in browser: input string → see results table

### Stage 5: Download/Export and Enhancements
- [ ] Add UI buttons: “Download CSV”, “Download JSON”
    - [ ] Implement backend endpoint `GET /api/download?format=csv&input=...` or front-end conversion
- [ ] Add reverse-mode toggle in UI (to switch from string → bytes/code points)
- [ ] Add copy-to-clipboard buttons next to each metadata field in the results table
- [ ] Add theme toggle (dark mode/light mode)
- [ ] Optional: Integrate with TCP echo server
    - [ ] Allow TCP client to send string → server uses `AnalyseString` and logs/prints result
    - [ ] On web UI show recent TCP inputs and their visualisation
- [ ] Write integration tests end-to-end (CLI + API + UI)

### Stage 6: Deployment & Documentation
- [ ] Write `README.md`: overview, setup instructions, CLI usage, API usage, UI usage, tests
- [ ] Create `Dockerfile` and optionally `docker-compose.yml`
- [ ] Deploy to cloud / container platform (e.g., Render, Fly.io) and get live demo link
- [ ] Write blog post or project summary: architecture, key features, learnings
- [ ] Prepare portfolio entry: link to live demo + source repo

### Optional Future Enhancements
- [ ] Add WebSocket live-input mode: stream characters and visualise in real time
- [ ] Add caching layer so repeated inputs return faster
- [ ] Extend support to other encodings (UTF-16, UTF-32, etc)
- [ ] Add user accounts / save history of visualisations
- [ ] Internationalisation: UI translations, wide range of scripts  
