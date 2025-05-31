import axios from "axios";
import { useEffect, useState } from "react";
import CityList from "../components/CityList";
import CityForm from "../components/CityForm";
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
      <h1>City List Page</h1>
      {cities && (
        <CityList cities={cities} />
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
