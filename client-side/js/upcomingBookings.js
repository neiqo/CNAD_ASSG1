document.addEventListener('DOMContentLoaded', () => {
    const upcomingBookingsContainer = document.getElementById('upcoming-bookings-container');

    // Get userID from localStorage
    const userDetails = JSON.parse(localStorage.getItem('userDetails'));
    
    const userID = userDetails?.user_id; // Safely accessing user_id

    if (!userID) {
        upcomingBookingsContainer.innerHTML = '<p>User is not logged in. Please log in to view upcoming bookings.</p>';
        return;
    }

    // Fetch upcoming bookings from the API
    fetch(`http://localhost:5002/api/v1/upcoming-bookings?userID=${userID}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(errorData => {
                throw new Error(`Error: ${errorData.error || 'An unexpected error occurred'}`);
            });
        }
        return response.json();
    })
    .then(data => {
        // Check if the response is empty
        if (!data || data.length === 0) {
            upcomingBookingsContainer.innerHTML = '<p>No upcoming bookings found.</p>';
            return;
        }

        // Create a list of upcoming bookings
        upcomingBookingsContainer.innerHTML = ''; // Clear loading message

        data.forEach(booking => {
            const bookingCard = document.createElement('div');
            bookingCard.classList.add('booking-card');
            const startTime = new Date(booking.startTime);
            const endTime = new Date(booking.endTime);

            bookingCard.innerHTML = `
                <h3>Booking ID: ${booking.bookingID}</h3>
                <p>Vehicle ID: ${booking.vehicleID}</p>
                <p>License Plate: ${booking.licensePlate}</p>
                <p>Model: ${booking.model}</p>
                <p>Rental Rate: $${booking.rentalRate}</p>
                <p>Start Time: ${startTime.toLocaleString()}</p>
                <p>End Time: ${endTime.toLocaleString()}</p>
                <p>Status: ${booking.status}</p>
                <button class="cancel-button" data-booking-id="${booking.bookingID}">Cancel Booking</button>
            `;

            upcomingBookingsContainer.appendChild(bookingCard);
        });

        // Add event listeners to cancel buttons
        const cancelButtons = document.querySelectorAll('.cancel-button');
        cancelButtons.forEach(button => {
            button.addEventListener('click', (event) => {
                const bookingID = event.target.dataset.bookingId;

                // Send PUT request to cancel the booking
                fetch(`http://localhost:5002/api/v1/cancel-booking?bookingID=${bookingID}&userID=${userID}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                })
                .then(response => response.json())
                .then(data => {
                    if (data.message === 'Booking successfully cancelled') {
                        alert('Booking cancelled successfully');
                        // Optionally, you can remove the booking from the list after cancellation
                        event.target.closest('.booking-card').remove();
                    } else {
                        alert('Error cancelling booking: ' + data.error);
                    }
                })
                .catch(error => {
                    console.error('Error cancelling booking:', error);
                    alert('Failed to cancel booking. Please try again later.');
                });
            });
        });

    })
    .catch(error => {
        console.error('Error fetching upcoming bookings:', error);
        upcomingBookingsContainer.innerHTML = `<p>${error.message}</p>`; // Display error message
    });
});
