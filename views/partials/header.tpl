<nav class="m-2 text-xl flex justify-between">
  <div>
    <a class="p-2 inline-block" href="/">Home</a>
    <a class="p-2 inline-block" href="/posts">Job Posts</a>
  </div>

  <form class="select-none" onchange="utils.setTheme(event)">
      <label class="flex items-center">
        <span class="text-slate-900 dark:text-slate-900">theme</span>
        <select name="themepicker" class="ml-2 border-none dark:bg-slate-500">
        {{range .ThemeOptions}}
        <option class="bg-transparent" {{if .Selected}}selected{{end}} value="{{.Value}}">{{.Label}}</option>
        {{end}}
        </select>
      </label>
    </form>
</nav>