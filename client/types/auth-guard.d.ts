import { LitElement } from 'lit';
import { LoginSubmit } from './login-form';
export declare class AuthGuard extends LitElement {
    isAuth: boolean;
    isLoaded: boolean;
    redirect?: string;
    errorMessage: string;
    connectedCallback(): void;
    render(): import("lit-html").TemplateResult<1>;
    _onSubmit: (event: CustomEvent<LoginSubmit>) => void;
    _setRedirect(): void;
    _login(username: string, password: string, rememberMe: boolean): Promise<void>;
    _fetchState(): Promise<void>;
    static styles: import("lit").CSSResult;
}
