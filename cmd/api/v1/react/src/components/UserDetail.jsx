import { useState } from "react";
import "./UserDetail.css";
import axios from "axios";

export default function UserDetail({ user, claims, token }) {
  const [images, setImages] = useState([]);

  const handleImageChange = (e) => {
    if (e.target.files) {
      setImages(Array.from(e.target.files));
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('data', JSON.stringify({name: "ayham"}));
    images.forEach((image, _) => {
      formData.append(`images`, image);
    });
    axios.put("http://localhost:8080/v1/users/" + claims.id, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        console.log(res.data);
      })
      .catch(err => {
        console.log(err.response.data)
      });
  };

  if (!user) {
    return <div>user not found</div>
  }

  return (
    <>
      <div className="user-card">
        <div className="user-info">
          <h2 className="user-name">{user.name}</h2>
          <p className="user-email">{user.email}</p>
          <div className="user-role">
            <span className={`role-badge ${user.role.toLowerCase()}`}>
              {user.role}
            </span>
          </div>
        </div>
        <div className="user-profile-detail">
          <img
            src={user.profileImageURL ? "http://localhost:8080" + user.profileImageURL : '/profile.png'}
            alt="Profile"
            className="profile-image"
          />
          {claims && user.id == claims.id && (
            <>
              <form onSubmit={handleSubmit}>
                <input
                  type="file"
                  name="image"
                  onChange={handleImageChange}
                  multiple
                  accept="image/*"
                  required
                />
                <button type="submit">
                  change
                </button>
              </form>
            </>
          )}
        </div>
      </div>
    </>
  );
}
