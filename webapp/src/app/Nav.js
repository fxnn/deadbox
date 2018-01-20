import { h } from "hyperapp"

export const Nav = ({ logo, items }, children) => (
  h("div", {}, [
    h("nav", {class: "navbar is-black"}, [
      h("div", {class:"navbar-brand"}, [
        h("h1", {class:"is-size-5 is-vertical-center"}, [
          logo
        ])
      ]),
      h("div", {class:"navbar-menu"}, [
        h("div", {class:"navbar-start"}),
        h("div", {class:"navbar-end"},
          items.map(item => h("div", {class:"navbar-item"}, [item]))
        )
      ])
    ]),
    ...children
  ])
);