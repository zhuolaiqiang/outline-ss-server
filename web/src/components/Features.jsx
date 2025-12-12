import React from 'react';
import './Features.css';

const Features = () => {
  const featureList = [
    {
      title: "Shadowsocks",
      description: "We use Shadowsocks, which is the best protocol for firewalls.",
      icon: "⚡"
    },
    {
      title: "CN2 GIA",
      description: "We use CN2 GIA/Direct Connect for the best speed possible.",
      icon: "🚀"
    },
    {
      title: "Reliability",
      description: "We provide highly reliable service with 99.9% uptime guarantee.",
      icon: "🛡️"
    }
  ];

  return (
    <section className="features">
      <div className="container">
        <h2 className="section-title">Why Choose Us</h2>
        <div className="features-grid">
          {featureList.map((feature, index) => (
            <div className="feature-card" key={index}>
              <div className="feature-icon">{feature.icon}</div>
              <h3 className="feature-title">{feature.title}</h3>
              <p className="feature-desc">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Features;
