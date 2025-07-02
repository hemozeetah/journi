import { useNavigate } from "react-router";
import "./ProgramList.css";

export default function ProgramList({ programs }) {
  const navigate = useNavigate();

  const handleProgramClick = (programID) => {
    navigate(`/programs/${programID}`);
  };

  return (
    <>
      {programs.map((program) => (
        <div
          key={program.companyID + program.caption}
          className="program-card"
          onClick={() => handleProgramClick(program.id)}
          style={{ cursor: 'pointer' }}
        >
          <h3 style={{ cursor: 'pointer', color: 'blue' }}>
            {program.companyName}
          </h3>
          <span style={{ fontStyle: 'italic' }}>
            {new Date(program.startDate).toLocaleDateString()} - {new Date(program.endDate).toLocaleDateString()}
          </span>
          <div className="program-caption">
            <pre>{program.caption.length > 100
              ? `${program.caption.substring(0, 100)}...`
              : program.caption}</pre>
          </div>
        </div>
      ))}
    </>
  );
}
