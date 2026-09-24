export type PageState<T> = StateLoading | StateLoaded<T> | StateError;

interface StateLoading {
    state: "LOADING";
}

interface StateLoaded<T> {
    state: "LOADED";
    data: T;
}

interface StateError {
    state: "ERROR";
    error: Error;
}
