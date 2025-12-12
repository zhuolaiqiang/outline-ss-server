import React, { useState } from 'react';
import './FAQ.css';

const FAQItem = ({ question, answer }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="faq-item">
      <button className={`faq-question ${isOpen ? 'active' : ''}`} onClick={() => setIsOpen(!isOpen)}>
        {question}
        <span className="faq-toggle">{isOpen ? '-' : '+'}</span>
      </button>
      <div className={`faq-answer ${isOpen ? 'open' : ''}`}>
        <p>{answer}</p>
      </div>
    </div>
  );
};

const FAQ = () => {
  const faqs = [
    {
      question: "How do I setup Shadowsocks?",
      answer: "Download the client for your device (Windows/Mac/Android/iOS), enter the server details provided in your dashboard, and connect."
    },
    {
      question: "Do you offer refunds?",
      answer: "Yes, we offer a 7-day money-back guarantee if you are not satisfied with our service."
    },
    {
      question: "What does CN2 GIA mean?",
      answer: "CN2 GIA (Global Internet Access) is the highest tier of China Telecom's international network, offering the lowest latency and highest stability."
    },
    {
      question: "Can I watch Netflix?",
      answer: "Yes, our Gold and Platinum plans include streaming optimization for Netflix, Disney+, and more."
    }
  ];

  return (
    <section className="faq-section">
      <div className="container">
        <h2 className="section-title">Frequently Asked Questions</h2>
        <div className="faq-list">
          {faqs.map((faq, index) => (
            <FAQItem key={index} {...faq} />
          ))}
        </div>
      </div>
    </section>
  );
};

export default FAQ;
