import {Component} from "@angular/core";
import {toSignal} from "@angular/core/rxjs-interop";
import {AuthState, SessionService} from "../services/session.service";
import {AuthService} from "../services/auth.service";

@Component({
    selector: "app-home-page",
    imports: [],
    templateUrl: "./home-page.html",
    styleUrl: "./home-page.scss",
})
export class HomePage {
    public readonly authData;

    constructor(
        private authService: AuthService,
        private readonly sessionSvc: SessionService,
    ) {
        this.authData = toSignal(this.sessionSvc.getAuthDataObservable(), {
            requireSync: true,
        });
    }

    public signIn(): void {
        this.authService.login();
    }

    protected readonly AuthState = AuthState;
}
