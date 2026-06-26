import {Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {InfoPage} from "./info-page/info-page";
import {AbvPage} from "./abv-page/abv-page";
import {EdiMapPage} from "./edi-map-page/edi-map-page";

export const routes: Routes = [
    {
        path: "",
        component: HomePage,
    },
    {
        path: "abvs",
        component: AbvPage,
    },
    {
        path: "info",
        component: InfoPage,
    },
    {
        path: "map",
        component: EdiMapPage,
    },
];
