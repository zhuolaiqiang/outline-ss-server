import React from 'react';
import './Hero.css';

const Hero = () => {
  return (
    <section className="hero">
      <div className="hero-content">
        <h1 className="hero-title">Just My Socks</h1>
        <p className="hero-subtitle">Shadowsocks & V2Ray & IPLC Service</p>
        <p className="hero-description">
          We provide high speed and reliable Shadowsocks & V2Ray service for you.<br/>
          CN2 GIA / Direct connections for best speed.
        </p>
        <div className="hero-buttons">
          <a href="#" className="btn btn-primary">Order Now</a>
          <a href="#" className="btn btn-secondary">Learn More</a>
        </div>
      </div>
    </section>
  );
};

export default Hero;
