document.addEventListener("DOMContentLoaded", () => {
    const vehicleDetailsContainer = document.getElementById("vehicleDetails");

    // Extract the vehicleID from the URL query parameters
    const params = new URLSearchParams(window.location.search);
    const vehicleID = params.get("vehicleID");

    if (!vehicleID) {
        vehicleDetailsContainer.innerHTML = "<p>Error: Vehicle ID not provided.</p>";
        return;
    }

    // Fetch vehicle details from the backend
    fetch(`http://localhost:5002/api/v1/vehicle?vehicleID=${vehicleID}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.error) {
                vehicleDetailsContainer.innerHTML = `<p>${data.error}</p>`;
                return;
            }

            // Populate the vehicle details on the page
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
        })
        .catch(error => {
            console.error("Error fetching vehicle details:", error);
            vehicleDetailsContainer.innerHTML = "<p>Failed to load vehicle details. Please try again later.</p>";
        });
});
