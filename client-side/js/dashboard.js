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
      const userDetails = JSON.parse(localStorage.getItem('userDetails'));
  
      console.log(userDetails);
  
      if (!userDetails) {
          window.location.href = 'login.html';
          return;
      }
  
      const { name, email, contact_no, membership_tier } = userDetails;
  
      document.getElementById('userName').textContent = name;
      document.getElementById('userEmail').textContent = email;
      document.getElementById('userContact').textContent = contact_no;
      document.getElementById('userMembership').textContent = membership_tier;
  
      document.getElementById('newName').value = name;
      document.getElementById('newContact').value = contact_no;
  
      document.getElementById('updateUserForm').addEventListener('submit', async function(event) {
          event.preventDefault(); 
  
          const updatedName = document.getElementById('newName').value;
          const updatedContact = document.getElementById('newContact').value;
          const updatedPassword = document.getElementById('newPassword').value;
  
          let hashedPassword = null;
          if (updatedPassword) {
              hashedPassword = await hashPassword(updatedPassword);
          }
  
          const updateData = {
              name: updatedName,
              contactNo: updatedContact,
          };
  
          if (hashedPassword) {
              updateData.password = hashedPassword; 
          }
  
          fetch(`http://localhost:5001/api/v1/user/${email}`, {
              method: 'PUT',
              headers: {
                  'Content-Type': 'application/json',
              },
              body: JSON.stringify(updateData),
          })
          .then(response => response.json())
          .then(data => {
              if (data.error) {
                  alert(`Error: ${data.error}`);
                  return;
              }
              alert('User details updated successfully');
  
              document.getElementById('userName').textContent = updatedName;
              document.getElementById('userContact').textContent = updatedContact;
  
              userDetails.name = updatedName;
              userDetails.contact_no = updatedContact;
              localStorage.setItem('userDetails', JSON.stringify(userDetails));
          })
          .catch(error => {
              console.error('Error:', error);
              alert('Error updating user details');
          });
      });
  
      fetch(`http://localhost:5003/api/v1/member-benefits/${membership_tier}`, {
          method: 'GET',
          headers: {
              'Content-Type': 'application/json',
          },
      })
      .then(response => response.json())
      .then(benefits => {
          if (benefits.error) {
              document.getElementById('benefits').textContent = benefits.error;
              return;
          }
  
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
  