import {Injectable} from "@angular/core";
import {ApiService} from "./api.service";
import {AuthService} from "./auth.service";

@Injectable({
    providedIn: "root",
})
export class SessionService {
    private authData?: AuthData;
    private sessionDataInProg = false;

    private changeCallbacks: Map<symbol, CallbackFn> = new Map();
    private progessCallbacks: Set<CallbackFn> = new Set<CallbackFn>();

    constructor(
        private readonly apiService: ApiService,
        private readonly authService: AuthService,
    ) {
        this.authService.isSignedIn.subscribe((signedIn) => {
            if (!signedIn && this.authData?.state === AuthState.SignedIn) {
                this.setAuthState({state: AuthState.NoAuth});
            }
        });
    }

    private setAuthState(ad: AuthData): void {
        this.authData = ad;
        for (const f of this.changeCallbacks.values()) {
            f(ad);
        }
    }

    public getSessionData(refresh = false): Promise<SessionData> {
        return this.getAuthData(refresh).then((ad) => {
            if (ad.state == AuthState.SignedIn) {
                return ad.sessionData;
            }
            throw new Error("not signed in");
        });
    }

    public getAuthData(refresh = false): Promise<AuthData> {
        if (this.authData && !refresh) {
            return Promise.resolve(this.authData);
        }

        if (this.sessionDataInProg) {
            return new Promise((r) => {
                this.progessCallbacks.add(r);
            });
        }

        this.sessionDataInProg = true;

        return this.authService
            .isAuthenticated()
            .then((isAuth) => {
                if (!isAuth) {
                    return {state: AuthState.NoAuth} as AuthData;
                }
                return this.apiService
                    .send("getSessionData", {})
                    .then(
                        (sd: SessionData) =>
                            ({
                                state: AuthState.SignedIn,
                                sessionData: sd,
                            }) as AuthData,
                    )
                    .catch(() => ({state: AuthState.NoAuth}) as AuthData);
            })
            .then((ad) => {
                this.setAuthState(ad);
                for (const f of this.progessCallbacks.values()) {
                    f(ad);
                }
                this.progessCallbacks.clear();
                this.sessionDataInProg = false;
                return ad;
            });
    }

    public isSignedIn(): Promise<AuthState> {
        if (this.authData) {
            return Promise.resolve(this.authData.state);
        }
        return this.getAuthData().then((ad) => ad.state);
    }

    public signOut(): Promise<void> {
        this.authService.logout();
        return Promise.resolve();
    }

    /**
     * Register a listener to be notified of auth state changes
     * @param fn callback function
     */
    public registerChangeCallback(fn: CallbackFn): symbol {
        const s = Symbol();
        this.changeCallbacks.set(s, fn);
        return s;
    }

    public deregisterChangeCallback(s: symbol): void {
        this.changeCallbacks.delete(s);
    }

    public setDisplayName(displayName: string): Promise<void> {
        return this.apiService
            .send("setDisplayName", {displayName})
            .then(() => {
                if (this.authData?.state == AuthState.SignedIn) {
                    this.authData.sessionData.displayName = displayName;
                }
            });
    }

    public deleteAccount(): Promise<void> {
        return this.apiService.send("deleteAccount", {}).then(() => {
            this.authService.clearSession();
        });
    }
}

type CallbackFn = (authData: AuthData) => void;

export type SessionData = {
    ID: string;
    email: string;
    displayName: string;
    points: number;
    letters: string[];
};

export type AuthData = NoAuthAuthData | SignedInAuthData;

export type NoAuthAuthData = {
    state: AuthState.Unknown | AuthState.NoAuth;
};

export type SignedInAuthData = {
    state: AuthState.SignedIn;
    sessionData: SessionData;
};

export enum AuthState {
    Unknown = "Unknown",
    SignedIn = "SignedIn",
    NoAuth = "NoAuth",
}
