import { h } from "hyperapp"
import {KeyProviderSection} from "./KeyProviderSection";

export const WithEncryptionKey = ({actions, state, whenKeyAvailable, logo}) => (
  h("div", {
    oncreate() { /*actions.setKeyAvailableDelayed(3000);*/ } // HINT: Mock key creation for now
  }, state.keyAvailable ? whenKeyAvailable : h(KeyProviderSection, {logo: logo, state: state, actions: actions}))
);
