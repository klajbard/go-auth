import { LitElement, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

interface Service {
  Name: string;
  URL: string;
}

@customElement("available-services")
export class AvailableServices extends LitElement {
  @property({ type: Boolean })
  isLoaded = false;

  @property()
  redirect?: string = undefined;

  @property()
  services: Service[] = [];

  connectedCallback() {
    this._fetchServices();
    super.connectedCallback();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
  }

  render() {
    return html`
      <div class="title">
        <h1>Available services</h1>
        <div class="logout">
          <button @click=${this._logOut}>Logout</button>
        </div>
      </div>
      <hr />
      ${this.services.length
        ? html`<ul>
            ${this.services.map(
              (service) => html`
                <li>
                  <a href=${service.URL}>${service.Name}</a>
                </li>
              `
            )}
          </ul> `
        : ``}
    `;
  }

  async _fetchServices() {
    this.isLoaded = false;
    try {
      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_URL}/services`,
        { credentials: "include" }
      );
      if (!response.ok) {
        console.log(response);
        throw new Error(`Error! status: ${response.status}`);
      }
      this.services = await response.json();
    } catch (error) {
      console.log(error);
    } finally {
      this.isLoaded = true;
    }
  }

  async _logOut() {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_URL}/logout`,
        { method: "post", credentials: "include" }
      );
      if (!response.ok) {
        console.log(response);
        throw new Error(`Error! status: ${response.status}`);
      }
      console.log(response);
      window.location.assign("/logout");
    } catch (error) {
      console.error("Something went wrong", error);
    }
  }

  static styles = css`
    h1 {
      font-size: 2rem;
      margin: 0;
    }

    hr {
      width: 100%;
    }

    :host {
      display: block;
      width: 100%;
    }

    .title {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      gap: 1rem;
      align-items: start;
    }

    ul {
      list-style: none;
      margin: 0;
      padding: 0;
    }

    a {
      color: var(--white);
    }

    .logout {
      text-align: end;
    }

    button {
      padding: 0.25rem 1rem;
      border: none;
      border-radius: 0.5rem;
      background: var(--purple-light);
      color: var(--white);
      font-weight: bold;
      font-size: 1.2rem;
      cursor: pointer;
    }
  `;
}
