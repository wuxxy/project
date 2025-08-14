import {writable} from "svelte/store";

export class WSCLient {
    private ws?: WebSocket;
    private url: string;
    private messageHandler = new Map<string, (data: any) => void>();
    private heartbeatInterval?: number;
    public connected = writable<boolean>(false);
    constructor(url: string) {
        this.url = url;
        this.sendHeartbeat = this.sendHeartbeat.bind(this);
    }

    public connect(){
        if(this.ws){
            return;
        }
        this.ws = new WebSocket(this.url);
        this.ws.onopen = (event) => {
            console.log("Connected to WebSocket server:", this.url);
            // Authenticate the user by sending a token using hello event
            this.send("hello", {t: `Bearer ${localStorage.getItem("token")}`});
            // Set up heartbeat to keep the connection alive
            this.on("hello", (data) => {
                this.heartbeatInterval = setInterval(this.sendHeartbeat, data.heartbeat);
                this.connected.set(true);
            })

        };
        this.ws.onclose = (event) => {
            clearInterval(this.heartbeatInterval);
            this.connected.set(false);
        }
        // Handle messages
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            if (message.e && message.d) {
                this.messageHandler.get(message.e)?.(message.d);
            } else {
                console.error("Invalid message format:", message);
            }
        };
    }

    public send(event: string, data: any) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            console.error("WebSocket is not open. Cannot send message.");
            return;
        }
        this.ws.send(JSON.stringify({
            e: event,
            d: data
        }))
    }

    public on(event: string, callback: (data: any) => void) {
        this.messageHandler.set(event, callback);
    }

    private sendHeartbeat() {
        this.send("ping", {})
    }
    public getConnected () {
        return this.connected;
    }

}