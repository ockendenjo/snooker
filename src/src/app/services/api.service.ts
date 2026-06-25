import {Injectable} from "@angular/core";
import {AuthService} from "./auth.service";

@Injectable({
    providedIn: "root",
})
export class ApiService {
    constructor(private readonly authService: AuthService) {}

    private async authHeaders(): Promise<Record<string, string>> {
        const token = await this.authService.getIdTokenAsync();
        return token ? {Authorization: `Bearer ${token}`} : {};
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    public async send(path: string, bodyObj: any): Promise<any> {
        const headers = await this.authHeaders();
        const r = await fetch(`/api/${path}`, {
            method: "POST",
            headers,
            body: JSON.stringify(bodyObj),
        });
        if (r.status > 299) {
            const t = await r.text();
            return Promise.reject(t);
        }
        if (r.status == 204) {
            return Promise.resolve();
        }
        return r.json();
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    public async get(path: string): Promise<any> {
        const headers = await this.authHeaders();
        const r = await fetch(`/api/${path}`, {method: "GET", headers});
        if (r.status > 299) {
            return Promise.reject(r);
        }
        if (r.status == 204) {
            return Promise.resolve();
        }
        return r.json();
    }
}
