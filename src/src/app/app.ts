import {
    ChangeDetectorRef,
    Component,
    OnDestroy,
    OnInit,
    computed,
    inject,
} from "@angular/core";
import {toSignal} from "@angular/core/rxjs-interop";
import {
    Router,
    RouterLink,
    RouterLinkActive,
    RouterOutlet,
} from "@angular/router";
import {
    MatSidenav,
    MatSidenavContainer,
    MatSidenavContent,
} from "@angular/material/sidenav";
import {MediaMatcher} from "@angular/cdk/layout";
import {MatIconButton} from "@angular/material/button";
import {MatIcon} from "@angular/material/icon";
import {MatMenu, MatMenuItem, MatMenuTrigger} from "@angular/material/menu";
import {AuthState, SessionService} from "./services/session.service";
import {UpdateService} from "./services/update.service";
import {NgClass} from "@angular/common";
import {environment} from "../environments/environment.dev";
import {gitCount} from "./build";

@Component({
    selector: "app-root",
    imports: [
        RouterOutlet,
        MatSidenavContainer,
        MatSidenav,
        MatSidenavContent,
        MatIconButton,
        MatIcon,
        MatMenuTrigger,
        MatMenu,
        MatMenuItem,
        RouterLink,
        RouterLinkActive,
        NgClass,
    ],
    templateUrl: "./app.html",
    styleUrl: "./app.scss",
})
export class App implements OnInit, OnDestroy {
    private readonly sessionSvc = inject(SessionService);
    private readonly authData = toSignal(
        this.sessionSvc.getAuthDataObservable(),
        {requireSync: true},
    );

    public readonly authState = computed(() => this.authData().state);
    public readonly points = computed(() => {
        const ad = this.authData();
        return ad.state === AuthState.SignedIn ? ad.sessionData.points : 0;
    });

    mobileQuery: MediaQueryList;
    darkModeQuery: MediaQueryList;
    public isDark = false;

    private readonly _mobileQueryListener: () => void;
    private readonly _darkModeQueryListener: () => void;
    private readonly _visibilityListener: () => void;

    constructor(
        changeDetectorRef: ChangeDetectorRef,
        media: MediaMatcher,
        private readonly router: Router,
        private readonly updateSvc: UpdateService,
    ) {
        this.mobileQuery = media.matchMedia("(max-width: 600px)");
        this._mobileQueryListener = () => changeDetectorRef.detectChanges();
        this.mobileQuery.addEventListener("change", this._mobileQueryListener);

        this.darkModeQuery = media.matchMedia("(prefers-color-scheme: dark)");
        this._darkModeQueryListener = () => {
            this.isDark = this.darkModeQuery.matches;
            changeDetectorRef.detectChanges();
        };
        this.darkModeQuery.addEventListener(
            "change",
            this._darkModeQueryListener,
        );

        this._visibilityListener = this.visibilityListener.bind(this);
    }

    public signOut(): void {
        this.sessionSvc
            .signOut()
            .catch(null)
            .finally(() => {
                this.router.navigate([""]);
            });
    }

    private visibilityListener(): void {
        this.updateSvc.checkForUpdate().catch(console.error);
    }

    ngOnInit() {
        document.addEventListener("visibilitychange", this._visibilityListener);
        window.addEventListener("focus", this._visibilityListener);
        this.isDark = this.darkModeQuery.matches;
    }

    ngOnDestroy(): void {
        document.removeEventListener(
            "visibilitychange",
            this._visibilityListener,
        );
        window.removeEventListener("focus", this._visibilityListener);
        this.mobileQuery.removeEventListener(
            "change",
            this._mobileQueryListener,
        );
        this.darkModeQuery.removeEventListener(
            "change",
            this._darkModeQueryListener,
        );
    }

    protected readonly AuthState = AuthState;
    protected readonly Math = Math;
    protected readonly environment = environment;
    protected readonly gitCount = gitCount;
}
