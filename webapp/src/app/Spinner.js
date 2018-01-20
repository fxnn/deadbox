import {h} from "hyperapp"
import "./Spinner.css"

export const Spinner = () => (
  h("div", {class:"sk-cube-grid"}, [...Array(9)].map((_,i)=>h("div", {class:"sk-cube sk-cube" + (i+1)})))
);
