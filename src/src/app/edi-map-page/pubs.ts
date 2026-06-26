import VectorSource from "ol/source/Vector";
import {Feature} from "ol";
import {Point} from "ol/geom";
import {fromLonLat} from "ol/proj";
import {Fill, RegularShape, Stroke, Style} from "ol/style";
import {Pub} from "../services/pubs.service";

export function addPubsToSource(
    selectedSource: VectorSource,
    pubs: Pub[],
): void {
    const filtered = pubs.filter((p) => p.hasRealAle);

    filtered.forEach((p, idx) => {
        const iconFeature = new Feature({
            geometry: new Point(fromLonLat([p.lon, p.lat])),
            pub: p,
        });

        iconFeature.setStyle(getStyle());
        selectedSource.addFeature(iconFeature);
    });
}

function getStyle() {
    return new Style({
        image: new RegularShape({
            points: 4,
            radius: 20,
            angle: Math.PI / 4,
            fill: new Fill({color: "red"}),
            stroke: new Stroke({color: "white", width: 1}),
        }),
    });
}
