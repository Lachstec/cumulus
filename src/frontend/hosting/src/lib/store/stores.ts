import { writable, derived, type Writable } from "svelte/store";

export const isAuthenticated = writable(false);
export const user: Writable<{} | undefined> = writable({});
export const popupOpen = writable(false);
export const error = writable();

