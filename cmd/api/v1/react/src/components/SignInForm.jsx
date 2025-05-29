import axios from "axios";
import { useState } from "react";

export default function SignInForm({ setClaims, setToken, setShowModal, setIsSignUp }) {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    axios.get("http://localhost:8080/v1/auth/token", { auth: formData })
      .then(res => {
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('claims', JSON.stringify(res.data.claims));
        setClaims(res.data.claims);
        setToken(res.data.token);
        setShowModal(false);
        // TODO flash message
      })
      .catch(err => {
        console.log(err)
      })
  };

  return (
    <>
      <h2>Sign In</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            name="username"
            value={formData.username}
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
        <button type="submit" className="submit-button">
          Sign In
        </button>
      </form>
      <p className="toggle-auth">
        Don't have an account? <button
          className="toggle-button"
          onClick={() => setIsSignUp(true)}
        >
          Sign Up
        </button>
      </p>
    </>
  );
}
