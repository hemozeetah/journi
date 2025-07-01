import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import UserDetail from "../components/UserDetail";

export default function SubscriberDetailPage({ claims, token }) {
  const { id } = useParams();

  const [user, setUser] = useState(null);
  const [program, setProgram] = useState(null);

  const toggleAccepted = async (referenceID, currentAccepted) => {
    try {
      const newAccepted = !currentAccepted;
      const res = await axios.put(
        `http://localhost:8080/v1/subscribers/${referenceID}`,
        { accepted: newAccepted },
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " + token
          }
        }
      );
      console.log(res.data);

      setUser(user => ({
        ...user,
        accepted: newAccepted
      }));
    } catch (err) {
      console.error('Failed to update subscriber:', err);
      alert('Failed to update subscriber status');
    }
  };

  useEffect(() => {
    axios.get(`http://localhost:8080/v1/subscribers/${id}`)
      .then(resSub => {
        console.log(resSub.data);
        axios.get(`http://localhost:8080/v1/users/${resSub.data.id}`)
          .then(res => {
            console.log(res.data);
            setUser({
              ...res.data,
              ...resSub.data,
            });
          })
          .catch(err => {
            console.log(err.response.data)
          });
        axios.get(`http://localhost:8080/v1/programs/${resSub.data.programID}`)
          .then(res => {
            console.log(res.data);
            setProgram(res.data);
          })
          .catch(err => {
            console.log(err.response.data)
          });
      })
      .catch(err => {
        console.log(err.response.data)
      });
  }, [id]);

  if (!program || claims.id !== program.companyID) {
    return <>NO Data</>;
  }

  return (
    <>
      <UserDetail
        user={user}
        claims={claims}
        token={token}
      />
      <div className="subscriber-list"
        style={{
          margin: "0 600px",
        }}
      >
        <table>
          <thead>
            <tr>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{user.accepted ? 'Accepted' : 'Pending'}</td>
              <td>
                <button
                  onClick={() => toggleAccepted(user.referenceID, user.accepted)}
                  className={`toggle-btn ${user.accepted ? 'active' : ''}`}
                >
                  {user.accepted ? 'Revoke' : 'Accept'}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </>
  );
};
