<!DOCTYPE html>
<html class="{{.Theme}}" lang="en">
<head>
  <title>{{.Title}}</title>
  <meta charset="UTF-8">
  <meta id="viewport" name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="{{.Description}}">
  <script src="/public/dist/index.js" defer></script>
  <link rel="stylesheet" href="/public/dist/stylesheet.css"></link>
  <style>
    .hide { display: none; }
  </style>
</head>

<body hx-boost="true" class="max-w-4xl mx-auto dark:bg-slate-500">
{{template "partials/header" .}}
<main class="p-2">
  <h1 class="text-3xl font-extrabold mb-4 text-black dark:text-black">{{.Title}}</h1>
  {{embed}}
</main>
{{template "partials/footer" .}}
</body>

</html>