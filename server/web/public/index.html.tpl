<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>MarNotes - JSON Server Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&display=swap" rel="stylesheet" />
    <link rel="stylesheet" href="https://unpkg.com/normalize.css@8.0.1/normalize.css" />
    <link rel="stylesheet" href="https://unpkg.com/@blueprintjs/icons@6/lib/css/blueprint-icons.css" />
    <link rel="stylesheet" href="https://unpkg.com/@blueprintjs/core@6/lib/css/blueprint.css" />
    <link rel="stylesheet" href="/static/dist/app.css" />
    <script>
      window.__JSON_SERVER_PORT__ = "{{.JSONPort}}";
    </script>
  </head>
  <body>
    <div id="root"></div>
    <script src="/static/dist/app.js"></script>
  </body>
</html>
