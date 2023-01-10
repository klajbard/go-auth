import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'


const services = [
  {name: 'Home Assistant', url: 'http://localhost:8123'},
  {name: 'Media Server', url: 'http://localhost:8111'},
  {name: 'Image Storage', url: 'http://localhost:8000'},
]

@customElement('available-services')
export class AvailableServices extends LitElement {
  @property()
  redirect?: string = undefined;

  connectedCallback() {
    super.connectedCallback();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
  }

  render() {
    return html`
      <ul>
        ${services.map((service) => html`
          <li>
            <a href=${service.url}>${service.name}</a>
          </li>
        `)}
      </ul>
      <div class="logout">
        <button @click=${this._logOut}>Logout</button>
      </div>
    `
  }

  async _logOut() {
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/logout`, {        method: "post",
      credentials: 'include',
    })

      if (!response.ok) {
        console.log(response);
        throw new Error(`Error! status: ${response.status}`);
      }
      console.log(response)
      window.location.assign('/logout');
    } catch (error) {
      console.error("Something went wrong", error)
    }
  }

  static styles = css`
    :host {
      display: block;
      width: 100%;
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
      margin-top: 3rem;
      padding: 0.5rem;
      border: none;
      border-radius: 0.5rem;
      background: var(--purple-light);
      color: var(--white);
      font-weight: bold;
      font-size: 1.2rem;
      cursor: pointer;
    }
    `
}
