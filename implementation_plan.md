# Implementation Plan - Enterprise Reporting & Metrics

To make the tool suitable for "enterprise testing" (like testing Microsoft servers), we need more than just raw load. We need **observability**. Enterprise testing requires knowing _how_ the server failed, not just that it received traffic.

## Proposed Changes

### 1. Detailed Metrics Collection

Instead of just counting `ops` (requests), we will track:

- **Status Codes**: Count 2xx (Success), 4xx (Client Error/Rate Limited), 5xx (Server Error/Down), and Connection Errors.
- **Latency**: Measure the average response time to understand if the server is slowing down under load.

### 2. Test Duration Control

- Add a `-duration` flag (e.g., `30s`, `5m`) to run the test for a specific time and then stop automatically.
- This allows for precise, reproducible test scenarios.

### 3. Final Report

- At the end of the test, print a summary report:
  - Total Requests
  - Success Rate (%)
  - Status Code Breakdown
  - Average RPS

## File Changes

### [dosattack.go](file:///c:/Users/dollarhunter/Documents/github/Killer/dosattack.go)

- **Structs**: Create a `Metrics` struct with atomic counters for `Requests`, `Success`, `Fail`, `Status2xx`, `Status4xx`, `Status5xx`.
- **Attack Function**:
  - Update to increment specific counters based on `resp.StatusCode`.
  - Measure `time.Since(start)` for each request to track latency (optional, might be too heavy for high-perf, maybe sample it).
- **Main Function**:
  - Add `-duration` flag.
  - Implement a timer to stop the attack.
  - Print the "Enterprise Report" on exit.

## Verification Plan

- **Run Test**: Run against a target for 10 seconds.
- **Check Output**: Verify the live stats show the breakdown and the final report appears.
