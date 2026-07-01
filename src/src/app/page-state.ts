export type PageState<T> = StateLoading | StateLoaded<T> | StateError;

export interface StateLoading {
    state: "LOADING";
}

export interface StateLoaded<T> {
    state: "LOADED";
    data: T;
}

export interface StateError {
    state: "ERROR";
    error: Error;
}
