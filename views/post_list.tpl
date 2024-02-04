{{range .Posts}}
<div
  {{if .IsPinned}}
  class="search_row group relative mb-2 rounded-md border-0 border-pink-200 bg-pink-100 dark:border-prussian-blue-900 dark:bg-black dark:text-white"
  {{else}}
  class="search_row group relative mb-2 rounded-md border-0 border-orange-200  bg-orange-200 text-black dark:border-slate-700 dark:bg-slate-700 dark:text-white"
  {{end}}
  onclick="utils.toggleOpenState('cbx{{.ID}}', 'desc{{.ID}}')"
>
  {{if .IsPinned}}
  <img
    alt="pin"
    src="/public/images/pin.svg"
    class="absolute right-[-8px] top-[-8px] h-6 w-6 text-red-400"
  />
  {{end}}
  <div class="cursor-pointer flex h-32 items-center space-x-2 px-2">
      {{if .Thumbnail}}
      <span
        class="rounded-full initials inline-flex h-[40px] w-[40px] my-2 shrink-0 items-center justify-center overflow-hidden"
      >
        <img
          alt="{{.CompanyName}} logo"
          class="overflow-hidden"
          loading="lazy"
          src="{{.Thumbnail}}"
          width="40"
          height="40"
        />
      </span>
      {{else}}
      <span
        class="bg-teal-300 rounded-full initials inline-flex h-[40px] w-[40px] my-2 shrink-0 items-center justify-center overflow-hidden"
      >
        {{.GetInitials}}
      </span>
      {{end}}
      <div class="flex grow flex-col min-w-10">
        <p class="text-black line-clamp-1 font-semibold lg:line-clamp-2">
          {{- .CompanyName}}
        </p>

        <p class="text-black line-clamp-1 font-bold md:line-clamp-2">
          {{- .Title}}</p>
        <p class="text-black">{{.Location}}</p>

        <!-- <p>
          {isRemote ? (
            <>
              Remote&nbsp;
              {post.location && (
                <span class="text-orange-700 dark:text-orange-300">
                  ({post.location.trim()})
                </span>
              )}
            </>
          ) : (
            <>{post.location}</>
          )}
        </p> -->
    </div>
    <div class="tag-container">
    {{range .Tags}}
    <button
      class="inline cursor-pointer rounded-md bg-white px-2 font-semibold text-pink-950 transition-colors duration-300 hover:bg-blue-100 hover:text-black my-[2px]"
      type="button"
      onclick="utils.halt(event)"
    >
    {{.}}
    </button>
    {{end}}
    </div>
    <span class="m-2">{{.TimeSinceCreated}}</span>
    <span class="btn-apply done">Applied</span>
    <!-- {alreadyApplied ? (
    ) : (
      <a
        class="btn-apply"
        href={applyUrl}
        target={isExternalApplyLink ? "_blank" : undefined}
        onclick="utils.halt(event)"
      >
        <span>Apply</span>
      </a>
    )} -->
  </div>

  <div onclick="utils.halt(event)" class="flex flex-col items-center justify-center">
    <form>
    <input
      aria-label="toggle show description"
      class="peer hidden"
      type="checkbox"
      id="cbx{{.ID}}"
    />
    </hidden>
    <div
      class="hidden items-center justify-center peer-checked:flex"
      hx-get="/partials/posts/details/{{.ID}}"
      hx-indicator="#htmx{{.ID}}"
      hx-swap="outerHTML transition:true"
      hx-trigger="change"
      id="desc{{.ID}}"
    >
      <img
        alt="loading"
        id="htmx{{.ID}}"
        class="htmx-indicator mb-24 h-12"
        src="/public/images/bars-loader.svg"
        height="48"
      />
    </div>
  </div>
</div>
{{end}}

<div
  {{if .SelectedTagsStr}}
  hx-post="/partials/posts/search/{{.Page}}?tags={{.SelectedTagsStr}}"
  {{else}}
  hx-post="/partials/posts/search/{{.Page}}"
  {{end}}
  hx-swap="outerHTML"
  hx-trigger="revealed"
  id="nextPageLoaderId_{{.Page}}"
>
  <div class="htmx-indicator flex flex-col gap-2 items-center justify-center">
    <div class="bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"></div>
    <div class="bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"></div>
    <div class="bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"></div>
    <div class="bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"></div>
    <div class="bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"></div>
    <div class="bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"></div>
  </div>
</div>
