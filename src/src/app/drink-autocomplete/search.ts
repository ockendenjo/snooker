import {Beer} from "../services/beers.service";

export function searchBeers(allBeers: Beer[], search: string): Beer[] {
    if (!search || search.length < 3) {
        return [];
    }
    const lower = search.toLowerCase().trim();
    const parts = lower.split(" ");

    const filtered = allBeers.filter((b) => {
        return parts.every((p) => {
            return (
                b.name.toLowerCase().includes(p) ||
                b.brewery.toLowerCase().includes(p)
            );
        });
    });
    if (filtered.length > 10) {
        return filtered.slice(0, 10);
    }
    return filtered;
}
