document.addEventListener('DOMContentLoaded', function () {
    // Get userID from localStorage
    const userDetails = JSON.parse(localStorage.getItem('userDetails'));
    
    userId = userDetails.user_id
    
    if (!userId) {
        alert('User ID not found. Please log in.');
        return;
    }

    // Get the booking history from the server
    fetchBookingHistory(userId);
});

function fetchBookingHistory(userId) {
    console.log(userId)
    const url = `http://localhost:5002/api/v1/past-bookings?userID=${userId}`;

    fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            // Check if the response contains an error or is empty
            if (!data || data.length === 0) {
                document.getElementById("booking-history-body").innerHTML = "<tr><td colspan='5'>No past bookings found.</td></tr>";
                return;
            }

            // Display the fetched bookings
            data.forEach(booking => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${booking.bookingID}</td>
                    <td>${booking.vehicleID}</td>
                    <td>${new Date(booking.startTime).toLocaleString()}</td>
                    <td>${new Date(booking.endTime).toLocaleString()}</td>
                    <td>${booking.status}</td>
                `;
                document.getElementById("booking-history-body").appendChild(row);
            });
        })
        .catch(error => {
            console.error('Error fetching bookings:', error);
            document.getElementById("booking-history-body").innerHTML = "<tr><td colspan='5'>Failed to load booking history. Please try again later.</td></tr>";
        });
}
