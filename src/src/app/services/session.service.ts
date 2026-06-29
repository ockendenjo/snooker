import {Injectable} from "@angular/core";
import {ApiService} from "./api.service";
import {AuthService} from "./auth.service";
import {BehaviorSubject, Observable} from "rxjs";

@Injectable({
    providedIn: "root",
})
export class SessionService {
    private readonly authDataSubject = new BehaviorSubject<AuthData>({
        state: AuthState.Unknown,
    });

    private readonly authData = this.authDataSubject.asObservable();

    public getAuthDataObservable(): Observable<AuthData> {
        return this.authData;
    }

    private sessionDataInProg = false;
    private readonly progressCallbacks: Set<(ad: AuthData) => void> = new Set();

    constructor(
        private readonly apiService: ApiService,
        private readonly authService: AuthService,
    ) {
        this.authService.isSignedIn.subscribe((signedIn) => {
            if (
                !signedIn &&
                this.authDataSubject.getValue().state === AuthState.SignedIn
            ) {
                this.authDataSubject.next({state: AuthState.NoAuth});
            }
        });

        this.getAuthData();
    }

    public getSessionData(refresh = false): Promise<SessionData> {
        return this.getAuthData(refresh).then((ad) => {
            if (ad.state == AuthState.SignedIn) {
                return ad.sessionData;
            }
            throw new Error("not signed in");
        });
    }

    private getAuthData(refresh = false): Promise<AuthData> {
        const current = this.authDataSubject.getValue();
        if (current.state !== AuthState.Unknown && !refresh) {
            return Promise.resolve(current);
        }

        if (this.sessionDataInProg) {
            return new Promise((r) => {
                this.progressCallbacks.add(r);
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
                this.authDataSubject.next(ad);
                for (const f of this.progressCallbacks.values()) {
                    f(ad);
                }
                this.progressCallbacks.clear();
                this.sessionDataInProg = false;
                return ad;
            });
    }

    public isSignedIn(): Promise<AuthState> {
        const current = this.authDataSubject.getValue();
        if (current.state !== AuthState.Unknown) {
            return Promise.resolve(current.state);
        }
        return this.getAuthData().then((ad) => ad.state);
    }

    public signOut(): Promise<void> {
        this.authService.logout();
        return Promise.resolve();
    }

    public setDisplayName(displayName: string): Promise<void> {
        return this.apiService
            .send("setDisplayName", {displayName})
            .then(() => {
                const current = this.authDataSubject.getValue();
                if (current.state === AuthState.SignedIn) {
                    this.authDataSubject.next({
                        ...current,
                        sessionData: {...current.sessionData, displayName},
                    });
                }
            });
    }

    public deleteAccount(): Promise<void> {
        return this.apiService.send("deleteAccount", {}).then(() => {
            this.authService.clearSession();
        });
    }
}

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
