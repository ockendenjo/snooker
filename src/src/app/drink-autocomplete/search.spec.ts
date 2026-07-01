import {describe, it, expect} from "vitest";
import {searchBeers} from "./search";
import {Beer} from "../services/beers.service";

function makeBeer(name: string, brewery: string): Beer {
    return {
        name,
        brewery,
        style: "",
        abv: 0,
        untappd: 0,
    };
}

describe("searchBeers", () => {
    const beers: Beer[] = [
        makeBeer("Plum Porter", "Titanic"),
        makeBeer("Cherry Porter", "Titanic"),
        makeBeer("XPA", "Moonwake"),
    ];

    it("returns all beers when search is empty", () => {
        const got = searchBeers(beers, "");
        expect(got).toEqual([]);
    });

    it("matches by beer name (case-insensitive)", () => {
        const got = searchBeers(beers, "porter");
        expect(got).toEqual([beers[0], beers[1]]);
    });

    it("matches by brewery name (case-insensitive)", () => {
        const got = searchBeers(beers, "moonwake");
        expect(got).toEqual([beers[2]]);
    });

    it("returns empty array when no beers match", () => {
        const got = searchBeers(beers, "stout");
        expect(got).toEqual([]);
    });

    it("returns empty array when input list is empty", () => {
        const got = searchBeers([], "IPA");
        expect(got).toEqual([]);
    });

    it("matches name and brewery", () => {
        const got = searchBeers(beers, "moonwake XPA");
        expect(got).toEqual([beers[2]]);
    });

    it("matches partial brewery", () => {
        const got = searchBeers(beers, "moon");
        expect(got).toEqual([beers[2]]);
    });
});
