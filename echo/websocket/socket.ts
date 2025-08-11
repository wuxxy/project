import {SOCKETS, STARTING_CREDIT} from "./main.ts";
import {ChannelID, EVENT} from "./types.ts";
export class Socket {
    instance: WebSocket;
    public isReady: boolean = false; // Indicates if the socket is ready to send messages
    public rooms: Set<ChannelID>;
    public credit: number;
    public lastPing: number;
    public id: string; // Unique identifier for the socket
    public userID: string; // User ID associated with the socket
    constructor(instance: WebSocket) {
        this.instance = instance;
        this.rooms = new Set();
        this.credit = STARTING_CREDIT;
        this.lastPing = Date.now();
        this.id = crypto.randomUUID(); // Generate a unique ID for the socket
        this.userID = ""; // Initialize userID as an empty string
        SOCKETS.set(this.id, this);
    }
    public send(event: EVENT, data: any) {
        if (this.instance.readyState === WebSocket.OPEN) {

        }
    }
    public sendDirectly(event: EVENT, data: any) {
        this.instance.send(JSON.stringify({
            e: event,
            d: data
        }));
    }

    public error(error: any) {
        this.sendDirectly("error", {message: error})
    }

}