import { h } from "hyperapp";

/**
 * Like Object.assign(), but concats props that are in both objects
 * @param target
 * @param source
 * @return modified target object
 */
function mergeProps(target, source) {
  if (source) {
    for (const propName in source) {
      if (source.hasOwnProperty(propName)) {
        if (target[propName]) {
          target[propName] = source[propName] + " " + target[propName]
        } else {
          target[propName] = source[propName]
        }
      }
    }
  }
  return target;
}

// based on https://github.com/hyperapp/html/blob/master/src/html.js#L3
function vnode(name, additionalProps) {
  return function (props, children) {
    return typeof props === "object" && !Array.isArray(props)
      ? h(name, mergeProps(props,  additionalProps), children)
      : h(name, additionalProps, props)
  }
}

export function container(props, children) { return vnode("div", { class: "container" })(props, children); }

export function card(props, children) { return vnode("div", { class: "card" })(props, children); }
export function cardContent(props, children) { return vnode("div", { class: "card-content" })(props, children); }

export function tabs(props, children) { return vnode("div", { class: "tabs" })(props, children); }

export function formField(props, children) { return vnode("div", { class: "field" })(props, children); }
export function formControl(props, children) { return vnode("div", { class: "control" })(props, children); }
