import {Component, OnInit, Signal, signal} from "@angular/core";
import {toSignal} from "@angular/core/rxjs-interop";
import {PageState} from "../page-state";
import {
    FormBuilder,
    FormGroup,
    FormsModule,
    ReactiveFormsModule,
    Validators,
} from "@angular/forms";
import {
    MatError,
    MatFormField,
    MatInput,
    MatLabel,
} from "@angular/material/input";
import {MatSnackBar, MatSnackBarConfig} from "@angular/material/snack-bar";
import {Router} from "@angular/router";
import {SessionData, SessionService} from "../services/session.service";

@Component({
    selector: "app-profile-page",
    imports: [
        FormsModule,
        MatError,
        MatFormField,
        MatInput,
        MatLabel,
        ReactiveFormsModule,
    ],
    templateUrl: "./profile-page.html",
    styleUrl: "./profile-page.scss",
})
export class ProfilePage implements OnInit {
    public pg = signal<PageState<SessionData>>({state: "LOADING"});
    public deleting = signal(false);
    public form: FormGroup;
    protected readonly formStatus: Signal<string>;

    private config: MatSnackBarConfig = {
        verticalPosition: "top",
        horizontalPosition: "center",
        duration: 1000,
    };

    constructor(
        private readonly sessionSvc: SessionService,
        private readonly matSnackBar: MatSnackBar,
        private readonly router: Router,
        formBuilder: FormBuilder,
    ) {
        this.form = formBuilder.group({
            displayName: ["", [Validators.required, Validators.maxLength(50)]],
        });
        this.formStatus = toSignal(this.form.statusChanges, {
            initialValue: this.form.status,
        });
    }

    ngOnInit() {
        this.sessionSvc
            .getSessionData(true)
            .then((sd) => {
                this.form.get("displayName")?.setValue(sd.displayName);
                this.pg.set({state: "LOADED", data: sd});
            })
            .catch((e) => {
                this.pg.set({state: "ERROR", error: e});
            });
    }

    public updateDisplayName(): void {
        if (this.form.invalid) {
            this.form.markAllAsTouched();
            return;
        }

        const displayName = this.form.get("displayName")?.value;
        this.form.disable();

        this.sessionSvc
            .setDisplayName(displayName)
            .then(() => {
                const ref = this.matSnackBar.open(
                    `Display name updated`,
                    "Dismiss",
                    this.config,
                );
                ref.afterDismissed()
                    .pipe()
                    .subscribe(() => {
                        this.router.navigate(["/"]);
                    });
            })
            .catch((e) => {
                alert(e);
            })
            .finally(() => {
                this.form.enable();
            });
    }

    public deleteAccount(): void {
        if (!confirm("Delete account?\nData will not be recoverable")) {
            return;
        }

        this.deleting.set(true);
        this.sessionSvc
            .deleteAccount()
            .then(() => {
                this.router.navigate(["/"]);
            })
            .catch((e) => {
                alert(e);
            });
    }
}
