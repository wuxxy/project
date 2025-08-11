import {connect} from "@nats-io/transport-deno";
import { NatsConnection } from "@nats-io/nats-core";

export let nc: NatsConnection;
export async function connectNats() {
    if (nc) {
        return nc;
    }
    try {
        nc = await connect({
            servers: Deno.env.get("NATS_URL") || "nats://localhost:4222",
            name: "api-websocket",
        });
        console.log("Connected to NATS");
        return nc;
    } catch (error) {
        console.error("Failed to connect to NATS:", error);
        throw error;
    }
}