import {Socket} from "../socket.ts";
import {Handler, HEARTBEAT_INTERVAL, SOCKETS} from "../main.ts";

export const handler: Handler = (data: any, socket: Socket) => {
    if(!socket.isReady){
        socket.error("Socket not ready");
        socket.instance.close();
        SOCKETS.delete(socket.id);
        return;
    }
    if (Date.now() - socket.lastPing > HEARTBEAT_INTERVAL+500 || Date.now() - socket.lastPing < HEARTBEAT_INTERVAL+500) {
        socket.error("Incorrect ping interval");
        socket.instance.close();
        SOCKETS.delete(socket.id);
        return;
    }
}