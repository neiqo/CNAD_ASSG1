// function to encrypt/hash the password using SHA-256
async function hashPassword(password) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
  
    // hash the password using SHA-256
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return hashHex; // return the hashed password as a hex string
  }
  
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
  
      // Prefill the form with the current details
      document.getElementById('newName').value = name;
      document.getElementById('newContact').value = contact_no;
  
      // Handle form submission to update user details
      document.getElementById('updateUserForm').addEventListener('submit', async function(event) {
          event.preventDefault(); // Prevent form from submitting normally
  
          const updatedName = document.getElementById('newName').value;
          const updatedContact = document.getElementById('newContact').value;
          const updatedPassword = document.getElementById('newPassword').value;
  
          // If a new password is provided, hash it
          let hashedPassword = null;
          if (updatedPassword) {
              hashedPassword = await hashPassword(updatedPassword);
          }
  
          // Prepare the data to send in the PUT request
          const updateData = {
              name: updatedName,
              contactNo: updatedContact,
          };
  
          if (hashedPassword) {
              updateData.password = hashedPassword; // Only send password if it's provided
          }
  
          // Send the update request to the user service
          fetch(`http://localhost:5001/api/v1/user/${email}`, {
              method: 'PUT',
              headers: {
                  'Content-Type': 'application/json',
              },
              body: JSON.stringify(updateData),
          })
          .then(response => response.json())
          .then(data => {
              // Handle the response
              if (data.error) {
                  alert(`Error: ${data.error}`);
                  return;
              }
              alert('User details updated successfully');
  
              // Update the details on the page with the new values
              document.getElementById('userName').textContent = updatedName;
              document.getElementById('userContact').textContent = updatedContact;
  
              // Optionally, you can update the userDetails in localStorage
              userDetails.name = updatedName;
              userDetails.contact_no = updatedContact;
              localStorage.setItem('userDetails', JSON.stringify(userDetails));
          })
          .catch(error => {
              console.error('Error:', error);
              alert('Error updating user details');
          });
      });
  
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
              console.log(benefit);
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
  