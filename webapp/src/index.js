import { app } from "hyperapp"
import './index.css';
import {state, actions, view} from './App';

app(state, actions, view, document.body);
