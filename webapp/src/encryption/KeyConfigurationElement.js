import { div, h1, label, p } from "@hyperapp/html";
import { button, box, formField, formControl, formFieldLabel, formFieldBody, input } from "util/bulma";

const keyGenerationPasswordName = "keyGenerationPassword";

const blockElement = (title, children) => (
  box([
    h1({ class: "subtitle has-text-black" }, [title]),
    ...children
  ])
);

export const KeyConfigurationElement = ({ state, actions }) => (
  div([
    blockElement("Passphrase based key", [
      formField([
        p(["Delayed responses can be retrieved even after closing this session, using the same passphrase again."]),
      ]),
      formField({ class: "is-horizontal" }, [
        formFieldLabel({ class: "is-normal" }, [label({ class: "label" }, "Passphrase")]),
        formFieldBody([
          formField([
            formControl([
              input({ type: "password", name: keyGenerationPasswordName })
            ])
          ])
        ])
      ]),
      formField({ class: "is-horizontal" }, [
        formFieldLabel(),
        formFieldBody([
          formField([
            formControl([
              button({ class: "is-success", onclick: () => actions.setKeyConfigured() }, [
                "Generate key from passphrase"
              ])
            ])
          ])
        ])
      ])
    ]),
    blockElement("One-time key", [
      formField([
        p(["After closing this session and deleting all data associated therewith, responses become unreadable."]),
      ]),
      formField({ class: "is-horizontal" }, [
        formFieldLabel(),
        formFieldBody([
          formField([
            formControl([
              button({ class: "is-success", onclick: () => actions.setKeyConfigured() }, [
                "Generate one-time key"
              ])
            ])
          ])
        ])
      ])
    ])
  ])
);