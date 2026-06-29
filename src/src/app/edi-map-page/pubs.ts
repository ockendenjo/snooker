import VectorSource from "ol/source/Vector";
import {Feature} from "ol";
import {Point} from "ol/geom";
import {fromLonLat} from "ol/proj";
import {Fill, Stroke, Style} from "ol/style";
import {Pub} from "../services/pubs.service";
import CircleStyle from "ol/style/Circle";

export function addPubsToSource(
    nearSource: VectorSource,
    farSource: VectorSource,
    pubs: Pub[],
): void {
    const filtered = pubs.filter((p) => {
        if (!p.hasRealAle || p.realAles < 1) {
            return false;
        }
        if (p.tempClosed) {
            return false;
        }
        if (p.chain === "Wetherspoons") {
            return false;
        }
        return true;
    });
    filtered.sort((a, b) => a.realAles - b.realAles);

    filtered.forEach((p, idx) => {
        const iconFeature = new Feature({
            geometry: new Point(fromLonLat([p.lon, p.lat])),
            pub: p,
        });

        iconFeature.setStyle(getStyle(p));
        if (isNear(p)) {
            nearSource.addFeature(iconFeature);
        } else {
            farSource.addFeature(iconFeature);
        }
    });
}

function isNear(p: Pub): boolean {
    const addr = p.address.toLowerCase();
    if (addr.includes("eh1 ")) {
        return true;
    }
    if (addr.includes("eh2 ")) {
        return true;
    }
    if (addr.includes("eh3 ")) {
        return true;
    }
    return false;
}

function getStyle(p: Pub) {
    return new Style({
        image: new CircleStyle({
            radius: 8,
            fill: new Fill({color: getFill(p.realAles)}),
            stroke: new Stroke({color: "#888", width: 1}),
        }),
    });
}

function getFill(numRealAles: number): string {
    switch (numRealAles) {
        case 0:
            return "#ffffe5";
        case 1:
            return "#feeaa1";
        case 2:
            return "#feba4a";
        case 3:
            return "#ee7918";
        case 4:
            return "#b74304";
        default:
            return "#662506";
    }
}
