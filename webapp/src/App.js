import { h } from "hyperapp"
import './App.css';

const AppLogo = () => (
    h("div", {class: 'app-logo'}, "deadbox")
);

export const state = {};
export const actions = {};
export const view = (state, actions) => (
    h(AppLogo)
);
