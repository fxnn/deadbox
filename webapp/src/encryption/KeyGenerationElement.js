import { h } from "hyperapp";
import { div } from "@hyperapp/html";
import { section } from "util/bulma";
import { Spinner } from "app/Spinner";

export const KeyGenerationElement = ({ state, actions }) => (
  div([
    section({ oncreate() { actions.setKeyAvailableDelayed(3000); } }, [h(Spinner)]),
    section([h("p", { class: "has-text-centered" }, "Your key is generated, please wait ...")]),
  ])
);
