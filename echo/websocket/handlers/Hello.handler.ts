import {Socket} from "../socket.ts";
import {nc} from "../ipc.ts";
import { encode, decode } from "jsr:@std/msgpack";
import {HEARTBEAT_INTERVAL} from "../main.ts";
type VerifyTokenResponse = {
    error: string;
    session_id: string;
    user_id: string;
}
export const handler = (data: any, socket: Socket) => {
    if (Date.now() - socket.lastPing > 5000) {
        socket.error("Ping timeout");
        socket.instance.close();
        return;
    }
    if (socket.userID){
        socket.error("Already authenticated");
        return;
    }
    try{
        nc.request("auth.verify_token", data.t).then((msg) => {
            const decodedRes = decode(msg.data) as VerifyTokenResponse;
            if (decodedRes.error && !decodedRes.session_id && !decodedRes.user_id) {
                socket.error(decodedRes.error);
                return;
            }
            socket.userID = decodedRes.user_id;
            socket.lastPing = Date.now();
            socket.sendDirectly("hello", {
                ssid: socket.id,
                user_id: decodedRes.user_id,
                heartbeat: HEARTBEAT_INTERVAL,
            });
            socket.isReady = true;
        })
    }catch(err){
        console.error(err)
    }
}