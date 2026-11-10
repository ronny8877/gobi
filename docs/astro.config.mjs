import { defineConfig } from 'astro/config';
import svelte from '@astrojs/svelte';
import mdx from '@astrojs/mdx';
import tailwind from '@tailwindcss/vite';
export default defineConfig({
  devToolbar: { enabled: false },
  integrations: [svelte(), mdx()],
  vite: {
    plugins: [tailwind()],
    server: {
      proxy: {
        '/gobi-api': {
          target: process.env.GOBI_API_URL || 'http://127.0.0.1:8080',
          rewrite: (path) => path.replace(/^\/gobi-api/, ''),
        },
      },
    },
  },
});
