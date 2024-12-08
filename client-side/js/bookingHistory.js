document.addEventListener('DOMContentLoaded', function () {
    const userDetails = JSON.parse(localStorage.getItem('userDetails'));
    
    userId = userDetails.user_id
    
    if (!userId) {
        alert('User ID not found. Please log in.');
        return;
    }

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
            if (!data || data.length === 0) {
                document.getElementById("booking-history-body").innerHTML = "<tr><td colspan='5'>No past bookings found.</td></tr>";
                return;
            }

            data.forEach(booking => {

            const startTime = new Date(booking.startTime);
            const endTime = new Date(booking.endTime);

            const options = { timeZone: 'UTC', hour12: false };
            const startTimeString = startTime.toLocaleString('en-US', options);
            const endTimeString = endTime.toLocaleString('en-US', options);

                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${booking.bookingID}</td>
                    <td>${booking.vehicleID}</td>
                    <td>${startTimeString}</td>
                    <td>${endTimeString}</td>
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
