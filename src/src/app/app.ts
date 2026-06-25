import {ChangeDetectorRef, Component, OnDestroy, OnInit} from "@angular/core";
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
import {AuthData, AuthState, SessionService} from "./services/session.service";
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
    public authState: AuthState = AuthState.Unknown;
    private symbol?: symbol;

    mobileQuery: MediaQueryList;
    darkModeQuery: MediaQueryList;

    public points: number = 0;
    public isDark = false;

    private readonly _mobileQueryListener: () => void;
    private readonly _darkModeQueryListener: () => void;
    private readonly _visibilityListener: () => void;

    constructor(
        changeDetectorRef: ChangeDetectorRef,
        media: MediaMatcher,
        private readonly sessionSvc: SessionService,
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

        this.symbol = this.sessionSvc.registerChangeCallback((ad) => {
            this.processAuthData(ad);
        });

        this.sessionSvc.getAuthData().then((ad: AuthData) => {
            this.processAuthData(ad);
        });

        this.isDark = this.darkModeQuery.matches;
    }

    private processAuthData(ad: AuthData) {
        this.authState = ad.state;
        if (ad.state == AuthState.SignedIn) {
            this.points = ad.sessionData.points;
        } else {
            this.points = 0;
        }
    }

    ngOnDestroy(): void {
        if (this.symbol) {
            this.sessionSvc.deregisterChangeCallback(this.symbol);
        }
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
