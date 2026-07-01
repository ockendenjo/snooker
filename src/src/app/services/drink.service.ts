import {Injectable} from "@angular/core";
import {ApiService} from "./api.service";

@Injectable({
    providedIn: "any",
})
export class DrinkService {
    constructor(private readonly apiService: ApiService) {}

    public logDrink(drink: NewDrink): Promise<LoggedDrink> {
        return this.apiService.send("logDrink", drink);
    }

    public getDrinks(): Promise<Drink[]> {
        return this.apiService.get("users/me/drinks").then((j: DrinksFile) => {
            return j.drinks;
        });
    }

    public listUnknownDrinks(): Promise<Drink[]> {
        return this.apiService.get("unknownDrinks").then((j: DrinksFile) => {
            return j.drinks;
        });
    }

    public getDrinksForUser(userID: string): Promise<Drink[]> {
        return this.apiService
            .get(`users/${userID}/drinks`)
            .then((j: DrinksFile) => {
                return j.drinks;
            });
    }

    public updateDrink(
        drink: Drink,
        endOfWord?: boolean,
        notInWord?: boolean,
    ): Promise<void> {
        const body: Record<string, unknown> = {
            userID: drink.userID,
            timestamp: drink.timestamp,
        };
        if (endOfWord !== undefined) body["endOfWord"] = endOfWord;
        if (notInWord !== undefined) body["notInWord"] = notInWord;
        return this.apiService.send("updateDrink", body);
    }

    public updateDrinkBeer(
        drink: Drink,
        name: string,
        brewery: string,
        untappdID: number,
    ): Promise<void> {
        const body: Record<string, unknown> = {
            userID: drink.userID,
            timestamp: drink.timestamp,
            name,
            brewery,
            untappdID,
        };
        return this.apiService.send("updateDrink", body);
    }
}

export type NewDrink = {
    timestamp?: string;
    pubID: number;
    name: string;
    brewery: string;
    untappdID?: number;
    with: string;
    endOfWord: boolean;
    notInWord: boolean;
    notes?: string;
};

export type Drink = otherDrink | endOfWordDrink;

interface baseDrink {
    userID: string;
    timestamp: string;
    pubID: number;
    name: string;
    brewery: string;
    untappdID?: number;
    points: number;
    notInWord: boolean;
    with: string;
    letter: string;
    actualPoints: number;
    tooShort?: boolean;
    notes?: string;
}

type otherDrink = baseDrink & {
    endOfWord: false;
};

type endOfWordDrink = baseDrink & {
    endOfWord: true;
    multiplier: number;
    sumLetters: number;
    wordLength: number;
    wordPoints: number;
};

export type DrinksFile = {
    drinks: Drink[];
};

export interface LoggedDrink {
    cid: string;
}
