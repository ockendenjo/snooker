import {Injectable} from "@angular/core";
import {BehaviorSubject} from "rxjs";
import {environment} from "../../environments/environment.dev";

interface CognitoTokenResponse {
    id_token: string;
    access_token: string;
    refresh_token: string;
    expires_in: number;
    token_type: string;
}

@Injectable({
    providedIn: "root",
})
export class AuthService {
    private readonly cognitoConfig = {
        domain: environment.cognito.domain,
        clientId: environment.cognito.clientId,
        redirectUri: `${window.location.origin}/callback`,
        scope: "openid email profile",
    };

    private readonly _isSignedInSubject = new BehaviorSubject<boolean>(false);
    public readonly isSignedIn = this._isSignedInSubject.asObservable();

    constructor() {
        this.isAuthenticated().then((res) => {
            this._isSignedInSubject.next(res);
        });
    }

    async getIdTokenAsync(
        thresholdMs: number = 60_000,
    ): Promise<string | null> {
        const idToken = this.getIdToken();
        const expiresAtStr = localStorage.getItem("cognito_expires_at");
        const refreshToken = this.getRefreshToken();

        if (!idToken || !expiresAtStr) {
            return null;
        }

        const expiresAt = parseInt(expiresAtStr, 10);
        const now = Date.now();

        if (now + thresholdMs < expiresAt) {
            return idToken;
        }

        if (refreshToken) {
            try {
                await this.refreshTokens(refreshToken);
                return this.getIdToken();
            } catch (e) {
                console.error("Failed to refresh tokens:", e);
                this.clearSession();
                return null;
            }
        }

        return idToken;
    }

    async login(): Promise<void> {
        const codeVerifier = this.generateCodeVerifier();
        const codeChallenge = await this.generateCodeChallenge(codeVerifier);

        localStorage.setItem("pkce_code_verifier", codeVerifier);

        const params = new URLSearchParams({
            response_type: "code",
            client_id: this.cognitoConfig.clientId,
            redirect_uri: this.cognitoConfig.redirectUri,
            scope: this.cognitoConfig.scope,
            code_challenge_method: "S256",
            code_challenge: codeChallenge,
        });

        window.location.href = `https://${this.cognitoConfig.domain}/oauth2/authorize?${params.toString()}`;
    }

    async handleCallback(code: string): Promise<void> {
        const codeVerifier = localStorage.getItem("pkce_code_verifier");
        if (!codeVerifier) {
            throw new Error("No code verifier found in session");
        }

        const tokens = await this.exchangeCodeForTokens(code, codeVerifier);
        this.storeTokens(tokens);
        localStorage.removeItem("pkce_code_verifier");
    }

    private async exchangeCodeForTokens(
        code: string,
        codeVerifier: string,
    ): Promise<CognitoTokenResponse> {
        const params = new URLSearchParams({
            grant_type: "authorization_code",
            client_id: this.cognitoConfig.clientId,
            code: code,
            redirect_uri: this.cognitoConfig.redirectUri,
            code_verifier: codeVerifier,
        });

        const response = await fetch(
            `https://${this.cognitoConfig.domain}/oauth2/token`,
            {
                method: "POST",
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: params.toString(),
            },
        );

        if (!response.ok) {
            const error = await response.text();
            throw new Error(`Token exchange failed: ${error}`);
        }

        return response.json();
    }

    private storeTokens(tokens: CognitoTokenResponse): void {
        localStorage.setItem("cognito_id_token", tokens.id_token);
        localStorage.setItem("cognito_access_token", tokens.access_token);
        localStorage.setItem("cognito_refresh_token", tokens.refresh_token);
        localStorage.setItem("cognito_token_type", tokens.token_type);

        const expiresAt = Date.now() + tokens.expires_in * 1000;
        localStorage.setItem("cognito_expires_at", expiresAt.toString());

        this._isSignedInSubject.next(true);
    }

    getIdToken(): string | null {
        return localStorage.getItem("cognito_id_token");
    }

    async isAuthenticated(): Promise<boolean> {
        try {
            const token = await this.getIdTokenAsync();
            if (!token) return false;

            const expiresAt = localStorage.getItem("cognito_expires_at");
            if (!expiresAt) return false;

            return Date.now() < parseInt(expiresAt, 10);
        } catch {
            return false;
        }
    }

    logout(): void {
        this.clearSession();

        const params = new URLSearchParams({
            client_id: this.cognitoConfig.clientId,
            logout_uri: window.location.origin,
        });

        window.location.href = `https://${this.cognitoConfig.domain}/logout?${params.toString()}`;
    }

    public clearSession(): void {
        const keys = [
            "cognito_id_token",
            "cognito_access_token",
            "cognito_refresh_token",
            "cognito_token_type",
            "cognito_expires_at",
        ];

        keys.forEach((key) => {
            localStorage.removeItem(key);
            sessionStorage.removeItem(key);
        });

        this._isSignedInSubject.next(false);
    }

    private generateCodeVerifier(): string {
        const array = new Uint8Array(32);
        crypto.getRandomValues(array);
        return this.base64UrlEncode(array);
    }

    private async generateCodeChallenge(verifier: string): Promise<string> {
        const encoder = new TextEncoder();
        const data = encoder.encode(verifier);
        const hash = await crypto.subtle.digest("SHA-256", data);
        return this.base64UrlEncode(new Uint8Array(hash));
    }

    private base64UrlEncode(array: Uint8Array): string {
        const base64 = btoa(String.fromCharCode(...array));
        return base64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
    }

    private getRefreshToken(): string | null {
        return localStorage.getItem("cognito_refresh_token");
    }

    private async refreshTokens(refreshToken: string): Promise<void> {
        const params = new URLSearchParams({
            grant_type: "refresh_token",
            client_id: this.cognitoConfig.clientId,
            refresh_token: refreshToken,
        });

        const response = await fetch(
            `https://${this.cognitoConfig.domain}/oauth2/token`,
            {
                method: "POST",
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: params.toString(),
            },
        );

        if (!response.ok) {
            const error = await response.text();
            throw new Error(`Refresh token request failed: ${error}`);
        }

        const data: Partial<CognitoTokenResponse> = await response.json();
        const merged: CognitoTokenResponse = {
            id_token: data.id_token as string,
            access_token: data.access_token as string,
            refresh_token: (data.refresh_token as string) ?? refreshToken,
            expires_in: data.expires_in as number,
            token_type: data.token_type as string,
        };
        this.storeTokens(merged);
    }
}
