document.addEventListener("DOMContentLoaded", () => {
    const vehicleDetailsContainer = document.getElementById("vehicleDetails");
    const bookingForm = document.getElementById("bookingForm");
    const bookingDateInput = document.getElementById("bookingDate");
    const timeSlotsContainer = document.getElementById("timeSlots");
    const bookingErrorDiv = document.getElementById("bookingError");
    const bookingSuccessDiv = document.getElementById("bookingSuccess");
    const promotionsContainer = document.getElementById("promotions"); // Container for displaying promotions

    const params = new URLSearchParams(window.location.search);
    const vehicleID = params.get("vehicleID");

    if (!vehicleID) {
        vehicleDetailsContainer.innerHTML = "<p>Error: Vehicle ID not provided.</p>";
        return;
    }

    // Fetch vehicle details
    fetch(`http://localhost:5002/api/v1/vehicle?vehicleID=${vehicleID}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errorData => {
                    throw new Error(errorData.error || `Error fetching vehicle details. Status: ${response.status}`);
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.error) {
                vehicleDetailsContainer.innerHTML = `<p>${data.error}</p>`;
                return;
            }

            const { vehicle, status } = data;

            vehicleDetailsContainer.innerHTML = `
                <h2>Vehicle Details</h2>
                <p><strong>Model:</strong> ${vehicle.model}</p>
                <p><strong>License Plate:</strong> ${vehicle.licensePlate}</p>
                <p><strong>Rental Rate:</strong> $${vehicle.rentalRate}/hour</p>
                <h3>Current Status</h3>
                ${status ? ` 
                    <p><strong>Location:</strong> ${status.location}</p>
                    <p><strong>Charge Level:</strong> ${status.chargeLevel}%</p>
                    <p><strong>Cleanliness:</strong> ${status.cleanlinessStatus}</p>
                ` : "<p>No status available for this vehicle.</p>"}
            `;

            // Fetch available promotions
            fetch("http://localhost:5003/api/v1/promotions", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            })
                .then(response => {
                    if (!response.ok) {
                        return response.json().then(errorData => {
                            throw new Error(errorData.error || `Error fetching promotions. Status: ${response.status}`);
                        });
                    }
                    return response.json();
                })
                .then(promotions => {
                    if (promotions.error) {
                        promotionsContainer.innerHTML = `<p>${promotions.error}</p>`;
                        return;
                    }

                    promotionsContainer.innerHTML = "<h3>Select a Promotion</h3>";
                    promotions.forEach(promotion => {
                        const promotionOption = document.createElement("label");
                        const promotionCheckbox = document.createElement("input");
                        promotionCheckbox.type = "radio";
                        promotionCheckbox.name = "promotion";
                        promotionCheckbox.value = promotion.promotionID;
                        promotionOption.appendChild(promotionCheckbox);
                        promotionOption.appendChild(document.createTextNode(`${promotion.name} - ${promotion.discount}% off`));
                        promotionsContainer.appendChild(promotionOption);
                        promotionsContainer.appendChild(document.createElement("br"));
                    });
                })
                .catch(error => {
                    console.error("Error fetching promotions:", error);
                    promotionsContainer.innerHTML = `<p>${error.message}</p>`;
                });

            generateTimeSlots();
        })
        .catch(error => {
            console.error("Error fetching vehicle details:", error);
            vehicleDetailsContainer.innerHTML = `<p>${error.message}</p>`;
        });

    function generateTimeSlots() {
        const timeSlotStart = 6; 
        const timeSlotEnd = 22; 
        const timeSlotDuration = 4; 

        for (let hour = timeSlotStart; hour < timeSlotEnd; hour += timeSlotDuration) {
            const label = document.createElement("label");
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.name = "timeSlot";
            checkbox.value = `${hour}:00-${hour + timeSlotDuration -1}:59`;
            label.appendChild(checkbox);
            label.appendChild(document.createTextNode(`${hour}:00 - ${hour + timeSlotDuration}:00`));

            timeSlotsContainer.appendChild(label);
            timeSlotsContainer.appendChild(document.createElement("br"));
        }
    }

    bookingForm.addEventListener("submit", function(event) {
        event.preventDefault();
    
        const selectedSlots = [];
        const checkboxes = document.querySelectorAll('input[name="timeSlot"]:checked');
    
        checkboxes.forEach(checkbox => {
            selectedSlots.push(checkbox.value);
        });
    
        if (selectedSlots.length === 0) {
            bookingErrorDiv.textContent = "Please select at least one time slot.";
            return;
        }
    
        const bookingDate = bookingDateInput.value;
        if (!bookingDate) {
            bookingErrorDiv.textContent = "Please select a date.";
            return;
        }
    
        const slot = selectedSlots[0];
    
        const [start, end] = slot.split("-");
    
        const startTimeStr = `${start.padStart(2, '0')}:00`;  
        const endTimeStr = `${end.padStart(2, '0')}:00`;    

        const startTimeString = `${bookingDate}T${startTimeStr}Z`;  
        const endTimeString = `${bookingDate}T${endTimeStr}Z`;     
    

        const userDetails = JSON.parse(localStorage.getItem('userDetails'));

        const selectedPromotion = document.querySelector('input[name="promotion"]:checked');
        const promotionID = selectedPromotion ? selectedPromotion.value : null;
    
        console.log(promotionID)

        const booking = {
            vehicleID: Number(vehicleID),
            userID: userDetails.user_id,  
            startTime: startTimeString,
            endTime: endTimeString,
            promotionID: Number(promotionID),
        };

        fetch("http://localhost:5002/api/v1/bookings", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(booking)
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errorData => {
                    throw new Error(errorData.error || `Booking failed. Status: ${response.status}`);
                });
            }
            return response.json();
        })
        .then(data => {
            bookingSuccessDiv.textContent = "Booking successful! Payment is pending";
            bookingErrorDiv.textContent = "";
        })
        .catch(error => {
            bookingErrorDiv.textContent = `Error: ${error.message}`;
            bookingSuccessDiv.textContent = "";
        });
    });
});
