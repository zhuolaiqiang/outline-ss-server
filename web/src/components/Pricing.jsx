import React from 'react';
import './Pricing.css';

const Pricing = () => {
  const plans = [
    {
      name: "Silver",
      price: "$5.88",
      period: "/mo",
      features: [
        "500GB Bandwidth",
        "3 Devices",
        "CN2 GIA Network",
        "IPLC Enabled",
        "24/7 Support"
      ],
      recommended: false
    },
    {
      name: "Gold",
      price: "$9.88",
      period: "/mo",
      features: [
        "1TB Bandwidth",
        "5 Devices",
        "CN2 GIA Network",
        "IPLC Enabled",
        "Priority Support"
      ],
      recommended: true
    },
    {
      name: "Platinum",
      price: "$15.88",
      period: "/mo",
      features: [
        "Unlimited Bandwidth",
        "Unlimited Devices",
        "CN2 GIA Network",
        "IPLC Enabled",
        "Dedicated Manager"
      ],
      recommended: false
    }
  ];

  return (
    <section className="pricing">
      <div className="container">
        <h2 className="section-title">Choose Your Plan</h2>
        <div className="pricing-grid">
          {plans.map((plan, index) => (
            <div className={`pricing-card ${plan.recommended ? 'recommended' : ''}`} key={index}>
              {plan.recommended && <div className="badge">Best Value</div>}
              <h3 className="plan-name">{plan.name}</h3>
              <div className="plan-price">
                <span className="amount">{plan.price}</span>
                <span className="period">{plan.period}</span>
              </div>
              <ul className="plan-features">
                {plan.features.map((feature, idx) => (
                  <li key={idx}>✓ {feature}</li>
                ))}
              </ul>
              <button className="btn-plan">Order Now</button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Pricing;
