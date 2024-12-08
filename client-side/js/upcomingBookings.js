document.addEventListener('DOMContentLoaded', () => {
    const upcomingBookingsContainer = document.getElementById('upcoming-bookings-container');
    const userDetails = JSON.parse(localStorage.getItem('userDetails'));
    const userID = userDetails?.user_id;

    if (!userID) {
        upcomingBookingsContainer.innerHTML = '<p>User is not logged in. Please log in to view upcoming bookings.</p>';
        return;
    }

    fetch(`http://localhost:5002/api/v1/upcoming-bookings?userID=${userID}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => response.json())
    .then(data => {
        if (!data || data.length === 0) {
            upcomingBookingsContainer.innerHTML = '<p>No upcoming bookings found.</p>';
            return;
        }

        upcomingBookingsContainer.innerHTML = '';

        data.forEach(booking => {
            const bookingCard = document.createElement('div');
            bookingCard.classList.add('booking-card');

            // Parse the time strings into Date objects
            const startTime = new Date(booking.startTime);
            const endTime = new Date(booking.endTime);

            const options = { timeZone: 'UTC', hour12: false };
            const startTimeString = startTime.toLocaleString('en-US', options);
            const endTimeString = endTime.toLocaleString('en-US', options);

            bookingCard.innerHTML = `
                <h3>Booking ID: ${booking.bookingID}</h3>
                <p>Vehicle ID: ${booking.vehicleID}</p>
                <p>License Plate: ${booking.licensePlate}</p>
                <p>Model: ${booking.model}</p>
                <p>Rental Rate: $${booking.rentalRate}</p>
                <p>Start Time (UTC): ${startTimeString}</p>
                <p>End Time (UTC): ${endTimeString}</p>
                <p>Status: ${booking.status}</p>
                <button class="cancel-button" data-booking-id="${booking.bookingID}">Cancel Booking</button>
                <button class="modify-button" data-booking-id="${booking.bookingID}" data-booking-vehicle-id="${booking.vehicleID}">Modify Timeslot</button>
            `;

            if (booking.status === "Pending") {
                bookingCard.innerHTML += `
                    <button class="make-payment-button" data-booking-id="${booking.bookingID}" data-user-id="${userID}">Make Payment</button>
                `;
            }

            upcomingBookingsContainer.appendChild(bookingCard);
        });

        // Cancel booking functionality
        document.querySelectorAll('.cancel-button').forEach(button => {
            button.addEventListener('click', (event) => {
                const bookingID = event.target.dataset.bookingId;
                cancelBooking(bookingID, userID, event.target);
            });
        });

        // Modify timeslot functionality
        document.querySelectorAll('.modify-button').forEach(button => {
            button.addEventListener('click', (event) => {
                const bookingID = event.target.dataset.bookingId;
                const vehicleID = event.target.dataset.bookingVehicleId;
                displayModifyTimeslotModal(bookingID, userID, Number(vehicleID));
            });
        });

        // Make Payment functionality
        document.querySelectorAll('.make-payment-button').forEach(button => {
            button.addEventListener('click', (event) => {
                const bookingID = event.target.dataset.bookingId;
                const userID = event.target.dataset.userId;
                makePayment(bookingID, userID, event.target);
            });
        });

    })
    .catch(error => {
        console.error('Error fetching upcoming bookings:', error);
        upcomingBookingsContainer.innerHTML = `<p>${error.message}</p>`;
    });
});

function cancelBooking(bookingID, userID, cancelButton) {
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
            cancelButton.closest('.booking-card').remove();
        } else {
            alert('Error cancelling booking: ' + data.error);
        }
    })
    .catch(error => {
        console.error('Error cancelling booking:', error);
        alert('Failed to cancel booking. Please try again later.');
    });
}
function makePayment(bookingID, userID, paymentButton) {
    console.log(userID);
    fetch(`http://localhost:5004/api/v1/make-payment/${bookingID}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            userID: Number(userID),
        }),
    })
    .then(response => {
        // Check if the response status is OK (2xx)
        if (!response.ok) {
            // Attempt to read the error message from the response
            return response.text().then(text => {
                throw new Error(`HTTP error! Status: ${response.status} - ${response.statusText}. Response: ${text}`);
            });
        }
        // Attempt to parse the JSON response
        return response.json();
    })
    .then(data => {
            alert('Payment successful and booking status updated to Active');
            paymentButton.closest('.booking-card').remove(); // Remove the payment button
            location.reload();

        
    })
    .catch(error => {
        console.error('Error making payment:', error);
        alert('Failed to make payment. Please try again later. Error: ' + error.message);
    });
}



function displayModifyTimeslotModal(bookingID, userID, vehicleID) {
    const modal = document.createElement('div');
    modal.classList.add('modal');
    modal.innerHTML = `
        <div class="modal-content">
            <h3>Modify Timeslot for Booking ID: ${bookingID}</h3>
            <label for="newDate">Select a New Date:</label>
            <input type="date" id="newDate" required>
            
            <label for="newTimeslot">Select a New Timeslot:</label>
            <select id="newTimeslot" required>
                <option value="6:00-10:00">6:00 - 10:00</option>
                <option value="10:00-14:00">10:00 - 14:00</option>
                <option value="14:00-18:00">14:00 - 18:00</option>
                <option value="18:00-22:00">18:00 - 22:00</option>
            </select>
            
            <button id="confirmModify">Confirm</button>
            <button id="closeModal">Cancel</button>
        </div>
    `;
    document.body.appendChild(modal);

    const today = new Date().toISOString().split('T')[0]; 
    document.getElementById('newDate').setAttribute('min', today);

    document.getElementById('confirmModify').addEventListener('click', () => {
        const newDate = document.getElementById('newDate').value;
        const newTimeslot = document.getElementById('newTimeslot').value;

        if (!newDate || !newTimeslot) {
            alert('Please select both a date and a timeslot.');
            return;
        }

        const [start, end] = newTimeslot.split('-');
        
        const newStartTime = `${newDate}T${start}:00Z`;  
        const newEndTime = `${newDate}T${end}:00Z`;

        fetch(`http://localhost:5002/api/v1/modify-booking?bookingID=${bookingID}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                vehicleID,
                userID,
                newStartTime,
                newEndTime
            }),
        })
        .then(response => response.json())
        .then(data => {
            console.log('Response Data:', data); 
            if (data.message === 'Booking successfully modified') { 
                alert('Booking updated successfully');
                modal.remove();
                location.reload();
            } else {
                alert('Error modifying booking: ' + (data.error || 'Unknown error'));
            }
        })
        .catch(error => {
            console.error('Error modifying booking:', error);
            alert('Failed to modify booking. Please try again later.');
        });
              
    });

    document.getElementById('closeModal').addEventListener('click', () => {
        modal.remove();
    });
}
