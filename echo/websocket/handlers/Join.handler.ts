import {Handler, SOCKETS} from "../main.ts";
import {Socket} from "../socket.ts";

export const handler: Handler = (data: any, socket: Socket) => {
    if(!socket.isReady){
        socket.error("Socket not ready");
        socket.instance.close();
        SOCKETS.delete(socket.id);
        return;
    }
    socket.lastPing = Date.now();
}