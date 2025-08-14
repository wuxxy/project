import {Configs, CpConfig, scanSync} from "./reader.ts";
import { load } from "https://deno.land/std/dotenv/mod.ts";
async function prompt(message: string = "") {
    const buf = new Uint8Array(1024);
    await Deno.stdout.write(new TextEncoder().encode(message + ": "));
    const n = <number>await Deno.stdin.read(buf);
    return new TextDecoder().decode(buf.subarray(0, n)).trim();
}
const workers = new Map<string, Worker>(); // keyed by v.name (or whatever you prefer)

// ChatGPT'd:
const C = {
    reset: "\x1b[0m",
    dim: "\x1b[2m",
    bold: "\x1b[1m",
    fg: {
        gray: "\x1b[90m",
        red: "\x1b[31m",
        green: "\x1b[32m",
        yellow: "\x1b[33m",
        blue: "\x1b[34m",
        magenta: "\x1b[35m",
        cyan: "\x1b[36m",
        white: "\x1b[37m",
    },
};
const PALETTE = [
    C.fg.cyan, C.fg.magenta, C.fg.yellow, C.fg.blue, C.fg.green, C.fg.red, C.fg.white, C.fg.gray,
];

function hash(s: string) { // deterministic color per name
    let h = 2166136261;
    for (let i = 0; i < s.length; i++) h = (h ^ s.charCodeAt(i)) * 16777619;
    return Math.abs(h) >>> 0;
}
function colorFor(name: string) { return PALETTE[hash(name) % PALETTE.length]; }
function timeStamp() {
    const d = new Date();
    const pad = (n: number) => String(n).padStart(2, "0");
    return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}
function padName(name: string, width = 14) {
    return name.length <= width ? name.padEnd(width, " ") : name.slice(0, width - 1) + "…";
}
function prefix(name: string) {
    const col = colorFor(name);
    return `${C.dim}${timeStamp()}${C.reset} ${col}[${padName(name)}]${C.reset}`;
}
function tint(kind: string, s: string) {
    // stderr/error lines get a subtle red tint; stdout stays normal
    const isErr = /err/i.test(kind);
    return isErr ? `${C.fg.red}${s}${C.reset}` : s;
}
function statusBadge(ok: boolean, code?: number) {
    const col = ok ? C.fg.green : C.fg.red;
    return `${col}${ok ? "OK" : "FAIL"}${code !== undefined ? ` (code ${code})` : ""}${C.reset}`;
}
export function runScriptPrefixed(path: string, v: CpConfig, script: string = "DEV") {
    const command = v.scripts?.[script];
    if (!command) throw new Error(`Script "${script}" not found for ${v.name}`);

    const worker = new Worker(new URL("./worker.ts", import.meta.url).href, {
        type: "module",
        deno: {
            permissions: { run: true, read: true, env: true },
        },
    });


    worker.addEventListener("message", (ev: MessageEvent) => {
        const msg = ev.data;
        if (msg.type === "line") {
            console.log(`${prefix(msg.name)} │ ${tint(msg.kind ?? "", msg.line)}`);
        } else if (msg.type === "exit") {
            console.log(`${prefix(msg.name)} ┆ ${statusBadge(msg.success, msg.code)}`);
            workers.delete(v.name);
        } else if (msg.type === "error") {
            console.error(`${prefix(msg.name)} │ ${C.fg.red}ERROR:${C.reset} ${msg.error}`);
        }
    });

    worker.postMessage({
        type: "start",
        payload: { path, name: v.name, command },
    });

    workers.set(v.name, worker);
    return worker;
}

export function stopScript(name: string) {
    const w = workers.get(name);
    if (!w) return Promise.resolve(false);
    return new Promise<boolean>((resolve) => {
        const handler = (ev: MessageEvent) => {
            const msg = ev.data;
            if (msg?.type === "exit" && msg.name === name) {
                w.removeEventListener("message", handler);
                w.terminate();
                workers.delete(name);
                resolve(true);
            }
        };
        w.addEventListener("message", handler);
        w.postMessage({ type: "stop" });
    });
}
function startProcesses(mode: "dev" | "prod") {
    for (const [k, v] of Configs.entries()) {
        if (workers.has(v.name)) {
            console.error(`Refusing to start duplicate worker for ${v.name}`);
            continue;
        }
        runScriptPrefixed(k, v, mode === "dev" ? "DEV" : "PROD");
    }
}
while (true) {
    scanSync("../")
    await load({export: true}).then(() => {
        console.log("\n\nLoaded environment variables\n")
    });
    const mode = await prompt("Dev or Prod? (d/p)")
    if (mode == "d") {
        console.log("Starting processes")
        startProcesses("dev")
    }
    const answer = await prompt("Next command (q to quit)");
    if (answer === "q") {
        console.log("Quitting...");
        for (const [name] of workers) {
            await stopScript(name);
        }
        Deno.exit(0);
    }

}
