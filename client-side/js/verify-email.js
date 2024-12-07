document.getElementById('verifyForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const verificationCode = document.getElementById('verificationCode').value;

    const requestData = {
        email: email,
        verificationCode: verificationCode
    };

    // Send the verification request to the backend
    fetch('http://localhost:5001/api/v1/verify-email', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.json())
    .then(data => {
        // Handle success response
        document.getElementById('responseMessage').innerHTML = `Verification Success: ${data.message}`;
    })
    .catch(error => {
        // Handle error response
        document.getElementById('responseMessage').innerHTML = `Error: ${error.message}`;
    });
});
