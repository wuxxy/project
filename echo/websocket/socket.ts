import {CHANNELS, SOCKETS, STARTING_CREDIT} from "./main.ts";
import {ChannelID, EVENT} from "./types.ts";
export class Socket {
    instance: WebSocket;
    public isReady: boolean = false; // Indicates if the socket is ready to send messages
    public credit: number;
    public lastPing: number;
    public id: string; // Unique identifier for the socket
    public userID: string; // User ID associated with the socket
    constructor(instance: WebSocket) {
        this.instance = instance;
        this.credit = STARTING_CREDIT;
        this.lastPing = Date.now();
        this.id = crypto.randomUUID(); // Generate a unique ID for the socket
        this.userID = ""; // Initialize userID as an empty string
        SOCKETS.set(this.id, this);
        this.instance.onclose = () => {
            SOCKETS.delete(this.id);
        }
    }
    public send(event: EVENT, data: any) {

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

    public subscribe(channelID: ChannelID) {
        let set = CHANNELS.get(channelID);
        if (!set) {
            set = new Set<Socket>();
            CHANNELS.set(channelID, set);
        }
        set.add(this);
    }
    public unsubscribe(channelID: ChannelID) {
        const set = CHANNELS.get(channelID);
        if (set) {
            set.delete(this);
            if (set.size === 0) {
                CHANNELS.delete(channelID);
            }
        }
    }

}