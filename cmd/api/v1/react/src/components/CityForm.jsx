import axios from "axios";
import { useState, useEffect } from "react";
import "./CityForm.css";

export default function CityForm({ token, setCities, setShowModal, city = null, setCity }) {
  const [data, setData] = useState({
    name: '',
    caption: ''
  });
  const [images, setImages] = useState([]);

  useEffect(() => {
    if (city) {
      setData({
        name: city.name,
        caption: city.caption
      });
    }
  }, [city]);

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

  const removeImage = (index) => {
    const updatedImages = [...images];
    updatedImages.splice(index, 1);

    setImages(updatedImages);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const formData = new FormData();
    formData.append('data', JSON.stringify(data));

    images.forEach((image) => {
      formData.append(`images`, image);
    });

    const url = city
      ? `http://localhost:8080/v1/cities/${city.id}`
      : "http://localhost:8080/v1/cities";

    const method = city ? 'put' : 'post';

    try {
      const res = await axios[method](url, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': "Bearer " + token
        }
      });

      if (city) {
        setCity(res.data);
      } else {
        setCities(cities => [...cities, res.data]);
      }
      setShowModal(false);
    } catch (err) {
      console.error(err.response?.data || err.message);
    }
  };

  return (
    <>
      <h2>{city ? 'Edit City' : 'Add City'}</h2>
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
          />
        </div>

        <div className="form-group">
          <label>Select Images</label>
          <input
            type="file"
            name="images"
            onChange={handleImageChange}
            multiple
            accept="image/*"
            required={!city}
          />

          {/* All image previews */}
          {images.length > 0 && (
            <div className="image-grid">
              {images.map((image, index) => (
                <div key={index} className="image-preview">
                  <img
                    src={URL.createObjectURL(image)}
                    alt={`image ${index}`}
                  />
                  <button
                    type="button"
                    onClick={() => removeImage(index)}
                    className="remove-btn"
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>

        <button type="submit" className="submit-button">
          {city ? 'Update City' : 'Add City'}
        </button>
      </form>
    </>
  );
}
