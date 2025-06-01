import axios from "axios";
import { useEffect, useState } from "react";
import CityForm from "../components/CityForm";
import CityList from "../components/CityList";
import FloatingModal from "../components/FloatingModal";

export default function CityListPage({ claims, token }) {
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
      <div className="city-list-container">
        <h1>City List Page</h1>
        <div className="city-list">
          {cities && (
            <CityList cities={cities} />
          )}
          {isAdmin && (
            <div className="city-card add-city-card" onClick={() => setShowModal(true)}>
              <div className="add-city-content">
                <div className="plus-sign">+</div>
              </div>
            </div>
          )}
        </div>
        {showModal && (
          <FloatingModal setShowModal={setShowModal}>
            <CityForm
              token={token}
              setCities={setCities}
              setShowModal={setShowModal}
            />
          </FloatingModal>
        )}
    </div >
    </>
  );
}
