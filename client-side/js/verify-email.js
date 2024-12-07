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
    .then(response => {
        if (!response.ok) {
            // If the response status is not OK (not 2xx), throw an error with the message from the response
            return response.json().then(errorData => {
                throw new Error(errorData.error || 'An error occurred');
            });
        }
        // If the response is OK, return the response as JSON
        return response.json();
    })
    .then(data => {
        // Handle success response
        document.getElementById('responseMessage').innerHTML = `Verification Success: ${data.message}`;
    })
    .catch(error => {
        // Handle error response
        document.getElementById('responseMessage').innerHTML = `Error: ${error.message}`;
    });
});
