import axios from "axios";
import { useEffect, useState } from "react";

export default function PlaceForm({ token, city, setPlaces, setShowModal, place = null, setPlace }) {
  const [data, setData] = useState({
    cityID: '',
    name: '',
    caption: '',
    type: ''
  });
  const [images, setImages] = useState([]);

  useEffect(() => {
    if (place) {
      setData({
        cityID: place.cityID,
        name: place.name,
        caption: place.caption,
        type: place.type
      });
    } else {
      setData({
        ...data,
        cityID: city.id
      });
    }
  }, [place]);

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

    images.forEach((image, _) => {
      formData.append(`images`, image);
    });

    const url = place
      ? `http://localhost:8080/v1/places/${place.id}`
      : "http://localhost:8080/v1/places";

    const method = place ? 'put' : 'post';

    try {
      const res = await axios[method](url, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': "Bearer " + token
        }
      });
      console.log(res.data);

      if (place) {
        setPlace(res.data);
      } else {
        setPlaces(places => [...places, res.data]);
      }
      setShowModal(false);
    } catch (err) {
      console.error(err.response?.data || err.message);
    }
  };

  return (
    <>
      <h2>{place ? 'Edit Place' : 'Add Place'}</h2>
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
          <textarea
            name="caption"
            value={data.caption}
            onChange={handleInputChange}
            rows={5}
            required
          />
        </div>
        <div className="form-group">
          <label>Type:</label>
          <select
            name="type"
            value={data.type}
            onChange={handleInputChange}
            required
          >
            <option disabled defaultValue="" value=""> -- select an option -- </option>
            <option value="restaurant">restaurant</option>
            <option value="hotel">hotel</option>
          </select>
        </div>
        <div className="form-group">
          <label>Select Images:</label>
          <input
            type="file"
            name="images"
            onChange={handleImageChange}
            multiple
            accept="image/*"
            required={!place}
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
          {place ? 'Update Place' : 'Add Place'}
        </button>
      </form>
    </>
  );
}
