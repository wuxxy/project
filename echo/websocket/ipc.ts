import {connect} from "@nats-io/transport-deno";
import { NatsConnection } from "@nats-io/nats-core";
function natsOptsFromEnv() {
    const raw = (Deno.env.get("NATS_URL") ?? "").trim();
    // fallback if empty
    if (!raw) return { servers: "nats://localhost:4222" } as const;

    // strip accidental wrapping [] (from bad merges)
    const cleaned = raw.replace(/^\[(.*)\]$/, "$1");

    try {
        const u = new URL(cleaned);
        const servers = `${u.protocol}//${u.hostname}${u.port ? ":" + u.port : ""}`;
        // normalize/encode password once
        const user = u.username || Deno.env.get("NATS_USER") || undefined;
        const pass = u.password
            ? encodeURIComponent(decodeURIComponent(u.password))
            : Deno.env.get("NATS_PASS") || undefined;

        return { servers, user, pass };
    } catch {
        // If NATS_URL is just "host:port" or similar, use as-is.
        return { servers: cleaned };
    }
}

const opts = natsOptsFromEnv();
export let nc: NatsConnection;
export async function connectNats() {
    if (nc) {
        return nc;
    }
    try {
        console.log(Deno.env.get("NATS_URL"))
        nc = await connect({
            ...opts,
            name: "api-websocket",
        });
        console.log("Connected to NATS");
        return nc;
    } catch (error) {
        console.error("Failed to connect to NATS:", error);
        throw error;
    }
}
