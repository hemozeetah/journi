import axios from "axios";
import { useState } from "react";
import "./UserDetail.css";

export default function UserDetail({ user, claims, token }) {
  const [images, setImages] = useState([]);
  const [editingRole, setEditingRole] = useState(false);
  const [selectedRole, setSelectedRole] = useState(user?.role || 'user');
  const [isUpdating, setIsUpdating] = useState(false);

  const handleImageChange = (e) => {
    if (e.target.files) {
      setImages(Array.from(e.target.files));
    }
  };

  const handleImageChangeSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('data', JSON.stringify({}));
    images.forEach((image, _) => {
      formData.append(`images`, image);
    });

    axios.put("http://localhost:8080/v1/users/" + user.id, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        console.log(res.data);
        // TODO
        window.location.reload();
      })
      .catch(err => {
        console.log(err.response.data)
      });
  };

  const handleRoleChangeSubmit = async () => {
    if (!user || selectedRole === user.role) {
      setEditingRole(false);
      return;
    }

    setIsUpdating(true);
    try {
      const formData = new FormData();
      formData.append('data', JSON.stringify({ role: selectedRole }));
      await axios.put(
        `http://localhost:8080/v1/users/${user.id}`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${token}`
          }
        }
      );
      // TODO
      window.location.reload();
    } catch (error) {
      console.error("Failed to update role:", error);
    } finally {
      setIsUpdating(false);
    }
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
            {claims?.role === 'admin' && editingRole ? (
              <div className="role-edit-container">
                <select
                  value={selectedRole}
                  onChange={(e) => setSelectedRole(e.target.value)}
                  className="role-select"
                >
                  <option value="user">User</option>
                  <option value="company">Company</option>
                  <option value="admin">Admin</option>
                </select>
                <button
                  onClick={handleRoleChangeSubmit}
                  className="role-save-button"
                  disabled={isUpdating}
                >
                  {isUpdating ? 'Saving...' : 'Save'}
                </button>
                <button
                  onClick={() => setEditingRole(false)}
                  className="role-cancel-button"
                >
                  Cancel
                </button>
              </div>
            ) : (
              <>
                <span className={`role-badge ${user.role.toLowerCase()}`}>
                  {user.role}
                </span>
                {claims && claims.role === "admin" && user.role !== "admin" && (
                  <button
                    onClick={() => {
                      setSelectedRole(user.role);
                      setEditingRole(true);
                    }}
                    className="role-edit-button"
                  >
                    Edit Role
                  </button>
                )}
              </>
            )}
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
              <form onSubmit={handleImageChangeSubmit}>
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
