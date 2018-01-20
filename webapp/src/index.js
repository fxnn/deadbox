import { app } from "hyperapp"
import './index.css';
import {view} from './View';
import encryption from "encryption";

const state = {
  encryption: encryption.state
};

const actions = {
  encryption: encryption.actions
};

app(state, actions, view, document.body);
