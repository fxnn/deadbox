import { h } from "hyperapp"
import { Spinner } from "app/Spinner";

export const KeyGeneratingSection = ({ logo }) => (
  h("section", {class:"hero is-success is-fullheight"}, [
    h("div", {class:"hero-body"}, [
      h("div", {class:"container has-text-centered"}, [
        h("h1", {class:"title"}, [
          logo
        ]),
        h(Spinner),
        h("p", {}, "Your key is generated, please wait ...")
      ])
    ])
  ])
);