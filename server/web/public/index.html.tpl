<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>MarNotes - JSON Server Dashboard</title>
    <link rel="stylesheet" href="/static/dist/app.css" />
    <script
      crossorigin
      src="https://unpkg.com/react@18/umd/react.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"
    ></script>
    <script>
      window.__JSON_SERVER_PORT__ = "{{.JSONPort}}";
    </script>
  </head>
  <body>
    <div id="root"></div>
    <script src="/static/dist/app.js"></script>
  </body>
</html>
