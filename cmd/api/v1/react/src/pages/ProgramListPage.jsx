import { useEffect, useState } from "react";
import ProgramList from "../components/ProgramList";
import FloatingModal from "../components/FloatingModal";
import axios from "axios";
import ProgramForm from "../components/ProgramForm";

export default function ProgramListPage({ claims, token }) {
  const isCompanyOrAdmin = claims && (claims.role === "company" || claims.role === "admin");
  const [showModal, setShowModal] = useState(false);

  const [programs, setPrograms] = useState([]);
  const [cities, setCities] = useState([]);
  const [places, setPlaces] = useState([]);
  useEffect(() => {
    axios.get("http://localhost:8080/v1/programs")
      .then(res => {
        setPrograms(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
    axios.get("http://localhost:8080/v1/cities")
      .then(res => {
        setCities(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
    axios.get("http://localhost:8080/v1/places")
      .then(res => {
        setPlaces(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
  }, []);

  return (
    <>
      <div className="program-list-container">
        <h1>Program List Page</h1>
        <div className="program-list">
          {isCompanyOrAdmin && (
            <div className="program-card add-program-card" onClick={() => setShowModal(true)}>
              <div className="add-program-content">
                <div className="plus-sign">+</div>
              </div>
            </div>
          )}
          <ProgramList programs={programs} />
        </div>
        {showModal && (
        <FloatingModal setShowModal={setShowModal}>
            <ProgramForm
              cities={cities}
              places={places}
              token={token}
              setPrograms={setPrograms}
              setShowModal={setShowModal}
            />
        </FloatingModal>
        )}
      </div>
    </>
  );
}
