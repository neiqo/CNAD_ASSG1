window.addEventListener('DOMContentLoaded', function() {
    // Fetch the user details from localStorage (which should be set after login)
    const userDetails = JSON.parse(localStorage.getItem('userDetails'));

    console.log(userDetails);

    // If userDetails are not in localStorage, redirect to login page
    if (!userDetails) {
        window.location.href = 'login.html';
        return;
    }

    // Get user details from the stored object
    const { name, email, contact_no, membership_tier } = userDetails;

    // Set user details to the HTML elements
    document.getElementById('userName').textContent = name;
    document.getElementById('userEmail').textContent = email;
    document.getElementById('userContact').textContent = contact_no;
    document.getElementById('userMembership').textContent = membership_tier;

    // Fetch the member benefits based on the membership tier
    fetch(`http://localhost:5003/api/v1/member-benefits/${membership_tier}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => response.json())
    .then(benefits => {
        // Check if there is an error or empty response
        if (benefits.error) {
            document.getElementById('benefits').textContent = benefits.error;
            return;
        }

        // Display the fetched benefits
        const benefitsContainer = document.getElementById('benefits');
        benefits.forEach(benefit => {
            console.log(benefit)
            const benefitElement = document.createElement('div');
            benefitElement.innerHTML = `
                <strong>${benefit.name}</strong>: <br>   - ${benefit.description}
            `;
            benefitsContainer.appendChild(benefitElement);
        });
    })
    .catch(error => {
        console.error(error);
        document.getElementById('benefits').textContent = 'Error fetching member benefits';
    });
});
