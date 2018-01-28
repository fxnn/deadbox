import { h } from "hyperapp";
import { Spinner } from "app/Spinner";

export const KeyGenerationElement = ({ state, actions }) => (
  h("div", {},
    h(Spinner),
    h("p", {}, "Your key is generated, please wait ...")
  )
);
