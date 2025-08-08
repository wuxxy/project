import type { PageServerLoad } from './$types';
import axios from "axios";
import { redirect } from '@sveltejs/kit';
export const load: PageServerLoad = async ({ params, url }) => {
    const service_id = url.searchParams.get('s');
    const redirect_url = url.searchParams.get('r');
    if(!service_id || !redirect_url) {
        return {
            error: true
        }
    }
    let response;
    try {
        response = await axios.get(`http://localhost:5000/api/service?service_id=${service_id}&redirect_url=${redirect_url}`, {
            headers: {
                'Accept': 'application/json',
                'X-CSRF-Token': "WATERBOTTLETABLETOMATOES"
            }
        });
    } catch (e) {
        return {
            error: true,
        }
    }
    return {
        service: response.data
    };
};