import axios from "axios";
import { QRCodeSVG } from "qrcode.react";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import FloatingModal from "../components/FloatingModal";
import ProgramDetail from "../components/ProgramDetail";
import SettingsButton from "../components/SettingsButton";
import SubscriberList from "../components/SubscriberList";
import ProgramForm from "../components/ProgramForm";

export default function ProgramDetailPage({ claims, token }) {
  const { id } = useParams();
  const navigate = useNavigate();

  const [showModal, setShowModal] = useState(false);
  const [showModalEdit, setShowModalEdit] = useState(false);

  const [program, setProgram] = useState(null);
  const [company, setCompany] = useState(null);
  const [currentPlaces, setCurrentPlaces] = useState([]);
  const [subscribers, setSubscribers] = useState([]);

  const [isRegister, setIsRegister] = useState(false);
  const [subscriber, setSubscriber] = useState(null);

  const [cities, setCities] = useState([]);
  const [places, setPlaces] = useState([]);

  const handleRegisteration = () => {
    if (!isRegister) {
      const formData = new FormData();
      formData.append('programID', id);
      axios.post("http://localhost:8080/v1/subscribers", formData, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': "Bearer " + token
        }
      })
        .then(res => {
          console.log(res.data);
          setSubscriber(res.data);
          setIsRegister(true);
        })
        .catch(err => {
          console.log(err.response.data);
        });
    } else {
      axios.delete(`http://localhost:8080/v1/subscribers/${subscriber.referenceID}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': "Bearer " + token
        }
      })
        .then(res => {
          console.log(res.data);
          setIsRegister(false);
        })
        .catch(err => {
          console.log(err.response.data);
        });
    }
  };


  useEffect(() => {
    const fetchData = async () => {
      try {
        const programRes = await axios.get(`http://localhost:8080/v1/programs/${id}`);
        setProgram(programRes.data);

        const companyRes = await axios.get(`http://localhost:8080/v1/users/${programRes.data.companyID}`);
        setCompany(companyRes.data);

        const journeysRes = await axios.get(`http://localhost:8080/v1/journeys?program_id=${id}`);
        const placesData = await Promise.all(
          journeysRes.data.map(journey =>
            axios.get(`http://localhost:8080/v1/places/${journey.placeID}`)
              .then(res => {
                return {
                  ...res.data,
                  startDatetime: journey.startDatetime,
                  endDatetime: journey.endDatetime
                }
              })
              .catch(err => {
                console.log(err.response.data);
                return null; // or handle error differently
              })
          )
        );

        setCurrentPlaces(placesData.filter(place => place !== null));

      } catch (err) {
        console.log(err.response?.data || err.message);
      }
    };

    fetchData();

    const fetchSubscribers = async () => {
      try {
        // Fetch subscribers for the program
        const subResponse = await axios.get(
          `http://localhost:8080/v1/subscribers?program_id=${id}`
        );

        // Fetch user details for each subscriber
        const subscribersWithUserData = await Promise.all(
          subResponse.data.map(async (subscriber) => {
            try {
              const userResponse = await axios.get(
                `http://localhost:8080/v1/users/${subscriber.id}`
              );
              return {
                ...subscriber,
                userName: userResponse.data.name
              };
            } catch (userError) {
              console.error(`Failed to fetch user ${subscriber.id}:`, userError);
              return {
                ...subscriber,
                userName: 'Unknown User'
              };
            }
          })
        );

        setSubscribers(subscribersWithUserData);

        if (claims && subscribersWithUserData.find(sub => sub.id === claims.id)) {
          setSubscriber(subscribersWithUserData.find(sub => sub.id === claims.id));
          setIsRegister(true);
        }
      } catch (err) {
        console.log(err);
      }
    };

    fetchSubscribers();

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

  const handleEdit = () => {
    setShowModalEdit(true);
  };

  const handleDelete = () => {
    axios.delete(`http://localhost:8080/v1/programs/${id}`, {
      headers: {
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        console.log(res.data);
        navigate("/programs");
      })
      .catch(err => {
        console.log(err.response.data);
      });
  };

  if (!program || !company) {
    return <div>Program not found</div>;
  }

  return (
    <>
      <ProgramDetail
        program={program}
        company={company}
        places={currentPlaces}
      />
      {claims && claims.role === "user" && (
        <>
          {isRegister && subscriber.accepted === false && (
            <div className="show-subscribers">
              <div className="show-subscribers-btn" onClick={handleRegisteration}>
                Unregister
              </div>
            </div>
          )}
          {!isRegister && (
            <div className="show-subscribers">
              <div className="show-subscribers-btn" onClick={handleRegisteration}>
                Register
              </div>
            </div>
          )}
          {isRegister && (
            <div className="qr-code-container">
              <QRCodeSVG
                value={`http://localhost:5173/subscribers/${subscriber.referenceID}`}
                size={256}
                level="H"
              />
              <p>QR code for verification</p>
            </div>
          )}
        </>
      )}
      {claims && claims.id === company.id && (
        <>
          <div className="show-subscribers">
            <div className="show-subscribers-btn" onClick={() => setShowModal(true)}>
              Show Subscribers
            </div>
          </div>
          {showModal && <FloatingModal setShowModal={setShowModal}>
            <SubscriberList
              subscribers={subscribers}
              setSubscribers={setSubscribers}
              token={token}
            />
          </FloatingModal>}
          {showModalEdit && (
            <FloatingModal setShowModal={setShowModalEdit}>
              <ProgramForm
                cities={cities}
                places={places}
                claims={claims}
                token={token}
                program={program}
                setProgram={setProgram}
                currentPlaces={currentPlaces}
                setShowModal={setShowModalEdit}
              />
            </FloatingModal>
          )}
          <SettingsButton
            onEdit={handleEdit}
            onDelete={handleDelete}
          />
        </>
      )}
    </>
  );
}
