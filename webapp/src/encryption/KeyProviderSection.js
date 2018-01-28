import { h } from "hyperapp";
import {div, h1, section} from "@hyperapp/html";
import {container} from "util/bulma";
import {KeyGenerationElement} from "./KeyGenerationElement";
import {KeyConfigurationElement} from "./KeyConfigurationElement";

export const KeyProviderSection = ({ logo, state, actions }) => (
  section({class:"hero is-success is-fullheight"}, [
    div({class:"hero-body"}, [
      container([
        h1({class:"title has-text-centered"}, [logo]),
        state.keyConfigurationAvailable
          ? h(KeyGenerationElement, { state: state, actions: actions })
          : h(KeyConfigurationElement, { state: state, actions: actions })
      ])
    ])
  ])
);
