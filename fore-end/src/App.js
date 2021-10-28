import {Redirect, Route, Switch} from 'react-router-dom'
import './App.less';
import Home from "./components/home";
import HomeLogin from "./components/homelogin";

function App() {
  return (
      <div className='app'>
          <Switch>
              <Route path='/homelogin' component={HomeLogin}/>
              <Route path='/home' component={Home}/>
              <Redirect to='/home'/>
          </Switch>
      </div>
  );
}

export default App;
