import { ul, li, a, span, input, label } from "@hyperapp/html";
import { card, cardContent, container, tabs, formField, formControl } from "util/bulma";

const keyGenerationModeName = "keyGenerationMode";
const keyGenerationModeWithPassphraseValue = "withPassphrase";
const keyGenerationModeOneTimeValue = "oneTime";

const keyGenerationPasswordName = "keyGenerationPassword";

export const KeyConfigurationElement = ({ state, actions }) => (
  tabs({ class: "is-boxed" }, [
    ul([
      li({ class: "is-active" }, [
        a([
          span([
            "Passphrase based-key"
          ])
        ])
      ]),
      li([
        a([
          span([
            "One-time key"
          ])
        ])
      ])
    ])
  ])
);

/*
formField([
            formControl([
              label({ class: "radio" }, [
                input({ type: "radio", name: keyGenerationModeName, value: keyGenerationModeWithPassphraseValue }),
                "I want to generate the key using a passphrase, so that I can receive responses later on."
              ])
            ])
          ]),
          formField([
            label({ class: "label" }, "Key Passphrase"),
            formControl([
              input({ type: "password", name: keyGenerationPasswordName })
            ])
          ])

          formField([
          formControl([
            label({ class: "radio" }, [
              input({ type: "radio", name: keyGenerationModeName, value: keyGenerationModeOneTimeValue }),
              "I do only want to generate a one-time key."
            ])
          ])
        ])
 */