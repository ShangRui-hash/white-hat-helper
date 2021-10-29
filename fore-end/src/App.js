import {Redirect, Route, Switch} from 'react-router-dom'
import './App.less';
import Home from "./pages/home";
import HomeLogin from "./pages/homelogin";

function App() {
  return (
      <div className='app'>
          <Switch>
              <Route path='/homelogin' component={HomeLogin}/>
              <Route path='/home' component={Home}/>
              <Redirect to='/homelogin'/>
          </Switch>
      </div>
  );
}

export default App;
