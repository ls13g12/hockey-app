import { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';
import { GetBackendHealth } from './helpers/backend_healthcheck';

function App() {
  const [backendHealthMessage, setBackendHealthMessage] = useState('');

  useEffect(() => {
    const GetBackendHealthMessage = async () => {
      const message = await GetBackendHealth();
      setBackendHealthMessage(message);
    }
    GetBackendHealthMessage();
  }, [])

  const healthcheckElement = () => {
    if (backendHealthMessage) {
      return (
        <p>
          {backendHealthMessage}
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
          {healthcheckElement()}
        </p>
      </header>
    </div>
  );
}

export default App;
