// Worker that runs an arbitrary command and streams stdout/stderr back line-by-line
// deno run perms needed by the *main* process: --allow-run --allow-read --allow-env
// (the worker inherits what you grant to it via Worker options)

import { TextLineStream } from "jsr:@std/streams@0.224.0/text-line-stream";

type StartMsg = {
    type: "start";
    payload: { path: string; name: string; command: string };
};
type StopMsg = { type: "stop" };

let child: Deno.Child | null = null;
let name = "";

self.onmessage = async (ev: MessageEvent<StartMsg | StopMsg>) => {
    const msg = ev.data;

    if (msg.type === "stop") {
        if (!child) return;

        try {
            if (Deno.build.os === "windows") {
                // Hard stop whole tree on Windows
                await new Deno.Command("taskkill", {
                    args: ["/PID", String(child.pid), "/T", "/F"],
                    stdout: "null",
                    stderr: "null",
                }).output();
            } else {
                // Graceful → forceful on POSIX
                try { child.kill("SIGINT"); } catch {}
                await new Promise(r => setTimeout(r, 400));
                try { child.kill("SIGTERM"); } catch {}
                await new Promise(r => setTimeout(r, 800));
                try { child.kill("SIGKILL"); } catch {}
            }
        } catch (e) {
            self.postMessage({ type: "error", name, error: String(e) });
        }
        return;
    }

    if (msg.type === "start") {
        const { path, command, name: n } = msg.payload;
        name = n;

        const isWin = Deno.build.os === "windows";
        const shell = isWin ? "cmd" : "sh";
        const args = isWin
            ? ["/d", "/s", "/c", command]
            : ["-c", `exec ${command}`]; // exec = replace shell with target

        child = new Deno.Command(shell, {
            args,
            cwd: path,
            stdout: "piped",
            stderr: "piped",
            stdin: "null",
        }).spawn();

        const sink = (kind: "OUT" | "ERR") =>
            new WritableStream<string>({ write: (line) => self.postMessage({ type: "line", name, kind, line }) });

        child.stdout.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream()).pipeTo(sink("OUT")).catch(() => {});
        child.stderr.pipeThrough(new TextDecoderStream()).pipeThrough(new TextLineStream()).pipeTo(sink("ERR")).catch(() => {});

        try {
            const status = await child.status;
            self.postMessage({ type: "exit", name, success: status.success, code: status.code });
        } catch (e) {
            self.postMessage({ type: "error", name, error: String(e) });
            self.postMessage({ type: "exit", name, success: false, code: -1 });
        }
    }
};
