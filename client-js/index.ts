/// <reference lib="dom" />
/// <reference lib="dom.iterable" />
import htmx from "htmx.org";
import Cookies from "js-cookie";
// import "./preact";
import "./react";

declare global {
  interface Window {
    __htmx: any;
    utils: any;
  }
}

window.__htmx = htmx;

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

  // Animate the page title
  if (pgTitle) {
    pgTitle.animate(
      [
        { opacity: 0, transform: "translateX(50px)" },
        { opacity: 0.5, transform: "translateX(-20px)" },
        { opacity: 1, transform: "translateX(0)" },
      ],
      {
        delay: 200,
        duration: 350,
        fill: "forwards",
        easing: "ease-in",
      }
    ).onfinish = () => {
      console.debug("Animation finished");
      pgTitle.style.removeProperty("opacity");
    };
  } else {
    console.error("pgTitle not found");
  }
}

window.document.addEventListener("DOMContentLoaded", () => {
  console.debug("DOMContentLoaded");
  animatePageTitle();
  utils.init();
  loadEventListeners();
});

htmx.onLoad(function (e) {
  console.debug("htmx.onLoad");
  if (e.id === "page-content") animatePageTitle();
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
