// public/js/features/admin-dashboard.js
import { Session } from '../auth/session.js';

const USER_API_ENDPOINT = '/users';

// ==========================================
// INITIALIZATION
// ==========================================
export function initAdminDashboard() {
    setTimeout(() => {
        document.getElementById('loadingOverlay')?.style.setProperty('display', 'none');
    }, 800);

    if (!Session.isAuthenticated() || !Session.hasRole('admin')) {
        window.location.href = '/login.html';
        return;
    }

    const name = localStorage.getItem('user_name');
    const roles = Session.getRoles();

    updateUIText(name, roles);
    filterSidebarByRole(roles);
    setupEventListeners();
    loadUsers();
    updateStats();
}

// ==========================================
// UI UPDATES
// ==========================================

function updateUIText(name, roles) {
    const miniName = document.getElementById('miniUserName');
    const miniRole = document.getElementById('miniUserRole');
    const welcomeName = document.getElementById('welcomeName');

    if (miniName) miniName.textContent = name || 'Guest';
    if (miniRole) miniRole.textContent = roles.join(', ') || 'Administrator';
    if (welcomeName) welcomeName.textContent = name || 'Admin';
}

function filterSidebarByRole(userRoles) {
    document.querySelectorAll('.sidebar-link[data-role]').forEach(link => {
        if (!userRoles.includes(link.getAttribute('data-role'))) {
            link.style.display = 'none';
        }
    });
}

// ==========================================
// EVENT LISTENERS
// ==========================================

function setupEventListeners() {
    // Logout
    document.getElementById('logout-btn')?.addEventListener('click', () => Session.logout());

    // ---- PROFILE MODAL EVENTS ----
    document.getElementById('miniUserName')?.addEventListener('click', openProfileModal);
    document.getElementById('close-profile-modal-btn')?.addEventListener('click', closeProfileModal);
    document.getElementById('cancel-profile-btn')?.addEventListener('click', closeProfileModal);

    // ---- USER MODAL EVENTS ----
    document.getElementById('close-user-modal-btn')?.addEventListener('click', closeUserModal);
    document.getElementById('cancel-user-btn')?.addEventListener('click', closeUserModal);
    document.getElementById('user-form')?.addEventListener('submit', handleUserFormSubmit);

    // ---- SIDEBAR TOGGLE ----
    setupSidebarToggle();
}

function setupSidebarToggle() {
    const sidebar = document.getElementById('sidebar');
    const toggleBtn = document.getElementById('sidebar-toggle');
    const toggleIcon = toggleBtn?.querySelector('.toggle-icon');
    const body = document.body;

    if (sidebar && toggleBtn && toggleIcon && body) {
        const isCollapsed = localStorage.getItem('sidebar_collapsed') === 'true';
        if (isCollapsed) {
            body.classList.add('sidebar-collapsed');
            sidebar.classList.add('collapsed');
            setToggleIcon('right');
        }

        toggleBtn.addEventListener('click', () => {
            body.classList.toggle('sidebar-collapsed');
            sidebar.classList.toggle('collapsed');
            const isNowCollapsed = body.classList.contains('sidebar-collapsed');
            setToggleIcon(isNowCollapsed ? 'right' : 'left');
            localStorage.setItem('sidebar_collapsed', isNowCollapsed);
        });
    }
}

function setToggleIcon(direction) {
    const icon = document.querySelector('.sidebar-toggle .toggle-icon');
    if (!icon) return;
    icon.classList.remove('fa-chevron-left', 'fa-chevron-right');
    icon.classList.add(direction === 'right' ? 'fa-chevron-right' : 'fa-chevron-left');
}

// ==========================================
// PROFILE MODAL (Self-Service)
// ==========================================

