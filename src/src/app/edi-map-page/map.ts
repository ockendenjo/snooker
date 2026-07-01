import {Map as OLMap, Overlay, View} from "ol";
import TileLayer from "ol/layer/Tile";
import OSM from "ol/source/OSM";
import VectorLayer from "ol/layer/Vector";
import VectorSource from "ol/source/Vector";
import {defaults as defaultControls, ScaleLine} from "ol/control";
import {addPubsToSource} from "./pubs";
import {Extent} from "ol/extent";
import {Pub} from "../services/pubs.service";

export function renderMap(pubs: Pub[], navigateToLog: (pubID: number) => void) {
    const osmLayer = new TileLayer({
        properties: {
            title: "OSM",
            displayInLayerSwitcher: false,
        },
        source: new OSM(),
        opacity: 0.5,
    });

    const nearSource = new VectorSource({wrapX: false});
    const nearLayer = new VectorLayer({
        source: nearSource,
        visible: true,
    });

    const farSource = new VectorSource({wrapX: false});
    const farLayer = new VectorLayer({
        source: farSource,
        visible: true,
    });

    const mapView = new View({maxZoom: 19});
    mapView.setZoom(10);
    let map: OLMap;

    function initialiseMap() {
        map = new OLMap({
            controls: defaultControls().extend([new ScaleLine()]),
            target: "map",
            layers: [osmLayer, farLayer, nearLayer],
            keyboardEventTarget: document,
            view: mapView,
        });

        const container = document.getElementById("popup") as HTMLElement;
        const content = document.getElementById("popup-content") as HTMLElement;
        const closer = document.getElementById("popup-closer") as HTMLElement;

        const overlay = new Overlay({
            element: container,
            autoPan: true,
        });

        map.on("click", function (e) {
            const feature = map.forEachFeatureAtPixel(e.pixel, (f) =>
                f.getProperties()["pub"] ? f : false,
            );
            if (!feature) {
                overlay.setPosition(undefined);
                closer.blur();
                return;
            }

            const properties = feature.getProperties();
            const pub = properties["pub"] as Pub;

            content.innerHTML = getPopupHTML(pub);

            const logBtn = document.createElement("button");
            logBtn.textContent = "Log a drink here";
            logBtn.style = "margin-top: 16px;";
            logBtn.onclick = () => navigateToLog(pub.camraID);
            content.appendChild(logBtn);

            overlay.setPosition(e.coordinate);
        });

        map.addOverlay(overlay);
        closer.onclick = () => {
            overlay.setPosition(undefined);
            closer.blur();
            return false;
        };
    }

    initialiseMap();

    addPubsToSource(nearSource, farSource, pubs);

    mapView.fit(nearSource.getExtent() as Extent, {
        padding: [20, 20, 20, 20],
    });
}

function getPopupHTML(pub: Pub): string {
    let html = `<div>${pub.name}</div>`;

    if (pub.hasRealAle) {
        html += `<div>We think this pub serves ${pub.realAles} real ale${pub.realAles == 1 ? "" : "s"}</div>`;
    } else {
        html += `<div class="closed">Real ale not available</div>`;
    }

    html += `<div><a href="https://camra.org.uk/pubs/${pub.camraID}" target="_blank">CAMRA listing for pub</a></div>`;

    return html;
}
