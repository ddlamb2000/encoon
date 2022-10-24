
class Login extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
        };
    }

    render() {
        return (
            <br/>
        );
    }
}


const rootElement = document.getElementById("root")
const root = ReactDOM.createRoot(rootElement);
const db = rootElement.getAttribute("db");
root.render(<Login />);
