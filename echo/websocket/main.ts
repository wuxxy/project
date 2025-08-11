import {Socket} from "./socket.ts";
import {ChannelID, EVENT, MESSAGE, MINUTE, SECOND, SocketID} from "./types.ts";
import {connectNats} from "./ipc.ts";


// A queue for messages that are not yet sent to channels
let QUEUE: Map<ChannelID, MESSAGE>
// A list of all active sockets
export const SOCKETS: Map<SocketID, Socket> = new Map();
// A list of all active channels
export const CHANNELS: Map<ChannelID, Set<Socket>> = new Map();


// Internal event handler registration
export type Handler = (data: any, socket: Socket) => void;
const HANDLERS: Map<EVENT, Handler> = new Map();
const handlerDir = Deno.readDirSync("./handlers")
for (const file of handlerDir.filter(file => file.name.includes(".handler"))) {
    import("./handlers/" + file.name).then(s => {
        if(s.handler){
            HANDLERS.set(file.name.split(".")[0].toLocaleLowerCase() as EVENT, s.handler);
        }
    })
}
// Connect to NATS server - IPC
connectNats()

// The amount of credit a socket starts with
// This is used to limit the amount of messages a socket can send
export const STARTING_CREDIT = 500;
// Heartbeat handling
export const HEARTBEAT_INTERVAL = 3000;
// Makes sure every socket is alive
setInterval(() => {
    SOCKETS.forEach((socket: Socket) => {
        if(!socket.isReady) {
            if(Date.now() - socket.lastPing > 15*SECOND) {
                socket.error("Socket not ready");
                socket.instance.close();
                SOCKETS.delete(socket.id);
                return;
            }
            return;
        };
        if (Date.now() - socket.lastPing > HEARTBEAT_INTERVAL+500 || Date.now() - socket.lastPing < HEARTBEAT_INTERVAL+500) {
            socket.error("Ping timeout");
            socket.instance.close();
            SOCKETS.delete(socket.id);
            return;
        }
    })
}, HEARTBEAT_INTERVAL);
// Start serving WebSocket connections
Deno.serve((req) => {
    console.log(HANDLERS);
    if (req.headers.get("upgrade") != "websocket") {
        return new Response(null, { status: 426 });
    }

    const { socket: wsConnection, response } = Deno.upgradeWebSocket(req);
    let socket: Socket | undefined;
    wsConnection.addEventListener("open", () => {
        socket = new Socket(wsConnection);

    });

    wsConnection.addEventListener("message", (event: MessageEvent) => {
        try {
            const incoming_message: MESSAGE = JSON.parse(event.data);
            if (incoming_message.d == null || incoming_message.e == null) {
                socket?.error("Invalid format")
                return;
            }
            HANDLERS.get(incoming_message.e)?.(incoming_message.d, socket!);
        }catch (error) {
            if(error instanceof SyntaxError ) {
                socket?.error("Invalid format")
            }
        }
    });

    return response;
});