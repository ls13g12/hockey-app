import { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';
import { GetWelcomeMsg } from './helpers/welcome_message';

function App() {
  const [welcomeMsg, setWelcomeMsg] = useState('');

  useEffect(() => {
    const getWelcomeMsg = async () => {
      const message = await GetWelcomeMsg();
      setWelcomeMsg(message);
    }
    getWelcomeMsg();
  }, [])

  const welcomeMsgElement = () => {
    if (welcomeMsg) {
      return (
        <p>
          {welcomeMsg}
        </p>
      )
    } else {
      return (
        <p>
          loading...
        </p>
      )
    }
  }

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          {welcomeMsgElement()}
        </p>
      </header>
    </div>
  );
}

export default App;
