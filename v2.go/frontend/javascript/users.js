
class Users extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: [],
        };
    }

    render() {
        const { items, isLoaded, error } = this.state;

        if (error) {
            return <div>Erreur : {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Chargement…</div>;
        } else {
            if(uuid !== "") {
                return (
                    <div>
                        <table className="table table-hover table-sm">
                            <thead className="table-light">
                                <tr>
                                    <th scope="col">Filed</th>
                                    <th scope="col">Value</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr><td>Uuid</td><td>{items.uuid}</td></tr>
                                <tr><td>Email</td><td>{items.email}</td></tr>
                                <tr><td>First Name</td><td>{items.firstName}</td></tr>
                                <tr><td>Last Name</td><td>{items.lastName}</td></tr>
                            </tbody>
                        </table>
                    </div>
                );
            }
            else {
                return (
                    <div>
                        <table className="table table-hover table-sm">
                            <thead className="table-light">
                                <tr>
                                    <th scope="col">Email</th>
                                    <th scope="col">First Name</th>
                                    <th scope="col">Last Name</th>
                                    <th scope="col">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-battery-full" viewBox="0 0 16 16">
                                            <path d="M2 6h10v4H2V6z"/>
                                            <path d="M2 4a2 2 0 0 0-2 2v4a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2H2zm10 1a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h10zm4 3a1.5 1.5 0 0 1-1.5 1.5v-3A1.5 1.5 0 0 1 16 8z"/>
                                        </svg>
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                {items.map(item => (
                                    <tr key={item.uuid}>
                                        <td>{item.email}</td>
                                        <td>{item.firstName}</td>
                                        <td>{item.lastName}</td>
                                        <th scope="row"><a href={`/${db}/user/${item.uuid}`}>{item.uuid}</a></th>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                );
            }
        }
    }

    componentDidMount() {
        if(uuid !== "") {
            fetch(`/${db}/api/v1/user/${uuid}`)
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        items: result.user
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
        }
        else {
            fetch(`/${db}/api/v1/users`)
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        items: result.users
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
        }
    }
}


const rootElement = document.getElementById("root")
const root = ReactDOM.createRoot(rootElement);
const db = rootElement.getAttribute("db");
const uuid = rootElement.getAttribute("uuid");
root.render(<Users />);
