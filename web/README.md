# Employee Directory Web

Astro static frontend with React islands for the directory, employee details, and analytics. It talks to the Go API at `http://localhost:8080` by default.

## 🚀 Project Structure

Important routes:

- `/` — searchable and paginated directory
- `/person?id={uuid}` — employee details
- `/analytics` — headcount and movement summaries

## 🧞 Commands

All commands are run from the root of the project, from a terminal:

| Command | Action |
| :--- | :--- |
| `npm install` | Install dependencies |
| `npm run dev` | Start local server at `localhost:4321` |
| `npm run check` | Run Astro and TypeScript checks |
| `npm run build` | Build static site to `./dist/` |
| `npm run preview` | Preview the production build |
