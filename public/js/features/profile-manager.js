// public/js/features/profile-manager.js
import { Session } from '../auth/session.js';

const API_BASE = ''; // Relative path

let isOpen = false;

export function openProfileModal() {
    if (isOpen) return;
    isOpen = true;
    
    const modal = document.getElementById('profile-modal');
    modal.classList.add('active');
    modal.style.display = 'flex'; // Ensure flex is applied
    
    loadProfileData();
}

export function closeProfileModal() {
    const modal = document.getElementById('profile-modal');
    modal.classList.remove('active');
    
    setTimeout(() => {
        modal.style.display = 'none';
        clearMessages();
        isOpen = false;
    }, 300); // Wait for animation
}

async function loadProfileData() {
    const token = Session.getToken();
    const msgEl = document.getElementById('profile-message');
    
    if (!token) {
        showMessage("Not authenticated", "error");
        return;
    }

    try {
        // Fetch current user data: GET /users/me
        const response = await fetch(`${API_BASE}/users/me`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) throw new Error("Failed to load profile");
        
        const data = await response.json();
        
        // Fill inputs
        document.getElementById('profile-email').value = data.user?.email || '';
        document.getElementById('profile-first-name').value = data.user?.first_name || data.user?.firstName || '';
        document.getElementById('profile-last-name').value = data.user?.last_name || data.user?.lastName || '';
        document.getElementById('profile-phone').value = data.user?.phone || '';
        document.getElementById('profile-password').value = ''; // Clear password field

    } catch (error) {
        console.error(error);
        showMessage("Could not load profile data.", "error");
    }
}

document.getElementById('profile-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const token = Session.getToken();
    const firstName = document.getElementById('profile-first-name').value.trim();
    const lastName = document.getElementById('profile-last-name').value.trim();
    const phone = document.getElementById('profile-phone').value.trim();
    const newPassword = document.getElementById('profile-password').value;
    
    if (!firstName || !lastName) {
        showMessage("First and Last name are required.", "error");
        return;
    }

    const payload = {
        first_name: firstName,
        last_name: lastName,
        phone: phone
    };

    // Only include password if user typed something
    if (newPassword) {
        payload.password = newPassword;
    }

    const msgEl = document.getElementById('profile-message');
    msgEl.style.display = 'block';
    msgEl.className = 'alert alert-success'; // Default to success while loading? No, wait.
    msgEl.textContent = "Saving...";
    msgEl.className = 'alert'; // Remove specific classes
    msgEl.style.backgroundColor = '#f0f0f0';

    try {
        // Save data: PUT /users/me
        const response = await fetch(`${API_BASE}/users/me`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            const result = await response.json();
            
            // Update Local Storage immediately
            const newName = `${firstName} ${lastName}`;
            localStorage.setItem('user_name', newName);
            localStorage.setItem('token', result.token)
            
            document.getElementById('miniUserName').textContent = newName;
            document.getElementById('welcomeName').textContent = newName;
            
            showMessage("Profile updated successfully!", "success");
            setTimeout(closeProfileModal, 1500);
        } else {
            const err = await response.json().catch(() => ({ message: 'Update failed' }));
            showMessage(err.message || "Failed to update profile.", "error");
        }
    } catch (error) {
        console.error(error);
        showMessage("Network error. Please try again.", "error");
    }
});

function showMessage(msg, type) {
    const el = document.getElementById('profile-message');
    el.style.display = 'block';
    el.textContent = msg;
    
    if (type === 'success') {
        el.style.color = '#27ae60';
        el.style.backgroundColor = '#e6ffe6';
    } else {
        el.style.color = '#c0392b';
        el.style.backgroundColor = '#ffe6e6';
    }
}

function clearMessages() {
    const el = document.getElementById('profile-message');
    el.style.display = 'none';
}

window.openProfileModal = openProfileModal;
window.closeProfileModal = closeProfileModal;
