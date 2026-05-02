import type { Config } from 'tailwindcss';

export default {
  content: [
    './app/**/*.{ts,tsx,mdx}',
    './content/**/*.{md,mdx}',
    './components/**/*.{ts,tsx}',
  ],
} satisfies Config;
