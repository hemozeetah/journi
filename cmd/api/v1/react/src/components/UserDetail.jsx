import { useState } from "react";
import "./UserDetail.css";

export default function UserDetail({ user, claims, token }) {
  const [image, setImage] = useState(null);

  const handleSubmit = (e) => {
    e.preventDefault();
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
            src={user.profile || '/profile.png'}
            alt="Profile"
            className="profile-image"
          />
          {user.id == claims.id && (
            <>
              <form onSubmit={handleSubmit}>
                <input
                  type="file"
                  name="image"
                  onChange={(e) => setImage(e.target.value)}
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
