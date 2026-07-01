import {Injectable} from "@angular/core";

@Injectable({
    providedIn: "root",
})
export class BeersService {
    private mapBeerFileToBeers(bf: BeerFile): Beer[] {
        return bf
            .map((brewery) => {
                return brewery.beers.map((b): Beer => {
                    return {
                        ...b,
                        brewery: brewery.name,
                        breweryImage: brewery.image,
                    };
                });
            })
            .flatMap((i) => i);
    }

    public async loadAll(): Promise<Beer[]> {
        const idxFile = (await fetch("beer/index.json").then((r) =>
            r.json(),
        )) as IndexFile;

        const promises: Promise<BeerFile>[] = idxFile.files.map((f) => {
            return fetch(`beer/${f}`)
                .then((r) => r.json())
                .catch(() => {
                    return [];
                });
        });

        return await Promise.all(promises).then((beerFiles) => {
            return beerFiles.map(this.mapBeerFileToBeers).flatMap((i) => i);
        });
    }
}

interface IndexFile {
    files: string[];
}

type BeerFile = Brewery[];

interface Brewery {
    name: string;
    image?: string;
    beers: BeerDetail[];
}

interface BeerDetail {
    name: string;
    abv: number;
    untappd: number;
    image?: string;
    style: string;
}

export type Beer = BeerDetail & {
    brewery: string;
    breweryImage?: string;
};
