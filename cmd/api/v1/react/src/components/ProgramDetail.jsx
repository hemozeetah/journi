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
            {new Date(program.startDate).toLocaleDateString()} - {new Date(program.endDate).toLocaleDateString()}
          </span>
          <p className="program-caption">{program.caption}</p>
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
                    <span style={{ fontStyle: 'italic' }}>
                      {new Date(place.startDatetime).toLocaleString()} - {new Date(place.endDatetime).toLocaleString()}
                    </span>
                    <p>{place.caption}</p>
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
