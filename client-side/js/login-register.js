// function to encrypt/hash the password using SHA-256
async function hashPassword(password) {
  const encoder = new TextEncoder();
  const data = encoder.encode(password);

  // hash the password using SHA-256
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
  return hashHex; // return the hashed password as a hex string
}

// register form
document.getElementById('registerForm').addEventListener('submit', async function(event) {
  event.preventDefault();

  const name = document.getElementById('name').value;
  const email = document.getElementById('registerEmail').value;
  const contactNo = document.getElementById('registerContactNo').value;
  const password = document.getElementById('registerPassword').value;

  // hash the password client-side using SHA-256
  const hashedPassword = await hashPassword(password);

  const requestData = {
      Name: name,
      Email: email,
      contact_no: contactNo,
      hashed_password: hashedPassword 
  };

  // send the data to the user service
  fetch('http://localhost:5001/api/v1/register', {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestData)
  })
  .then(response => response.json())
  .then(data => {
      document.getElementById('responseMessage').innerHTML = `Registration Success: ${data.message}`;
  })
  .catch(error => {
      document.getElementById('responseMessage').innerHTML = `Error: ${error.message}`;
  });
});

// login form
document.getElementById('loginForm').addEventListener('submit', async function(event) {
  event.preventDefault();

  const email = document.getElementById('loginEmail').value;
  const password = document.getElementById('loginPassword').value;

  // hash the password client-side using SHA-256
  const hashedPassword = await hashPassword(password);

  const loginData = {
      email: email,
      hashedPassword: hashedPassword 
  };

  // send the data to the user service
  fetch('http://localhost:5001/api/v1/login', {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      },
      body: JSON.stringify(loginData)
  })
  .then(response => response.json())
  .then(data => {
      document.getElementById('responseMessage').innerHTML = `Login Success: ${data.message}`;
  })
  .catch(error => {
      document.getElementById('responseMessage').innerHTML = `Error: ${error.message}`;
  });
});
