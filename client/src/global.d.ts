import { SuccessMessage } from "./success-message";
import { LoginForm } from "./login-form";
import { AuthGuard } from "./auth-guard";
import { AvailableServices } from "./available-services";

declare global {
  interface HTMLElementTagNameMap {
    "success-message": SuccessMessage;
    "login-form": LoginForm;
    "auth-guard": AuthGuard;
    "available-services": AvailableServices;
  }
}
