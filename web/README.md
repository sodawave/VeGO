# VeGo public showcase

Static site that presents the Alpha pre-release to the world: idea, Go ↔ `.vego` compare, pipeline, BMAD trail, and clone commands.

```bash
cd web
python3 -m http.server 8080
```

Open http://localhost:8080

Deploy the contents of `web/` to GitHub Pages (`/docs` alternative: point Pages at `/web`), Netlify, or Cloudflare Pages. Sample IR matches `src/pkg/transpiler/testdata/hello_http.vego`.
