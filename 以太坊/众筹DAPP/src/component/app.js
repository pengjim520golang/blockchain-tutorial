import React,{Component} from 'react';
import Projects from './projects';
import Details from './details';
import createProject from './createProject';
import {BrowserRouter,Switch,Route} from 'react-router-dom';

class App extends Component {

  constructor(...args){
    super(...args)
    console.log(this.props)
  }
  render() {
   
    return (

     

      <div>
        <nav className="navbar navbar-inverse navbar-fixed-top">
            <div className="container-fluid">
              <div className="navbar-header">
                <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
                <span className="sr-only">Toggle navigation</span>
                <span className="icon-bar"></span>
                <span className="icon-bar"></span>
                <span className="icon-bar"></span>
                </button>          
                <a className="navbar-brand" href="#">
                  众筹DAPP
                </a>
              </div>

              <div className="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
                  <ul className="nav navbar-nav">
                    <li><a href="/" className="active">项目列表</a></li>
                    <li><a href="/createProject">发起项目</a></li>
                  </ul>
              </div>
          </div>    
        </nav>


        <BrowserRouter>
        <div>
          <Switch>
            <Route exact path="/" component={Projects} />
            <Route exact path="/details/:address" component={Details} />
            <Route exact path="/createProject" component={createProject} />
          </Switch>
        </div>
      </BrowserRouter>

      </div>

    );
  }
}

export default App;