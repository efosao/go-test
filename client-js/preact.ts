import { LitElement, html } from "lit";
import { SignalWatcher, signal } from "@lit-labs/preact-signals";

export class SignalExample extends SignalWatcher(LitElement) {
  count = signal(0);
  render() {
    return html`
      <div
        class="flex gap-2 bg-pink-200 p-4 my-2 rounded-lg border-2 border-red-300"
      >
        <button class="px-4 py-2" @click=${this._onClick}>+</button>
        <p
          class="flex text-lg font-bold text-pink-800 dark:text-pink-800 w-36 items-center"
        >
          The count is ${this.count.value}
        </p>
      </div>
    `;
  }

  protected createRenderRoot() {
    return this;
  }

  private _onClick() {
    this.count.value = Math.floor(this.count.value + Math.random() * 10);
  }
}

customElements.define("x-greeting", SignalExample);
