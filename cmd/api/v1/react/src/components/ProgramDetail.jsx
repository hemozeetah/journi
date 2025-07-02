import { useNavigate } from "react-router";
import "./ProgramDetail.css";

export default function ProgramDetail({ program, company, places }) {
  const navigate = useNavigate();

  const handleCompanyClick = (companyID) => {
    navigate(`/users/${companyID}`);
  };

  const handlePlaceClick = (placeID) => {
    navigate(`/places/${placeID}`);
  };

  return (
    <>
      <div className="program-detail">
        <div className="program-info">
          <h2
            className="company-name"
            onClick={() => handleCompanyClick(company.id)}
            style={{ cursor: 'pointer', color: 'blue' }}
          >
            {company.name}
          </h2>
          <span style={{ fontStyle: 'italic' }}>
            {new Date(program.startDate).toDateString()} - {new Date(program.endDate).toDateString()}
          </span>
          <pre>{program.caption}</pre>
        </div>

        <div className="program-detail-container">
          <div className="places-section">
            <h3>Places</h3>
            <div className="places-column">
              <ul className="places-list">
                {places.map((place, index) => (
                  <li
                    key={place.id + index}
                    className="place-item"
                    onClick={() => handlePlaceClick(place.id)}
                  >
                    <strong>{place.name}</strong>
                    <p style={{ fontStyle: 'italic' }}>
                      from {new Date(place.startDatetime).toLocaleString()}
                    </p>
                    <p style={{ fontStyle: 'italic' }}>
                      to {new Date(place.endDatetime).toLocaleString()}
                    </p>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
