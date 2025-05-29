import axios from "axios";
import { useEffect, useState } from "react";
import FloatingModal from "../components/FloatingModal";
import "./Cities.css";

export default function Cities({ claims, token }) {
  const isAdmin = claims && claims.role === "admin";
  const [name, setName] = useState('');
  const [caption, setCaption] = useState('');
  const [images, setImages] = useState([]);
  const [showModal, setShowModal] = useState(false);


  const handleImageChange = (e) => {
    if (e.target.files) {
      setImages(Array.from(e.target.files));
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('data', JSON.stringify({
      name: name,
      caption: caption
    }));
    images.forEach((image, _) => {
      formData.append(`images`, image);
    });
    axios.post("http://localhost:8080/v1/cities", formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': "Bearer " + token
      }
    }).then(res => {
      // TODO flash message
      setCities(cities => [...cities, res.data]);
      setShowModal(false);
    }).catch(err => {
      console.log(err.response.data)
    })
  };

  //

  const [cities, setCities] = useState([]);
  useEffect(() => {
    axios.get("http://localhost:8080/v1/cities")
      .then(res => {
        setCities(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      })
  }, [])

  return (
    <>
      <h1>Cities Page</h1>
      {cities && (
        <ul>
          {cities.map(city => (
            <li key={city.id}>
              {city.name}
            </li>
          ))}
        </ul>
      )}
      {isAdmin && (
        <>
          <button className="modal-toggle-button" onClick={() => setShowModal(true)}>
            +
          </button>
          {showModal && (
            <FloatingModal setShowModal={setShowModal}>
              <h2>Add City</h2>
              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label>Name:</label>
                  <input
                    type="text"
                    name="name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Caption:</label>
                  <input
                    type="text"
                    name="caption"
                    value={caption}
                    onChange={(e) => setCaption(e.target.value)}
                    required
                  />
                </div>

                <div className="form-group">
                  <label>Select Images:</label>
                  <input
                    type="file"
                    name="images"
                    onChange={handleImageChange}
                    multiple
                    accept="image/*"
                    required
                  />
                </div>

                <button type="submit" className="submit-button">
                  Add
                </button>
              </form>
            </FloatingModal>
          )}
        </>
      )}
    </>
  );
}
