# Volo Landing Page

Marketing and download page for Volo. Raycast-inspired dark theme with animated hero, product mockups, and download CTAs.

## Setup

```bash
npm install
```

## Development

```bash
npm run dev
```

Opens at `http://localhost:5173`.

## Build

```bash
npm run build
```

Outputs static files to `dist/`. Deploy anywhere (Vercel, Netlify, Caddy on VM).

## Preview (production build locally)

```bash
npm run build && npx vite preview
```

## Structure

```
src/
├── components/
│   ├── Nav.tsx              Fixed nav with download CTA
│   ├── Hero.tsx             Animated rays background, headline, CTAs
│   ├── ProductShowcase.tsx  Desktop app + extension mockups
│   ├── Features.tsx         6 feature cards
│   ├── HowItWorks.tsx      3-step guide
│   ├── Download.tsx         Final download CTAs
│   └── Footer.tsx           Links, credits
├── App.tsx                  Page composition
├── main.tsx                 Entry point
└── index.css                Tailwind + custom animations
```

## Design

- **Style:** Raycast-inspired (dark, bold headline, animated background)
- **Hero:** Rotating blue/teal light rays, gradient text, white CTA button
- **Mockups:** Actual app UI rendered in CSS (Discord colors, Codex layout)
- **Palette:** Near-black background (`#0a0a0f`), Discord-colored mockups

## Deployment

Static files — deploy to any static host:

```bash
# Caddy (on VM alongside API)
# Caddyfile: volo.yourdomain.com { root * /var/www/landing; file_server }

# Or Vercel
npx vercel --prod

# Or Netlify
npx netlify deploy --prod --dir=dist
```

## Tech Stack

- React 19 + TypeScript
- Vite
- Tailwind CSS
