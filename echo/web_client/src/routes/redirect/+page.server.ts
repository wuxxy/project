import type { PageServerLoad } from './$types';
import { fail } from '@sveltejs/kit';
import axios from 'axios';

export const load: PageServerLoad = async ({ url }) => {
    const code = url.searchParams.get('code');
    const state = url.searchParams.get('state');

    if (!code) {
        return { error: 'Missing code' };
    }

    try {
        const res = await axios.post('http://localhost:5000/api/token', {
            code,
            service_id: '53b7600d-80e0-4109-8a54-9be631691020',
            service_secret: '4bc3746b46dae9f816059ad5d9e80126e7a12369',
            state
        }, {
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': 'WATERBOTTLETABLETOMATOES'
            }
        });
        console.log(res.data)
        const { access_token } = res.data;

        if (!access_token || !state) {
            throw new Error('Missing token or redirect URI');
        }

        return {
            access_token,
            state
        };
    } catch (err: any) {
        console.error('Token exchange failed:', err?.response?.data || err.message);

        return {
            error: err?.response?.data?.message || err.message || 'Unknown error'
        };
    }
};
