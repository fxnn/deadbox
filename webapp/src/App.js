import React, {Component} from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {
    render() {
        return (
            <div className="App">
                <header className="hero is-primary is-bold">
                    <div className="hero-body">
                        <div className="container">
                            <div className="columns">
                                <div className="column">
                                    <h1 className="title">my app</h1>
                                    <h2 className="subtitle">built with react</h2>
                                </div>
                                <div className="column is-narrow">
                                    <img src={logo} className="App-logo" alt="logo"/>
                                </div>
                            </div>
                        </div>
                    </div>
                </header>
                <section className="section">
                    <div className="container">
                        <h1 className="title">Section</h1>
                        To get started, edit <code>src/App.js</code> and save to reload.
                    </div>
                </section>
            </div>
        );
    }
}

export default App;
