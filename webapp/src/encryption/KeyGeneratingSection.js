import { h } from "hyperapp"

export const KeyGeneratingSection = ({ logo }) => (
  h("section", {class:"hero is-success is-fullheight"}, [
    h("div", {class:"hero-body"}, [
      h("div", {class:"container has-text-centered"}, [
        h("h1", {class:"title"}, [
          logo
        ]),
        h("p", {}, "Your key is generated, please wait ...")
      ])
    ])
  ])
);