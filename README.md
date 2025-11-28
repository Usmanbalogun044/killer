# Killer Organisation - Enterprise Stress Tool

> [!CAUTION] > **DISCLAIMER: THIS TOOL IS FOR EDUCATIONAL PURPOSES ONLY.**
>
> This software is intended to demonstrate how Denial of Service (DoS) attacks work conceptually and to help understand network security and stress testing.
>
> **DO NOT** use this tool against any system, network, or website that you do not own or have explicit permission to test. Unauthorized use of this tool is illegal and unethical. The author accepts no responsibility for any misuse of this code.

## Overview

This is an **Enterprise-Grade Stress Tool** developed by **Killer Organisation**. It is designed for professional load testing and security auditing, providing detailed insights into server performance under high stress.

**Key Features:**

- **Enterprise Reporting**: Tracks detailed HTTP status codes (2xx, 3xx, 4xx, 5xx) to identify exactly how the server is failing.
- **Test Duration Control**: Set precise test durations (e.g., `30s`, `5m`) for reproducible benchmarks.
- **HTTP/2 Support**: Leverages HTTP/2 multiplexing for maximum throughput.
- **Advanced IP Rotation**: Supports HTTP/SOCKS proxies to rotate IP addresses for every request.
- **Real-time Monitoring**: Displays live Requests Per Second (RPS), Success counts, and Error counts.

## Prerequisites

- **Go (Golang)**: You need to have Go installed on your machine. [Download Go](https://go.dev/dl/).

## Installation

1.  **Clone the repository**:

    ```bash
    git clone https://github.com/Usmanbalogun044/killer.git
    cd killer
    ```

2.  **Build the program**:
    ```bash
    go build -o killer dosattack.go
    ```

## Usage

### Command Line Flags

- `-t`: **Target URL**. The full URL of the target (e.g., `https://example.com`).
- `-th`: **Threads**. Number of concurrent workers. Default: `100`.
- `-d`: **Duration**. Length of the test (e.g., `30s`, `1m`, `1h`). Default: `0` (Infinite).
- `-proxy`: **Proxy List**. Path to a file containing proxies (one per line).

### Examples

**1. Quick 30-Second Load Test**
Target `https://example.com` with 100 threads for 30 seconds:

```bash
./killer -t https://example.com -th 100 -d 30s
```

**2. Sustained Enterprise Stress Test**
Target `https://target-site.com` with 500 threads for 1 hour, using proxies:

```bash
./killer -t https://target-site.com -th 500 -d 1h -proxy proxies.txt
```

## Reporting & Analysis

At the end of the test, the tool generates a **Test Completion Report**:

```text
=================================================
             TEST COMPLETION REPORT
=================================================
 Target:          https://example.com
 Duration:        30.00s
 Total Requests:  15420
 Avg RPS:         514.00
-------------------------------------------------
 [2xx] Success:   15000
 [3xx] Redirects: 0
 [4xx] Client Err:420
 [5xx] Server Err:0
 [Err] Net Errors:0
=================================================
```

### Interpreting Results for cPanel / Shared Hosting

If you are testing your own cPanel or shared hosting environment, look for these specific signs:

- **508 Resource Limit Is Reached**: This is a common **5xx** error on cPanel. It means your test successfully hit the CPU/RAM limits assigned to your account (CloudLinux LVE limits).
- **503 Service Unavailable**: The web server (Apache/LiteSpeed) is overloaded or restarting.
- **403 Forbidden**: If you see many **4xx** errors, the server's firewall (like ModSecurity or Imunify360) has likely detected the attack and blocked your IP.

## License

This project is open-source and available for educational study.
