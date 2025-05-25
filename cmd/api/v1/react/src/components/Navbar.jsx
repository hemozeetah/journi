import axios from 'axios';
import React, { useRef, useState } from 'react';
import './Navbar.css';

export default function Navbar({ setClaims, setToken }) {
  const [showModal, setShowModal] = useState(false);
  const [isSignUp, setIsSignUp] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });

  const modalRef = useRef(null);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const toggleModal = () => {
    setShowModal(!showModal);
    if (!showModal) {
      setIsSignUp(false);
      setFormData({
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
      });
    }
  };

  const toggleSignUp = () => {
    setIsSignUp(!isSignUp);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (isSignUp) {
      if (formData.password !== formData.confirmPassword) {
        alert("Passwords don't match!");
        return;
      }
      axios.post("http://localhost:8080/v1/users", {
        name: formData.name,
        email: formData.email,
        password: formData.password,
        passwordConfirm: formData.confirmPassword
      }).then(_ => {
        // TODO flash message
        toggleSignUp()
      }).catch(err => {
        console.log(err)
      })
    } else {
      axios.get("http://localhost:8080/v1/auth/token", {
        auth: {
          username: formData.email,
          password: formData.password,
        }
      }).then(res => {
        setClaims(res.data.claims);
        setToken(res.data.token);
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('claims', JSON.stringify(res.data.claims));
        toggleModal()
      }).catch(err => {
        console.log(err)
      })
    }
  };

  const handleClickOutsideModal = (event) => {
    if (modalRef.current && !modalRef.current.contains(event.target)) {
      setShowModal(false);
    }
  };
  document.addEventListener('mousedown', handleClickOutsideModal);

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <a href="/" className="nav-link">Home</a>
        <a href="/programs" className="nav-link">Programs</a>
      </div>

      <div className="navbar-center">
        <input
          type="text"
          placeholder="Search..."
          className="search-bar"
        />
      </div>

      <div className="navbar-right">
        <button className="auth-button" onClick={toggleModal}>
          Sign In/Sign Up
        </button>
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal" ref={modalRef}>
            <h2>{isSignUp ? 'Sign Up' : 'Sign In'}</h2>
            <form onSubmit={handleSubmit}>
              {isSignUp && (
                <div className="form-group">
                  <label>Name</label>
                  <input
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleInputChange}
                    required
                  />
                </div>
              )}
              <div className="form-group">
                <label>Email</label>
                <input
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleInputChange}
                  required
                />
              </div>
              <div className="form-group">
                <label>Password</label>
                <input
                  type="password"
                  name="password"
                  value={formData.password}
                  onChange={handleInputChange}
                  required
                />
              </div>
              {isSignUp && (
                <div className="form-group">
                  <label>Confirm Password</label>
                  <input
                    type="password"
                    name="confirmPassword"
                    value={formData.confirmPassword}
                    onChange={handleInputChange}
                    required
                  />
                </div>
              )}
              <button type="submit" className="submit-button">
                {isSignUp ? 'Sign Up' : 'Sign In'}
              </button>
            </form>
            <p className="toggle-auth">
              {isSignUp ? 'Already have an account? ' : "Don't have an account? "}
              <button onClick={toggleSignUp} className="toggle-button">
                {isSignUp ? 'Sign In' : 'Sign Up'}
              </button>
            </p>
          </div>
        </div>
      )}
    </nav>
  );
};
