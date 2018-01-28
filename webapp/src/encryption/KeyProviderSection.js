import { h } from "hyperapp";
import {div, h1} from "@hyperapp/html";
import {container, section} from "util/bulma";
import {KeyGenerationElement} from "./KeyGenerationElement";
import {KeyConfigurationElement} from "./KeyConfigurationElement";

export const KeyProviderSection = ({ logo, state, actions }) => (
  div({class:"hero is-primary is-fullheight"}, [
    div({class:"hero-body"}, [
      container([
        section([
          h1({class:"title has-text-centered"}, [logo])
        ]),
        state.keyConfigurationAvailable
          ? h(KeyGenerationElement, { state: state, actions: actions })
          : h(KeyConfigurationElement, { state: state, actions: actions })
      ])
    ])
  ])
);
