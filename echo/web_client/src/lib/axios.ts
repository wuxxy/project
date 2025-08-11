import axios from 'axios';
import { browser } from '$app/environment';

// Create a custom axios instance
const client = axios.create({
    baseURL: 'http://localhost:5002',
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    }
})
const authxios = axios.create({
    baseURL: 'http://localhost:5000',
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'X-CSRF-Token': 'WATERBOTTLETABLETOMATOES' // Replace with your actual CSRF token if needed
    }
})

// Add a request interceptor to include auth token if available
client.interceptors.request.use(config => {
    // Only access localStorage in browser environment
    if (browser) {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
}, error => {
    return Promise.reject(error);
});
authxios.interceptors.request.use(config => {
    // Only access localStorage in browser environment
    if (browser) {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
}, error => {
    return Promise.reject(error);
});

export { client, authxios };
