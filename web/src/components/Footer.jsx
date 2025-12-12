import React from 'react';
import './Footer.css';

const Footer = () => {
  return (
    <footer className="footer">
      <div className="container">
        <div className="footer-columns">
          <div className="footer-col">
            <h4 className="footer-title">About Naiixi</h4>
            <p>Providing premium network acceleration services with industry-leading stability and speed. CN2 GIA & IPLC optimized.</p>
          </div>
          <div className="footer-col">
            <h4 className="footer-title">Products</h4>
            <ul>
              <li><a href="#">Shadowsocks</a></li>
              <li><a href="#">V2Ray / V-Less</a></li>
              <li><a href="#">IPLC Business</a></li>
              <li><a href="#">Dedicated Servers</a></li>
            </ul>
          </div>
          <div className="footer-col">
            <h4 className="footer-title">Support</h4>
            <ul>
              <li><a href="#">Knowledgebase</a></li>
              <li><a href="#">Open Ticket</a></li>
              <li><a href="#">Network Status</a></li>
              <li><a href="#">Contact Us</a></li>
            </ul>
          </div>
          <div className="footer-col">
            <h4 className="footer-title">Legal</h4>
            <ul>
              <li><a href="#">Terms of Service</a></li>
              <li><a href="#">Privacy Policy</a></li>
              <li><a href="#">Refund Policy</a></li>
              <li><a href="#">SLA</a></li>
            </ul>
          </div>
        </div>
        <div className="footer-bottom">
          <p className="copyright">&copy; {new Date().getFullYear()} Naiixi / Just My Socks. All Rights Reserved.</p>
          <div className="payment-methods">
            <span>💳 Visa</span>
            <span>💳 MasterCard</span>
            <span>💳 Alipay</span>
            <span>💳 Crypto</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
