import { LitElement, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

@customElement("success-message")
export class SuccessMessage extends LitElement {
  @property()
  redirect?: string = undefined;

  @property()
  timeout?: number = undefined;

  connectedCallback() {
    super.connectedCallback();
    this.timeout = window.setTimeout(() => {
      window.location.assign(this.redirect || "/services");
    }, 1500);
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    if (this.timeout) {
      window.clearTimeout(this.timeout);
      this.timeout = undefined;
    }
  }

  render() {
    return html`
      <h1>Successfully logged in!</h1>
      <hr />
      <span>
        Redirecting to
        ${this.redirect
          ? html`
              <a href=${this.redirect}>${new URL(this.redirect).hostname}</a>
            `
          : html` <a href="/services">Services</a> `}...
      </span>
    `;
  }

  static styles = css`
    :host {
      margin: auto;
      box-sizing: border-box;
      min-width: 320px;
      display: flex;
      padding: 2rem;
      background-color: var(--blue-dark);
      color: var(--white);
      flex-direction: column;
      align-items: flex-start;
      font-size: 1.25rem;
      border-radius: 0.5rem;
      box-shadow: var(--shadow);
    }
    div {
      display: flex;
      min-width: 320px;
      flex-direction: column;
      align-items: center;
    }

    h1 {
      margin: 0;
      font-size: 2rem;
      font-weight: bold;
    }

    hr {
      width: 100%;
    }

    a {
      font-style: normal;
      color: white;
      font-weight: bold;
    }
  `;
}
