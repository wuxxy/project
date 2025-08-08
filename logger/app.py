import asyncio
import json
import nats
from nats.errors import ConnectionClosedError, TimeoutError, NoServersError

async def main():
    try:
        nc = await nats.connect("nats://localhost:4222")

        # Subscribe to the "logs" subject
        async def log_handler(msg):
            subject = msg.subject
            data = msg.data.decode()
            try:
                log = json.loads(data)
            except:
                log = data
            print(f"[{subject}] {log}")

        await nc.subscribe("logs.>", cb=log_handler)
        print("Listening for logs...")

        # Keep running
        while True:
            await asyncio.sleep(1)

    except (ConnectionClosedError, TimeoutError, NoServersError) as e:
        print(f"Error while connecting: {e}")

if __name__ == '__main__':
    asyncio.run(main())
