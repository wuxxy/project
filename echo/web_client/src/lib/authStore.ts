import { writable } from 'svelte/store';
import {client} from "$lib/axios";

type User = {
    ID: number;
    CreatedAt: string; // ISO timestamp
    UpdatedAt: string;
    DeletedAt: string | null;
    id: string; // UUID
    username: string;
    verified: boolean;
    avatar_url: string;
    suspended: boolean;
    disable: boolean;
    premium: boolean;
    is_admin: boolean;

} | null;

function createUserStore() {
    const { subscribe, set, update } = writable<User>(null);

    async function fetchUser() {
        try {
            const res = await client.get('/api/me');
            set(res.data);
        } catch (err) {
            try{
                // Attempt to refresh token if available
                    // @ts-ignore
                if(err.response.data.error =="Invalid or expired token"){
                        const refreshRes = await client.post('/auth/token');
                        localStorage.setItem('token', refreshRes.data.access_token);
                        // Retry fetching user after refreshing token
                        return fetchUser();
                    }

            }catch(refreshErr) {
                set(null);
            }
             // not logged in or token invalid
        }
    }
    async function logout(){
        try {
            await client.post('/api/logout');
            localStorage.removeItem('token');
            set(null);
        } catch (err) {
            console.error('Logout failed:', err);
        }
    }
    return {
        subscribe,
        fetchUser,
        logout,
        set // optional manual set
    };
}

export const user = createUserStore();
