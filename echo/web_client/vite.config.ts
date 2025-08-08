import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
    server: {
        port: 6001, // Change this to your desired port
    },
    preview: {
        port: 6001, // Change this to your desired port for preview
    }
});
