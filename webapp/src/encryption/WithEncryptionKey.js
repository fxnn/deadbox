import { h } from "hyperapp"
import {KeyGeneratingSection} from "./KeyGeneratingSection";

export const WithEncryptionKey = ({actions, state, whenKeyAvailable, logo}) => (
  h("div", {
    oncreate() { actions.setKeyAvailableDelayed(3000); } // HINT: Mock key creation for now
  }, state.keyAvailable ? whenKeyAvailable : h(KeyGeneratingSection, {logo: logo}))
);
