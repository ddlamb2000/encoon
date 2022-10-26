// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class LoggedIn extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
          error: false,
          isLoaded: false,
          message: "",
          items: [],
      };
    }
  
    render() {
      const { items, isLoaded, error } = this.state;
  
      if(error) {
        alert(`Error for ${dbName}: ${this.state.message}`)
        this.logout()
      } else if (!isLoaded) {
          return <div>Loading…</div>;
      } else {
          if(uuid !== "") {
                    //Create Date Editor
                    var dateEditor = function(cell, onRendered, success, cancel){
                        //cell - the cell component for the editable cell
                        //onRendered - function to call when the editor has been rendered
                        //success - function to call to pass thesuccessfully updated value to Tabulator
                        //cancel - function to call to abort the edit and return to a normal cell

                        //create and style input
                        var cellValue = luxon.DateTime.fromFormat(cell.getValue(), "dd/MM/yyyy").toFormat("yyyy-MM-dd"),
                        input = document.createElement("input");

                        input.setAttribute("type", "date");

                        input.style.padding = "4px";
                        input.style.width = "100%";
                        input.style.boxSizing = "border-box";

                        input.value = cellValue;

                        onRendered(function(){
                            input.focus();
                            input.style.height = "100%";
                        });

                        function onChange(){
                            if(input.value != cellValue){
                                success(luxon.DateTime.fromFormat(input.value, "yyyy-MM-dd").toFormat("dd/MM/yyyy"));
                            }else{
                                cancel();
                            }
                        }

                        //submit new value on blur or change
                        input.addEventListener("blur", onChange);

                        //submit new value on enter
                        input.addEventListener("keydown", function(e){
                            if(e.keyCode == 13){
                                onChange();
                            }

                            if(e.keyCode == 27){
                                cancel();
                            }
                        });

                        return input;
                    };


                    //Build Tabulator
                    var table = new Tabulator("#example-table", {
                        height:"311px",
                        columns:[
                            {title:"Name", field:"name", width:150, editor:"input"},
                            {title:"Location", field:"location", width:130, editor:"list", editorParams:{autocomplete:"true", allowEmpty:true,listOnEmpty:true, valuesLookup:true}},
                            {title:"Progress", field:"progress", sorter:"number", hozAlign:"left", formatter:"progress", width:140, editor:true},
                            {title:"Gender", field:"gender", editor:"list", editorParams:{values:{"male":"Male", "female":"Female", "unknown":"Unknown"}}},
                            {title:"Rating", field:"rating",  formatter:"star", hozAlign:"center", width:100, editor:true},
                            {title:"Date Of Birth", field:"dob", hozAlign:"center", sorter:"date", width:140, editor:dateEditor},
                            {title:"Driver", field:"car", hozAlign:"center", editor:true, formatter:"tickCross"},
                        ],
                    });
              return (
                  <div>
                      <h2>User</h2>
                      <table className="table table-hover table-sm">
                          <thead className="table-light">
                          </thead>
                          <tbody>
                              <tr><td>Uuid</td><td>{items.uuid}</td></tr>
                              <tr><td>Email</td><td>{items.email}</td></tr>
                              <tr><td>First Name</td><td>{items.firstName}</td></tr>
                              <tr><td>Last Name</td><td>{items.lastName}</td></tr>
                          </tbody>
                      </table>
                      <span className="pull-right">
                        <a onClick={this.logout}>Log out</a>
                      </span>
                  </div>
              );
          }
          else {
              return (
                  <div>
                      <h2>Users</h2>
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
                                      <th scope="row"><a href={`/${dbName}/users/${item.uuid}`}>{item.uuid}</a></th>
                                  </tr>
                              ))}
                          </tbody>
                      </table>
                      <span className="pull-right">
                        <a onClick={this.logout}>Log out</a>
                      </span>
                  </div>
              );
          }
      }
    }
  
    logout() {
      localStorage.removeItem(`access_token_${dbName}`);
      location.reload();
    }
  
    componentDidMount() {
        var uri = `/${dbName}/api/v1/users`
        if(uuid !== "") uri = uri + `/${uuid}`
        fetch(uri, {
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
              'Authorization': 'Bearer ' + localStorage.getItem(`access_token_${dbName}`)
            }
          })
          .then(res => res.json())
          .then(
              (result) => {
                  this.setState({
                      isLoaded: true,
                      items: result.users,
                      error: result.error,
                      message: result.message
                  });
              },
              (error) => {
                  this.setState({
                      isLoaded: false,
                      items: [],
                      message: "Something happened.",
                      error: true
                  });
              }
          )
    }
}
