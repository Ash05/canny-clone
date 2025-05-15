import React, { useEffect } from 'react';

function Register() {
  useEffect(() => {
    window.google.accounts.id.initialize({
      client_id: 'YOUR_GOOGLE_CLIENT_ID',
      callback: handleCredentialResponse
    });
    window.google.accounts.id.renderButton(
      document.getElementById('google-signup-button'),
      { theme: 'outline', size: 'large' }
    );
  }, []);

  const handleCredentialResponse = async (response) => {
    console.log('Encoded JWT ID token: ' + response.credential);
    try {
      const res = await fetch('http://localhost:8080/verify-token', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token: response.credential }),
      });
      const data = await res.json();
      console.log('Backend response:', data);
    } catch (error) {
      console.error('Error verifying token:', error);
    }
  };

  return (
    <div>
      <h1>Register</h1>
      <div id="google-signup-button"></div>
    </div>
  );
}

export default Register;
