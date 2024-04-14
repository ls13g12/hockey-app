import { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';
import * as fs from 'fs';
const players = JSON.parse(fs.readFileSync('../mocked_data/players.json'));

function App() {
  const [welcomeMsg, setWelcomeMsg] = useState('No message');
  console.log(players);
  useEffect(() => {
    const getWelcomeMsg = async () => {
      const resp = await fetch('http://localhost:3333/hello', {
        "Content-Type": "application/json"
      });
      if (resp.ok) {
        const data = await resp.json();
        console.log(data);
        setWelcomeMsg(data.welcome_msg);
      }
    }
    getWelcomeMsg();
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          {welcomeMsg}
        </p>
      </header>
    </div>
  );
}

export default App;