function openProfileModal() {
    const modal = document.getElementById('profile-modal');
    if (!modal) return;

    modal.classList.add('active');
    modal.style.display = 'flex';

    // Load current user data into profile form
    const token = Session.getToken();
    fetch('/users/me', {
        headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => {
        document.getElementById('profile-email').value = data.user?.email || '';
        document.getElementById('profile-first-name').value = data.user?.first_name || '';
        document.getElementById('profile-last-name').value = data.user?.last_name || '';
        document.getElementById('profile-phone').value = data.user?.phone || '';
        document.getElementById('profile-password').value = '';
    })
    .catch(() => {
        // Silently fail — form will just be empty
    });
}

function closeProfileModal() {
    const modal = document.getElementById('profile-modal');
    if (!modal) return;

    modal.classList.remove('active');
    setTimeout(() => {
        modal.style.display = 'none';
        document.getElementById('profile-message').style.display = 'none';
    }, 300);
}

// Expose to global scope for profile-manager.js
window.openProfileModal = openProfileModal;
window.closeProfileModal = closeProfileModal;

// ==========================================
// USER MODAL (Admin CRUD)
// ==========================================

export async function openUserModal(userId = null) {

    const token = Session.getToken();
    if (!token) {
        console.error("No token found in session");
        return;
    }

    const modal = document.getElementById('user-modal');
    const form = document.getElementById('user-form');
    const title = document.getElementById('modal-title');
    const passwordField = document.getElementById('password-field');

    if (!modal) return;

    modal.classList.add('active');
    modal.style.display = 'flex';

    if (form) form.reset();
    document.getElementById('user-id').value = '';

    if (userId) {
        // EDIT MODE
        title.textContent = 'Edit User';
        document.getElementById('user-id').value = userId;
        passwordField.style.display = 'none';

        const response = await fetch(`${USER_API_ENDPOINT}/${userId}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`API Error (${response.status}): ${errorText}`);
        }

        const data = await response.json();

        const firstName = data.user.first_name || '';
        const lastName = data.user.last_name || '';
        const email = data.user.email || '';
        const phone = data.user.phone || '';


        document.getElementById('user-first-name').value = firstName;
        document.getElementById('user-last-name').value = lastName;
        document.getElementById('user-email').value = email;
        document.getElementById('user-phone').value = phone;
        document.querySelectorAll('.role-checkbox').forEach(cb => {
            cb.checked = cb.value === 'customer';
        });
    } else {
        // CREATE MODE
        title.textContent = 'Create New User';
        passwordField.style.display = 'block';
        document.querySelectorAll('.role-checkbox').forEach(cb => cb.checked = false);
        document.getElementById('user-active').checked = true;
    }
}

export function closeUserModal() {
    const modal = document.getElementById('user-modal');
    if (!modal) return;

    modal.classList.remove('active');
    setTimeout(() => {
        modal.style.display = 'none';
        document.getElementById('user-message').style.display = 'none';
    }, 300);
}

// Expose to global scope for HTML onclick handlers
window.openCreateUserModal = () => openUserModal();
window.closeUserModal = closeUserModal;
window.editUser = (id) => openUserModal(id);
window.deleteUser = async (id) => {
    if (!confirm('Are you sure you want to delete this user?')) return;
    try {
        // TODO: Real API call - DELETE /users/{id}
        alert('User deleted (Demo mode)');
        loadUsers();
    } catch (e) {
        alert('Error deleting user');
    }
};

// ==========================================
// FORM SUBMISSION (User CRUD)
// ==========================================

async function handleUserFormSubmit(e) {
    e.preventDefault();

    const token = Session.getToken();
    if (!token) {
        console.error("No token found in session");
        return;
    }

    const id = document.getElementById('user-id').value;
    const formData = {
        first_name: document.getElementById('user-first-name').value,
        last_name: document.getElementById('user-last-name').value,
        email: document.getElementById('user-email').value,
        phone: document.getElementById('user-phone').value,
        roles: Array.from(document.querySelectorAll('.role-checkbox:checked')).map(cb => cb.value),
        active: document.getElementById('user-active').checked
    };

    const pwd = document.getElementById('user-password').value;
    if (pwd) formData.password = pwd;

    const msgEl = document.getElementById('user-message');
    showMsg(msgEl, 'Saving...', 'loading');

    try {
        // TODO: Real API call
        // const res = await fetch(id ? `${USER_API_ENDPOINT}/${id}` : USER_API_ENDPOINT, { ... });
        if (id) {
            const response = await fetch(`${USER_API_ENDPOINT}/${id}`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ formData })
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`API Error (${response.status}): ${errorText}`);
            }
        } else {
            const response = await fetch(USER_API_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ formData })
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`API Error (${response.status}): ${errorText}`);
            }
        }
        showMsg(msgEl, 'User saved successfully!', 'success');
        setTimeout(closeUserModal, 1000);
        loadUsers();
    } catch (err) {
        showMsg(msgEl, err.message || 'Failed to save user', 'error');
    }
}

// ==========================================
// DATA LOADING
// ==========================================

async function updateStats() {
    // TODO: Replace with GET /stats/admin
    document.getElementById('totalUsers').textContent = '47';
    document.getElementById('activeUsers').textContent = '12';
    document.getElementById('pendingApprovals').textContent = '3';
    document.getElementById('systemAlerts').textContent = '2';
}

async function loadUsers() {
    const tbody = document.getElementById('users-table-body');
    if (!tbody) return;

    const token = Session.getToken();

    if (!token) {
        console.error("No token found in session");
        return;
    }

    try {
        const response = await fetch(USER_API_ENDPOINT, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`API Error (${response.status}): ${errorText}`);
        }

        const data = await response.json();

        renderUserTable(data.users);
    } catch (error) {
        tbody.innerHTML = `<tr><td colspan="4" style="text-align:center; color:red;">Error loading users</td></tr>`;
    }
}

function renderUserTable(users) {
    const tbody = document.getElementById('users-table-body');
    if (!tbody) return;

    if (!users || users.length === 0) {
        tbody.innerHTML = `<tr><td colspan="4" style="text-align:center; color:#888;">No users found.</td></tr>`;
        return;
    }

    tbody.innerHTML = users.map(user => {

        const firstName = user.first_name || '';
        const lastName = user.last_name || '';
        const email = user.email || '';
        const phone = user.phone || '';
        const roles = Array.isArray(user.roles) ? user.roles.join(', ') : 'No Role';
        const userId = user.id;
  
        return `
            <tr>
                <td><strong>${firstName} ${lastName}</strong></td>
                <td>${email}</td>
                <td><span class="role-badge">${roles}</span></td>
                <td>
                    <button class="btn-small edit-btn" onclick="window.editUser('${userId}')"><i class="fa-solid fa-edit"></i></button>
                    <button class="btn-small delete-btn" style="background:#ef4444;" onclick="window.deleteUser('${userId}')"><i class="fa-solid fa-trash"></i></button>
                </td>
            </tr>
        `;
    }).join('');
}

// ==========================================
// UTILITY
// ==========================================

function showMsg(el, text, type) {
    if (!el) return;
    el.style.display = 'block';
    if (type === 'loading') {
        el.className = 'alert';
        el.style.backgroundColor = '#f0f0f0';
    } else if (type === 'success') {
        el.className = 'alert alert-success';
    } else {
        el.className = 'alert alert-error';
    }
    el.textContent = text;
}

// ==========================================
// BOOT
// ==========================================
initAdminDashboard();
