import {derived, writable} from 'svelte/store';
import {WSCLient} from "$lib/client/ws";

// Writable store to hold the socket instance
export const socket = writable<WSCLient | null>(null);
// This subscribes to the *inner* connected store and swaps when socket changes
export const socketConnected = derived(socket, ($socket, set) => {
    if (!$socket?.connected) { set(false); return; }
    const unsub = $socket.connected.subscribe(set);
    return () => unsub();
});
export function initSocket(url: string) {
    if (typeof window === 'undefined') return; // skip SSR

    const ws = new WSCLient(url);
    ws.connect();
    socket.set(ws);
    console.log(`WebSocket client initialized with URL: ${url}`);
}
