import React, { useState } from 'react';
import './Announcements.css';

const Announcements = () => {
  const [isVisible, setIsVisible] = useState(true);

  if (!isVisible) return null;

  return (
    <div className="announcement-bar">
      <div className="container">
        <span className="icon">📢</span>
        <span className="message">Important: Server maintenance scheduled for Dec 15th, 02:00 UTC. Please anticipate brief interruptions.</span>
        <button className="close-btn" onClick={() => setIsVisible(false)}>×</button>
      </div>
    </div>
  );
};

export default Announcements;
