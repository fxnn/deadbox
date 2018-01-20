import { h } from "hyperapp"
import {Logo} from "app/Logo";
import {WithEncryptionKey} from "encryption/WithEncryptionKey";
import {Nav} from "app/Nav";
import {KeyTag} from "encryption/KeyTag";

export const view = (state, actions) => {
  return h(WithEncryptionKey, {
    actions: actions.encryption,
    state: state.encryption,
    whenKeyAvailable: h(Nav, {logo: h(Logo), items: [h(KeyTag)]}),
    logo: h(Logo)
  });
};
