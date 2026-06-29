import {Component, OnInit, signal} from "@angular/core";
import {Router} from "@angular/router";
import {CommonModule} from "@angular/common";
import {AuthService} from "../services/auth.service";
import {SessionService} from "../services/session.service";

@Component({
    selector: "app-callback",
    standalone: true,
    imports: [CommonModule],
    templateUrl: "./callback.component.html",
})
export class CallbackComponent implements OnInit {
    readonly loading = signal(true);
    readonly error = signal<string | null>(null);

    constructor(
        public readonly router: Router,
        private readonly authService: AuthService,
        private readonly sessionService: SessionService,
    ) {}

    ngOnInit(): void {
        this.handleCallback();
    }

    private async handleCallback(): Promise<void> {
        try {
            const searchParams = new URLSearchParams(window.location.search);

            const error = searchParams.get("error");
            const errorDescription = searchParams.get("error_description");

            if (error) {
                this.error.set(
                    `Authentication error: ${error}. ${errorDescription || ""}`,
                );
                this.loading.set(false);
                return;
            }

            const authCode = searchParams.get("code");

            if (!authCode) {
                this.error.set("No authorization code received");
                this.loading.set(false);
                return;
            }

            await this.authService.handleCallback(authCode);

            const sessionData = await this.sessionService.getSessionData(true);
            if (sessionData.displayName === "") {
                this.router.navigate(["/profile"]);
            } else {
                this.router.navigate(["/"]);
            }
        } catch (err) {
            this.error.set(
                err instanceof Error
                    ? err.message
                    : "Failed to process authentication callback",
            );
            this.loading.set(false);
        }
    }
}
