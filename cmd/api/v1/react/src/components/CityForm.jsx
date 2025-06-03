import axios from "axios";
import { useState } from "react";

export default function CityForm({ token, setCities, setShowModal }) {
  const [data, setData] = useState({
    name: '',
    caption: ''
  });
  const [images, setImages] = useState([]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setData({
      ...data,
      [name]: value
    });
  };

  const handleImageChange = (e) => {
    if (e.target.files) {
      setImages(Array.from(e.target.files));
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('data', JSON.stringify(data));
    images.forEach((image, _) => {
      formData.append(`images`, image);
    });
    axios.post("http://localhost:8080/v1/cities", formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        // TODO flash message
        setCities(cities => [...cities, res.data]);
        setShowModal(false);
      })
      .catch(err => {
        console.log(err.response.data)
      });
  };

  return (
    <>
      <h2>Add City</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name:</label>
          <input
            type="text"
            name="name"
            value={data.name}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="form-group">
          <label>Caption:</label>
          <input
            type="text"
            name="caption"
            value={data.caption}
            onChange={handleInputChange}
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
          Add City
        </button>
      </form>
    </>
  );
}
