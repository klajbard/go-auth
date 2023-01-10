import { LitElement } from 'lit';
export declare class SuccessMessage extends LitElement {
    redirect?: string;
    timeout?: number;
    connectedCallback(): void;
    disconnectedCallback(): void;
    render(): import("lit-html").TemplateResult<1>;
    static styles: import("lit").CSSResult;
}
