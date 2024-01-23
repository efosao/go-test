/// <reference lib="dom" />
/// <reference lib="dom.iterable" />
import htmx from "htmx.org";
import Cookies from "js-cookie";
import SlimSelect from "../node_modules/slim-select/src/slim-select/index";
import Alpine from "alpinejs";

declare global {
  interface Window {
    __htmx: any;
    utils: any;
  }
}

window.__htmx = htmx;

const ENABLE_HTMX_DEBUG = process.env.NODE_ENV !== "production" && false;

const utils = {
  halt: (e: Event) => {
    e.preventDefault();
    e.stopPropagation();
  },

  init: async () => {
    console.debug("Initializing utils");
    utils.configureDefaultTheme();

    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
    Cookies.set("tz", tz);
  },

  configureDefaultTheme: () => {
    const theme = Cookies.get("theme") || "system";
    let isDark =
      theme === "system"
        ? window.matchMedia("(prefers-color-scheme: dark)").matches
        : theme === "dark";

    if (isDark) document.getElementsByTagName("html")[0].classList.add("dark");
    else document.getElementsByTagName("html")[0].classList.remove("dark");
  },

  toggleOpenState: (checkboxId: string, descriptionId: string) => {
    const checkbox: HTMLInputElement = document.getElementById(
      checkboxId
    ) as HTMLInputElement;
    checkbox.checked = !checkbox.checked;
    window.__htmx.trigger("#" + descriptionId, "change");
  },

  setTheme: (event: { target: { value: "light" | "dark" | "system" } }) => {
    const theme = event?.target?.value;
    const options = ["light", "dark", "system"];

    if (!options.includes(theme)) throw new Error("Unknown theme: " + theme);
    Cookies.set("theme", theme);
    utils.configureDefaultTheme();
  },

  loadSlimSelect: () => {
    setTimeout(() => {
      for (const select of document.getElementsByClassName("slim-select")) {
        /* @ts-ignore */
        if (typeof select.slim !== "undefined") {
          /* @ts-ignore */
          select.slim.destroy();
        }
        select.classList.remove("hide");
        new SlimSelect({
          select: select as HTMLSelectElement,
          settings: {
            isMultiple: true,
            maxSelected: 5,
            searchHighlight: true
          }
        });
      }
    }, 50);
  }
};

window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", ({ matches }) => {
    utils.configureDefaultTheme();
  });

window.utils = utils;

window.document.addEventListener("DOMContentLoaded", () => {
  Alpine.start();
  utils.init();
});

if (htmx && ENABLE_HTMX_DEBUG) {
  window.__htmx.logger = function (
    elt: Element,
    event: string,
    data: Record<string, any>
  ) {
    if (console) {
      console.debug(event, elt, data);
    }
  };
}
