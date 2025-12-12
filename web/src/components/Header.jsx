import React, { useState, useEffect } from 'react';
import './Header.css';

const Header = () => {
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <>
      {/* Top Utility Bar */}
      <div className="top-bar">
        <div className="container">
          <div className="top-bar-left">
            <a href="#">Network Status</a>
            <a href="#">Knowledgebase</a>
          </div>
          <div className="top-bar-right">
            <span className="lang-selector">
              Language: 
              <select>
                <option value="en">English</option>
                <option value="cn">中文</option>
              </select>
            </span>
            <div className="auth-links">
              <a href="#" className="login-link">Login</a>
              <a href="#" className="register-link">Register</a>
            </div>
          </div>
        </div>
      </div>

      {/* Main Header */}
      <header className={`header ${isScrolled ? 'scrolled' : ''}`}>
        <div className="container">
          <div className="logo">
            <a href="/">
              <span className="logo-icon">🌐</span>
              <span className="logo-text">naiixi</span>
            </a>
          </div>
          <nav className="nav">
            <ul className="nav-list">
              <li><a href="#" className="active">Home</a></li>
              <li className="dropdown">
                <a href="#">Store ▾</a>
                <ul className="dropdown-menu">
                  <li><a href="#">Shadowsocks</a></li>
                  <li><a href="#">V2Ray</a></li>
                  <li><a href="#">IPLC Services</a></li>
                </ul>
              </li>
              <li><a href="#">Announcements</a></li>
              <li><a href="#">Support</a></li>
              <li><a href="#">Contact Us</a></li>
            </ul>
          </nav>
        </div>
      </header>
    </>
  );
};

export default Header;
