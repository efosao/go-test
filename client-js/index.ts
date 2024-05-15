/// <reference lib="dom" />
/// <reference lib="dom.iterable" />
import htmx from "htmx.org";
import Cookies from "js-cookie";
// import "./preact";
// import "@sudodevnull/datastar";
import "unpoly";
import Alpine from "alpinejs";
import "./react";
import "./components";

declare global {
  interface Window {
    __htmx: any;
    alpine: any;
    fn: any;
    utils: any;
  }
}

window.__htmx = htmx;
window.alpine = Alpine;
window.fn = (...args: any) => {
  console.log("fn", args);
};

const ENABLE_HTMX_DEBUG = process.env.NODE_ENV !== "production" && false;

const selectors = {
  pgTitle: "page-title",
};

const utils = {
  halt: (e: Event) => {
    e.preventDefault();
    e.stopPropagation();
  },

  init: async () => {
    console.debug("Initializing utils");
    utils.configureDefaultTheme(false);

    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
    Cookies.set("tz", tz);
  },

  isCurrentThemeDark: () => {
    const theme = Cookies.get("theme") || "system";
    let isDark =
      theme === "system"
        ? window.matchMedia("(prefers-color-scheme: dark)").matches
        : theme === "dark";
    return isDark;
  },

  configureDefaultTheme: (animate = true) => {
    const darkToggleId = "dark-mode-toggle";
    const lightToggleId = "light-mode-toggle";

    const darkToggle = document.getElementById(darkToggleId);
    const lightToggle = document.getElementById(lightToggleId);
    const duration = animate ? 300 : 0;

    function transitionToggles(
      outElement: HTMLElement,
      inElement: HTMLElement
    ) {
      outElement.animate(
        [
          { opacity: 1, transform: "rotate(0deg)" },
          { opacity: 0, transform: "rotate(360deg)" },
        ],
        { duration, fill: "forwards", easing: "ease-in-out" }
      );
      inElement.classList.remove("hidden");
      inElement.animate(
        [
          { opacity: 0, transform: "rotate(0deg)" },
          { opacity: 1, transform: "rotate(360deg)" },
        ],
        {
          duration,
          fill: "forwards",
        }
      );
    }

    if (!darkToggle || !lightToggle) return;

    if (utils.isCurrentThemeDark()) {
      document.getElementsByTagName("html")[0].classList.add("dark");
      transitionToggles(darkToggle, lightToggle);
    } else {
      document.getElementsByTagName("html")[0].classList.remove("dark");
      transitionToggles(lightToggle, darkToggle);
    }
  },

  toggleCurrentTheme: () => {
    const newTheme = utils.isCurrentThemeDark() ? "light" : "dark";
    utils.setTheme({ target: { value: newTheme } });
  },

  toggleOpenState: (checkboxId: string, descriptionId: string) => {
    const checkbox: HTMLInputElement = document.getElementById(
      checkboxId
    ) as HTMLInputElement;
    checkbox.checked = !checkbox.checked;
    const description: HTMLElement = document.getElementById(
      descriptionId
    ) as HTMLElement;
    if (description) {
      window.__htmx.trigger("#" + descriptionId, "change");
    }
  },

  setTheme: (event: { target: { value: "light" | "dark" | "system" } }) => {
    const theme = event?.target?.value;
    const options = ["light", "dark", "system"];

    if (!options.includes(theme)) throw new Error("Unknown theme: " + theme);
    Cookies.set("theme", theme);
    utils.configureDefaultTheme();
  },
};

window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", ({ matches }) => {
    utils.configureDefaultTheme();
  });

window.utils = utils;

function loadEventListeners() {
  console.debug("DOMContentLoaded");

  const showButton = document.getElementById("dialog_button");

  window.addEventListener("click", (e) => {
    const dialog = document.getElementById(
      "dialog"
    ) as HTMLDialogElement | null;
    if (e.target === dialog) {
      dialog?.close();
    }
  });

  showButton?.addEventListener("click", () => {
    const dialog = document.getElementById(
      "dialog"
    ) as HTMLDialogElement | null;
    dialog?.showModal();
  });

  const closeButtons = document.querySelectorAll("dialog button.close");

  closeButtons.forEach((el) => {
    el.addEventListener("click", () => {
      const dialog = document.getElementById(
        "dialog"
      ) as HTMLDialogElement | null;
      dialog?.close();
    });
  });
}

function animatePageTitle() {
  const pgTitle = document.getElementById(selectors.pgTitle);
  pgTitle?.style.setProperty("opacity", "0");
  pgTitle?.classList.add("fade-in-bottom");
}

window.document.addEventListener("DOMContentLoaded", () => {
  console.debug("DOMContentLoaded");
  Alpine.start();
  animatePageTitle();
  utils.init();
  loadEventListeners();
});

htmx.onLoad(function (e) {
  console.debug("htmx.onLoad");
  if (e.id === "page-content") {
    animatePageTitle();
    utils.configureDefaultTheme(false);
  }
  loadEventListeners();
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
