import axios from 'axios';
import "./SubscriberList.css";

export default function SubscriberList({ subscribers, setSubscribers, token }) {
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

      setSubscribers(subscribers.map(sub =>
        sub.referenceID === referenceID
          ? { ...sub, accepted: newAccepted }
          : sub
      ));
    } catch (err) {
      console.error('Failed to update subscriber:', err);
      alert('Failed to update subscriber status');
    }
  };

  return (
    <>
      <div className="subscriber-list">
        <table>
          <thead>
            <tr>
              <th>User</th>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {subscribers.map((subscriber) => (
              <tr key={subscriber.referenceID}>
                <td>{subscriber.userName}</td>
                <td>{subscriber.accepted ? 'Accepted' : 'Pending'}</td>
                <td>
                  <button
                    onClick={() => toggleAccepted(subscriber.referenceID, subscriber.accepted)}
                    className={`toggle-btn ${subscriber.accepted ? 'active' : ''}`}
                  >
                    {subscriber.accepted ? 'Revoke' : 'Accept'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
};
