import { useState } from 'react';
import './App.css';
import Navbar from './components/Navbar';
import CreatePost from './components/CreatePost';

function App() {
  const [claims, setClaims] = useState(null);
  const [token, setToken] = useState(null);

  return (
    <>
      <Navbar
        setClaims={setClaims}
        setToken={setToken}
      />
      {token && <p>{token}</p>}
      {claims && (
        <div>
          <p>{claims.name}</p>
          <p>{claims.role}</p>
        </div>
      )}

      <CreatePost />
    </>
  )
}

export default App
