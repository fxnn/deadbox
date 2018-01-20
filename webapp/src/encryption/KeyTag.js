import { h } from "hyperapp"

export const KeyTag = () => (
  h("div", {class:"tags has-addons"}, [
    h("span", {class:"tag"}, "personal key"),
    h("span", {class:"tag is-success"}, "AA:BB:CC:DD:EE:FF")
  ])
);