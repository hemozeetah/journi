import { useEffect, useRef } from "react";
import "./FloatingModal.css";

export default function FloatingModal({ children, setShowModal }) {
  const modalRef = useRef(null);
  const handleClickOutsideModal = (event) => {
    if (modalRef.current && !modalRef.current.contains(event.target)) {
      setShowModal(false);
    }
  };
  useEffect(() => {
    document.addEventListener('mousedown', handleClickOutsideModal);
  }, []);

  return (
    <>
      <div className="modal-overlay">
        <div className="modal" ref={modalRef}>
          {children}
        </div>
      </div>
    </>
  );
}
