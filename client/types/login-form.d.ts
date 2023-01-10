import { LitElement } from 'lit';
export interface LoginSubmit {
    username: string;
    password: string;
    rememberMe: boolean;
}
export declare class LoginForm extends LitElement {
    redirect?: string;
    errorMessage: string;
    render(): import("lit-html").TemplateResult<1>;
    _onSubmit: (event: SubmitEvent) => void;
    static styles: import("lit").CSSResult;
}
