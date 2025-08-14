import {join} from "https://deno.land/std@0.107.0/path/mod.ts";

const IGNORED = new Set([".git", "node_modules", ".idea", ".venv"]);
export type CpConfig = {
    version?: string;
    description?: string;
    scripts?: Record<string, string>;
    name: string;
}
export const Configs= new Map<string, CpConfig>()
export type FileJSON = {
    type: "file";
    name: string;
    path: string;
    isExecutable: boolean;
};

export type DirJSON = {
    type: "dir";
    name: string;
    path: string;
    children: Array<FileJSON | DirJSON>;
    scripts?: Record<string, string>; // from .cpconfig
};

export abstract class FSNode {
    constructor(public name: string, protected pathAbs: string) {}
    getPath(): string {
        return this.pathAbs;
    }
}

export class FSFile extends FSNode {
    constructor(name: string, pathAbs: string, public isExecutable = false) {
        super(name, pathAbs);
    }
    toJSON(): FileJSON {
        return {
            type: "file",
            name: this.name,
            path: this.getPath(),
            isExecutable: this.isExecutable,
        };
    }
}

export class FSDir extends FSNode {
    files = new Set<FSFile>();
    directories = new Set<FSDir>();
    scripts?: Map<string, string>; // parsed from .cpconfig

    constructor(name: string, pathAbs: string) {
        super(name, pathAbs);
    }

    addFile(f: FSFile) {
        this.files.add(f);
    }
    addDir(d: FSDir) {
        this.directories.add(d);
    }

    toJSON(): DirJSON {
        const children: Array<FileJSON | DirJSON> = [];
        for (const f of this.files) children.push(f.toJSON());
        for (const d of this.directories) children.push(d.toJSON());
        const scriptsObj =
            this.scripts && this.scripts.size ? Object.fromEntries(this.scripts) : undefined;
        return {
            type: "dir",
            name: this.name,
            path: this.getPath(),
            children,
            ...(scriptsObj ? { scripts: scriptsObj } : {}),
        };
    }
}

/** Build a tree synchronously from a root path */
export function scanSync(rootPath: string): FSDir {
    const root = new FSDir("/", rootPath);
    scanDirSync(root);
    return root;
}

function scanDirSync(dir: FSDir) {
    readCpConfigSync(dir);

    for (const entry of Deno.readDirSync(dir.getPath())) {
        if (IGNORED.has(entry.name)) continue;

        const full = join(dir.getPath(), entry.name);

        if (entry.isFile) {
            dir.addFile(new FSFile(entry.name, full, false));
        } else if (entry.isDirectory) {
            const child = new FSDir(entry.name, full);
            dir.addDir(child);
            scanDirSync(child);
        }
    }
}

function readCpConfigSync(dir: FSDir) {
    const cfg = join(dir.getPath(), ".cpconfig");
    try {
        const txt = Deno.readTextFileSync(cfg);
        const data = JSON.parse(txt) as CpConfig;
        const m = new Map<string, string>();

        for (const [k, v] of Object.entries(data.scripts ?? {})) {
            m.set(k, v);
        }
        Configs.set(dir.getPath(), data);
        if (m.size) dir.scripts = m;
    } catch {
    }
}
