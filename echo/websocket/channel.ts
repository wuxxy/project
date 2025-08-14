import {QUEUE} from "./main.ts";

export class Channel {
    public channelID: string;
    public sockets: Set<string>; // Set of socket IDs
    constructor(channelID: string) {
        this.channelID = channelID;
        this.sockets = new Set<string>(); // Initialize an empty set for sockets
    }
    public broadcast(event: string, data: any) {
    }
}