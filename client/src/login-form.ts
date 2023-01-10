import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

export interface LoginSubmit {
  username: string
  password: string
  rememberMe: boolean
}

@customElement('login-form')
export class LoginForm extends LitElement {
  @property()
  redirect?: string = undefined;

  @property({type: String})
  errorMessage: string = "";

  render() {
    return html`
      <div class="left">
        <h1>Welcome</h1>
        ${this.redirect ? html`<span class="message">Authenticating for <strong>${new URL(this.redirect).hostname}</strong></span>` : ``}
      </div>
      <div class="right">
        <h1>Login</h1>
        ${this.errorMessage ? html`<span class="error message">${this.errorMessage}</span>` : ``}
        <form @submit=${this._onSubmit}>
          <label>
            <span>Username</span>
            <input type="text" name="username" required />
          </label>
          <label>
            <span>Password</span>
            <input type="password" name="password" required />
          </label>
          <label class="checkbox">
            <input type="checkbox" />
            <span>Remember me</span>
          </label>
          <input type="submit" value="Login" />
        </form>
      </div>
    `
  }

  _onSubmit = (event: SubmitEvent) => {
    event.preventDefault();
    if (!event.currentTarget) {
      return
    }
    const formElement = event.currentTarget as HTMLFormElement
    const username = (formElement as any)[0].value
    const password = (formElement as any)[1].value
    const rememberMe = (formElement as any)[2].checked

    const customEv = new CustomEvent('login-submit', {
      bubbles: true,
      composed: true,
      detail: { username, password, rememberMe },
    })

    this.dispatchEvent(customEv)
  }


  static styles = css`
    :host {
      display: grid;
      margin: 2rem auto;
      grid-template-rows: auto 1fr;
      overflow: hidden;
      border-radius: 0.5rem;
      box-shadow: var(--shadow);
    }

    @media screen and (min-width: 640px) {
      :host {
        grid-template-columns: 1fr 1fr;
        grid-template-rows: auto;
      }
    }

    .left {
      color: var(--white);
      background-color: var(--blue-light);
    }

    .right {
      background-color: var(--white-transparent);
    }

    .right h1 {
      font-family: "Poppins Bold";
      color: var(--blue-dark);
    }

    .left, .right {
      box-sizing: border-box;
      min-width: 320px;
      padding: 2rem;
    }

    form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    label {
      display: flex;
      flex-direction: column;
    }

    input {
      box-sizing: border-box;
      height: 2.5rem;
      padding: 0.5rem;
      border-radius: 0.5rem;
      border: none;
      background: var(--white);
    }

    input[required]:invalid {
      border: 2px solid red;
    }

    input[required]:valid {
      border: 2px solid green;
    }

    input[type="submit"] {
      color: var(--white);
      font-size: 1.2rem;
      background: var(--blue-dark);
      cursor: pointer;
    }

    input[type="checkbox"] {
      width: 1rem;
      height: 1rem;
      padding: 0;
      margin-top: 2px;
      margin-right: 0.5rem;
    }

    .message {
      display: block;
      margin-bottom: 0.5rem;
    }

    .error {
      color: red;
      font-weight: bold;
    }

    .checkbox {
      flex-direction: row;
    }
  `
}
