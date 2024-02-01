<div class="mb-4 h-10">
  <select
    aria-hidden="true"
    class="hide slim-select"
    hx-post="/partials/posts/search/0"
    hx-target="#post-list"
    hx-trigger="change"
    id="tags"
    name="tags"
    multiple
    tabindex="-1"
    x-init="window.utils.loadSlimSelect"
  >
    {{range .Tags}}<option value="{{.GetEscapedName}}" {{- if .Selected}}selected{{end}}>{{.Name}}</option>{{end}}
  </select>
</div>

<div id="post-list">
{{template "post_list" .}}
</div>