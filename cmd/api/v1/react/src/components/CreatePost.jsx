import './CreatePost.css';
import { useRef, useState } from "react";

export default function CreatePost() {
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    caption: '',
    place: ''
  });

  const modalRef = useRef(null);

  const toggleModal = () => {
    setShowModal(!showModal);
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const handleClickOutsideModal = (event) => {
    if (modalRef.current && !modalRef.current.contains(event.target)) {
      setShowModal(false);
    }
  };
  document.addEventListener('mousedown', handleClickOutsideModal);

  const handleSubmit = (e) => {
    e.preventDefault();
  };

  return (
    <>
      <button className="post-button" onClick={toggleModal}>
        +
      </button>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal" ref={modalRef}>
            <h2>Create Post</h2>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Caption</label>
                <input
                  type="text"
                  name="caption"
                  value={formData.caption}
                  onChange={handleInputChange}
                  required
                />
              </div>
              <div className="form-group">
                <label>Place</label>
                <input
                  type="text"
                  name="place"
                  value={formData.place}
                  onChange={handleInputChange}
                  required
                />
              </div>
              <button type="submit" className="submit-button">
                Post
              </button>
            </form>
          </div>
        </div>
      )}
    </>
  )
}
