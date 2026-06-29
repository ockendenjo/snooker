import {Injectable} from "@angular/core";
import {ApiService} from "./api.service";

@Injectable({
    providedIn: "root",
})
export class PubsService {
    private pubs: Pub[] = [];

    constructor(private readonly apiSvc: ApiService) {}

    public getPubs(): Promise<Pub[]> {
        if (this.pubs.length > 0) {
            return Promise.resolve(this.pubs);
        }

        return fetch("pubs.json")
            .then((r) => r.json())
            .then((j: PubsFile) => {
                this.pubs = j.pubs;
                return j.pubs;
            });
    }
}

export interface PubsFile {
    pubs: Pub[];
}

export interface Pub {
    goodBeerID: number;
    lat: number;
    lon: number;
    name: string;
    address: string;
    realAles: number;
    numBeers: number;
    hasRealAle: boolean;
    camraID?: number;
    tempClosed?: boolean;
    chain?: string;
}
