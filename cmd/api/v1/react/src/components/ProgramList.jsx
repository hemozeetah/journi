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
          key={program.id}
          className="program-card"
          onClick={() => handleProgramClick(program.id)}
        >
          <h3>{program.caption}</h3>
          <p>Company ID: {program.companyID}</p>
          <p>
            Duration: {new Date(program.startDate).toLocaleDateString()} -
            {new Date(program.endDate).toLocaleDateString()}
          </p>
        </div>
      ))}
    </>
  );
}
