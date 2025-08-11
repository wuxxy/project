export type ChannelID = string;
export type SocketID = string;
export type EVENT = "ping" | "message" | "error" | "hello" | "join"
export type MESSAGE = {
    e: EVENT; // Event
    d: any; // Data
}
export const MS = 1;
export const SECOND = 1000*MS;
export const MINUTE = 60*SECOND;
export const HOUR = 60*MINUTE;
export const DAY = 24*HOUR;