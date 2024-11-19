import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, searchForWorkspaceRoot } from 'vite';
import path from 'path';

export default defineConfig({
  server: {
      fs: {
          allow: [
              searchForWorkspaceRoot(process.cwd()),
              './bindings/*'
          ]
      }
  },
  resolve: {
      alias: {
          '@clave/backend': path.resolve('./bindings/clave/backend/app'),
      }
  },
  plugins: [sveltekit()]
});