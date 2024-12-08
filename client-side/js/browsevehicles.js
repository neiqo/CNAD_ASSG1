document.addEventListener("DOMContentLoaded", () => {
    const vehiclesList = document.getElementById("vehiclesList");

    fetch("http://localhost:5002/api/v1/vehicles", {
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
            if (!data || data.length === 0) {
                vehiclesList.innerHTML = "<p>No vehicles available.</p>";
                return;
            }

            data.forEach(vehicle => {
                const vehicleCard = document.createElement("div");
                vehicleCard.classList.add("vehicle-card");

                vehicleCard.innerHTML = `
                    <h3>${vehicle.model}</h3>
                    <p>License Plate: ${vehicle.licensePlate}</p>
                    <p>Rental Rate: $${vehicle.rentalRate}/hour</p>
                    <button onclick="viewDetails(${vehicle.vehicleID})">View Details</button>
                `;

                vehiclesList.appendChild(vehicleCard);
            });
        })
        .catch(error => {
            console.error("Error fetching vehicles:", error);
            vehiclesList.innerHTML = "<p>Failed to load vehicles. Please try again later.</p>";
        });
});

function viewDetails(vehicleID) {
    window.location.href = `vehicledetails.html?vehicleID=${vehicleID}`;
}
