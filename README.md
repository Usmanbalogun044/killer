# ☠️ KILLER ☠️

```
  _  __  ___   _       _       ______   _____
 | |/ / |_ _| | |     | |     |  ____| |  __ \
 | ' /   | |  | |     | |     | |__    | |__) |
 |  <    | |  | |     | |     |  __|   |  _  /
 | . \  _| |_ | |____ | |____ | |____  | | \ \
 |_|\_\|_____||______||______||______| |_|  \_\

      >>> HIGH-PERFORMANCE HTTPS STRESSER <<<
           >>> CODED BY DOLLARHUNTER <<<
```

**Killer** is a lethal, high-performance HTTPS stress testing tool. It floods targets with a tsunami of requests, testing the absolute limits of their resilience.

> [!CAUTION] > **⚠️ WARNING: AUTHORIZED USE ONLY ⚠️**
> This tool is a weapon for stress testing. **DO NOT** use this on targets you do not own or have explicit permission to test. The author, **dollarhunter**, takes **NO RESPONSIBILITY** for the destruction you cause. You have been warned.

## ⚡ Capabilities

- **🔐 HTTPS Bypass**: Native TLS injection to penetrate SSL layers.
- **🚀 Hyper-Threading**: Spawns thousands of concurrent goroutines for maximum impact.
- **🎭 Ghost Mode**: Rotates User-Agents, Referers, and Headers to evade detection and mimic legitimate traffic.
- **🛡️ Proxy Support**: Automatically rotates IPs from `proxies.txt` to bypass rate limits and firewalls.
- **📊 Live Intel**: Real-time RPS (Requests Per Second) monitoring and status reports.

## 🛠️ Deployment

**Prerequisites:** [Go (Golang)](https://go.dev/dl/)

1.  **Infiltrate the Repository:**

    ```bash
    git clone https://github.com/Usmanbalogun044/killer.git
    cd killer
    ```

2.  **Execute Payload:**

    ```bash
    go run killer.go -t <TARGET_IP>
    ```

3.  **Proxy Configuration (Optional but Recommended):**
    Create a file named `proxies.txt` in the same directory.
    Format: `IP:Port` or `User:Pass@IP:Port`
    ```text
    127.0.0.1:8080
    user:pass@1.2.3.4:3128
    ```
    _Killer will automatically load these and rotate your IP to evade bans._

## 💀 Usage

Launch the attack from your terminal:

```bash
go run killer.go -t <TARGET> -p <PORT> -th <THREADS>
```

### 🚩 Command Flags

| Flag  | Description                                 | Default       |
| :---- | :------------------------------------------ | :------------ |
| `-t`  | **Target Host** (Domain/IP). No `https://`. | `example.com` |
| `-p`  | **Target Port**. Usually 443 for SSL.       | `443`         |
| `-th` | **Thread Count**. The power of the swarm.   | `100`         |

### 💣 Attack Scenarios

**Standard Strike:**

```bash
go run killer.go -t example.com
```

**Heavy Artillery (Custom Port + 500 Threads):**

```bash
go run killer.go -t test-server.com -p 8443 -th 500
```

## 👨‍💻 Operator

**Code by [dollarhunter]**
_Stay low, move fast._
