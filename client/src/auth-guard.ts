import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import {ifDefined} from 'lit/directives/if-defined.js';

import {LoginSubmit} from './login-form'

@customElement('auth-guard')
export class AuthGuard extends LitElement {
  @property({ type: Boolean })
  isAuth = false;

  @property({ type: Boolean })
  isLoaded = false;

  @property()
  redirect?: string = undefined;

  @property({type: String})
  errorMessage: string = ""

  connectedCallback() {
    super.connectedCallback();
    this._fetchState();
    this._setRedirect();
  }

  render() {
    return html`
      ${this.isLoaded ? (
        this.isAuth ? (
          html`
            <success-message redirect=${ifDefined(this.redirect)}></success-message>
          `
        ) : (
          html`
            <login-form redirect=${ifDefined(this.redirect)} @login-submit=${this._onSubmit} errorMessage=${this.errorMessage}>
            </login-form>
          `
        )
      ) : "loading..."}
    `
  }

  _onSubmit = (event: CustomEvent<LoginSubmit>) => {
    const { username, password, rememberMe } = event.detail
    this._login(username, password, rememberMe)
  }

  _setRedirect() {
    const params = new URLSearchParams(window.location.search)
    this.redirect = decodeURIComponent(params.get("redirect") || "")
    console.log(`Will redirect to ${this.redirect}`)
  }

  async _login(username: string, password: string, rememberMe: boolean) {
    this.isLoaded = false;
    this.errorMessage = "";
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/login`, {
        method: "post",
        credentials: 'include',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },

        body: JSON.stringify({
          username,
          password,
          rememberMe,
        })
      })
      if (!response.ok) {
        this.errorMessage = "Invalid username or password!"
      } else {
        this.isAuth = true;
      }
    } catch (error) {
      this.errorMessage = "Something went wrong!"
      this.isAuth = false;
    } finally {
      this.isLoaded = true;
    }
  }

  async _fetchState() {
    this.isLoaded = false;
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/hello`, {credentials: 'include'})
      if (!response.ok) {
        console.log(response);
        throw new Error(`Error! status: ${response.status}`);
      }
      console.log(response)
      this.isAuth = true;
    } catch (error) {
      this.isAuth = false;
    } finally {
      this.isLoaded = true;
    }
  }

  static styles = css`
  :host {
    margin: auto;
    display: flex;
  }
  `
}
