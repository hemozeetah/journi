import axios from "axios";
import { useEffect, useState } from "react";
import CityForm from "../components/CityForm";
import FloatingModal from "../components/FloatingModal";
import "./Cities.css";

export default function Cities({ claims, token }) {
  const isAdmin = claims && claims.role === "admin";
  const [showModal, setShowModal] = useState(false);

  const [cities, setCities] = useState([]);
  useEffect(() => {
    axios.get("http://localhost:8080/v1/cities")
      .then(res => {
        setCities(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      })
  }, []);

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
              <CityForm
                token={token}
                setCities={setCities}
                setShowModal={setShowModal}
              />
            </FloatingModal>
          )}
        </>
      )}
    </>
  );
}
