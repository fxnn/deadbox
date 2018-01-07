import React, {Component} from 'react';

class EncryptionKeyBox extends Component {

    constructor(props) {
        super(props);
        this.state = {
            keyAvailable: false
        };
    }

    componentDidMount() {
        const self = this;
        self.setState({keyAvailable: true});
    }

    render() {
        if(this.state.keyAvailable) {
            return (
                <div>
                    <nav className="navbar is-black">
                        <div className="navbar-brand">
                            <h1 className="is-size-4">
                                {this.props.logo}
                            </h1>
                        </div>
                        <div className="navbar-menu">
                            <div className="navbar-start" />
                            <div className="navbar-end">
                                <div className="navbar-item tags has-addons">
                                    <span className="tag">personal key</span>
                                    <span className="tag is-success">AA:BB:CC:DD:EE:FF</span>
                                </div>
                            </div>
                        </div>
                    </nav>
                    {this.props.children}
                </div>
            );
        }

        return (
            <section className="hero is-success is-fullheight">
                <div className="hero-body">
                    <div className="container has-text-centered">
                        <h1 className="title">
                            {this.props.logo}
                        </h1>
                        <p>Your key is generated, please wait ...</p>
                    </div>
                </div>
            </section>
        );
    }

}

export default EncryptionKeyBox;