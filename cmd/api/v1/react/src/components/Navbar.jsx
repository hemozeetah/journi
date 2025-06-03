import React, { useState } from 'react';
import { Link } from 'react-router';
import FloatingModal from './FloatingModal';
import './Navbar.css';
import SignInForm from './SignInForm';
import SignUpForm from './SignUpForm';

export default function Navbar({ claims, setClaims, setToken }) {
  const [showModal, setShowModal] = useState(false);
  const [isSignUp, setIsSignUp] = useState(false);

  const handleLogout = () => {
    setClaims(null);
    setToken(null);
    localStorage.removeItem('token');
    localStorage.removeItem('claims');
  }

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <Link to="/" className="nav-link">Home</Link>
        <Link to="/cities" className="nav-link">Cities</Link>
        <Link to="/programs" className="nav-link">Programs</Link>
      </div>

      <div className="navbar-center">
        <input
          type="text"
          placeholder="Search..."
          className="search-bar"
        />
      </div>

      <div className="navbar-right">
        {claims && (
          <button className="auth-button" onClick={handleLogout}>
            logout
          </button>
        )}
        {!claims && (
          <button className="auth-button" onClick={() => setShowModal(true)}>
            Sign In/Sign Up
          </button>
        )}
      </div>

      {showModal && (
        <FloatingModal setShowModal={setShowModal}>
          {isSignUp ? (
            <SignUpForm
              setIsSignUp={setIsSignUp}
            />
          ) : (
            <SignInForm
              setClaims={setClaims}
              setToken={setToken}
              setShowModal={setShowModal}
              setIsSignUp={setIsSignUp}
            />
          )
          }
        </FloatingModal>
      )}
    </nav>
  );
};
