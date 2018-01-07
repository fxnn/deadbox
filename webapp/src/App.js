import React, {Component} from 'react';
import EncryptionKeyBox from './EncryptionKeyBox';
import './App.css';

function AppLogo() {
    return (
        <div className="app-logo">
            deadbox
        </div>
    );
}

class App extends Component {
    render() {
        return (
            <EncryptionKeyBox logo={<AppLogo />}>
                <section className="hero is-primary is-fullheight">
                    <div className="hero-body">
                        <div className="container">
                            <h1 className="title">Connect</h1>
                        </div>
                    </div>
                </section>
            </EncryptionKeyBox>
        );
    }
}

export default App;
