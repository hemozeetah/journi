import { useNavigate } from "react-router";
import "./ProgramList.css";

export default function ProgramList({ programs }) {
  const navigate = useNavigate();

  const now = new Date();

  // Categorize programs
  const upcomingPrograms = programs.filter(program => new Date(program.startDate) > now);
  const currentPrograms = programs.filter(program =>
    new Date(program.startDate) <= now && new Date(program.endDate) >= now
  );
  const pastPrograms = programs.filter(program => new Date(program.endDate) < now);


  const handleProgramClick = (programID) => {
    navigate(`/programs/${programID}`);
  };

  const renderProgramCard = (program) => {
    const startDate = new Date(program.startDate);
    const endDate = new Date(program.endDate);
    const endDatePlusOne = new Date(endDate);
    endDatePlusOne.setDate(endDatePlusOne.getDate() + 1);
    const isCurrent = startDate <= now && endDatePlusOne >= now;
    const isUpcoming = startDate > now;

    return (
      <div
        key={program.id}
        className={`program-card ${isCurrent ? 'current' : ''} ${isUpcoming ? 'upcoming' : 'past'}`}
        onClick={() => handleProgramClick(program.id)}
      >
        <div className="program-header">
          <h3>{program.companyName}</h3>
          <span className="program-date">
            {startDate.toLocaleDateString()} - {endDate.toLocaleDateString()}
            {isCurrent && <span className="status-badge current-badge">Ongoing</span>}
            {isUpcoming && <span className="status-badge upcoming-badge">Upcoming</span>}
          </span>
        </div>
        <div className="program-caption">
          <p>{program.caption.length > 100
            ? `${program.caption.substring(0, 100)}...`
            : program.caption}</p>
        </div>
        <div className="program-meta">
          <span className="days-remaining">
            {isCurrent && `${Math.ceil((endDatePlusOne - now) / (1000 * 60 * 60 * 24))} days remaining`}
            {isUpcoming && `Starts in ${Math.ceil((startDate - now) / (1000 * 60 * 60 * 24))} days`}
            {!isCurrent && !isUpcoming && `Ended ${Math.ceil((now - endDatePlusOne) / (1000 * 60 * 60 * 24))} days ago`}
          </span>
        </div>
      </div>
    );
  };

  return (
    <div className="program-list-container">
      {upcomingPrograms.length > 0 && (
        <div className="program-section">
          <h2>Upcoming Programs</h2>
          <div className="program-grid">
            {upcomingPrograms.map(renderProgramCard)}
          </div>
        </div>
      )}

      {currentPrograms.length > 0 && (
        <div className="program-section">
          <h2>Current Programs</h2>
          <div className="program-grid">
            {currentPrograms.map(renderProgramCard)}
          </div>
        </div>
      )}

      {pastPrograms.length > 0 && (
        <div className="program-section">
          <h2>Past Programs</h2>
          <div className="program-grid">
            {pastPrograms.map(renderProgramCard)}
          </div>
        </div>
      )}
    </div>
  );
}
