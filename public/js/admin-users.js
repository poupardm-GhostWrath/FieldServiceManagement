// =========================================
// User Management Module
// Handles user listing, CRUD operations, and modals
// =========================================

import { Session } from './auth/session.js';

// Constants
const USER_API_ENDPOINT = '/users';
const ROLE_CHECKBOX_CLASS = 'role-checkbox'; // For consistent selection

// ===============================
// 1. INITIALIZATION & GUARDS
// ===============================

export function initUsersPage() {
    // Auth Guard
    if (!Session.isAuthenticated()) {
        window.location.href = '/login.html';
        return;
    }

    const roles = Session.getRoles();
    if (!roles.includes('admin')) {
        alert('Access Denied: Admin privileges required.');
        window.location.href = '/dashboards/customer.html';
        return;
    }

    // Load data when DOM is ready
    document.addEventListener('DOMContentLoaded', () => {
        updateUserInfo();
        loadUsers();
        setupEventListeners();
        hideLoadingOverlay();
    });
}

// ===============================
// 2. USER INFO & UI UPDATES
// ===============================

function updateUserInfo() {
    const name = localStorage.getItem('user_name') || 'Admin';
    const roles = Session.getRoles();
    
    const miniUserName = document.getElementById('miniUserName');
    const miniUserRole = document.getElementById('miniUserRole');
    
    if (miniUserName) miniUserName.textContent = name;
    if (miniUserRole) miniUserRole.textContent = roles.join(', ');
}

function hideLoadingOverlay() {
    const loader = document.getElementById('loadingOverlay');
    if (loader) {
        setTimeout(() => { loader.style.display = 'none'; }, 500);
    }
}

// ===============================
// 3. USER DATA FETCHING & RENDERING
// ===============================

/**
 * Fetches all users from the API and renders the table
 */
export async function loadUsers() {
    const tbody = document.getElementById('usersTableBody');
    if (!tbody) return;

    const token = Session.getToken();
    if (!token) {
        console.error("No token found in session");
        return;
    }

    try {
        // TODO: Replace with real API call
        const response = await fetch(USER_API_ENDPOINT, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP Error: ${response.status}`);
        }

        const data = await response.json();
        renderUsersTable(data.users);

    } catch (error) {
        console.error('Failed to load users:', error);
        // Fallback mock data for development
        const mockUsers = getMockUsers();
        renderUsersTable(mockUsers);
    }
}

/**
 * Renders the user table rows from an array of user objects
 */
function renderUsersTable(users) {
    const tbody = document.getElementById('usersTableBody');
    if (!tbody) return;

    tbody.innerHTML = '';

    if (!users || users.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="6" style="text-align:center; padding: 40px; color: #888;">
                    <i class="fa-solid fa-user-slash" style="font-size: 2rem; display: block; margin-bottom: 10px;"></i>
                    No users found. Click "Add User" to create one.
                </td>
            </tr>
        `;
        return;
    }

    users.forEach(user => {
        const row = document.createElement('tr');
        
        // Format Roles as tags
        const roleTags = Array.isArray(user.roles) 
            ? user.roles.map(r => `<span class="role-tag">${r}</span>`).join(' ')
            : '<span class="role-tag">N/A</span>';
        
        const statusClass = user.is_active ? 'status-active' : 'status-inactive';
        const statusText = user.is_active ? 'Active' : 'Inactive';
        const formattedDate = user.created_at 
            ? new Date(user.created_at).toLocaleDateString() 
            : '-';

        row.innerHTML = `
            <td>
                <div style="font-weight: 600;">${escapeHtml(user.first_name)} ${escapeHtml(user.last_name)}</div>
            </td>
            <td>${escapeHtml(user.email)}</td>
            <td>${roleTags}</td>
            <td><span class="status-badge ${statusClass}">${statusText}</span></td>
            <td>${formattedDate}</td>
            <td style="text-align: right;">
                <button class="action-btn-icon edit" onclick="window.openEditUser('${user.id}')" title="Edit">
                    <i class="fa-solid fa-pen-to-square"></i>
                </button>
                <button class="action-btn-icon delete" onclick="window.openDeleteUser('${user.id}', '${escapeHtml(user.firstName)} ${escapeHtml(user.lastName)}')" title="Deactivate">
                    <i class="fa-solid fa-trash-can"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Helper to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Mock data for development (remove in production)
function getMockUsers() {
    return [
        { id: 1, firstName: "Alice", lastName: "Johnson", email: "alice@example.com", roles: ["admin"], isActive: true, createdAt: "2026-01-15" },
        { id: 2, firstName: "Bob", lastName: "Smith", email: "bob@service.com", roles: ["dispatcher", "technician"], isActive: true, createdAt: "2026-02-10" },
        { id: 3, firstName: "Charlie", lastName: "Brown", email: "charlie@test.com", roles: ["customer"], isActive: false, createdAt: "2026-03-05" }
    ];
}

// ===============================
// 4. MODAL HANDLING
// ===============================

/**
 * Opens the Add User modal
 */
export function openUserModal() {
    resetUserForm();
    document.getElementById('modalTitle').innerHTML = '<i class="fa-solid fa-user-plus" style="margin-right: 10px; color: #667eea;"></i>Add New User';
    document.getElementById('pwd-note').textContent = '(Required for new users)';
    document.getElementById('user-password').required = true;
    document.getElementById('email-note').textContent = 'New user accounts require a unique email.';
    document.getElementById('user-active').checked = true;
    openModal('user-modal');
}

/**
 * Opens the Edit User modal with pre-filled data
 */
export async function openEditUser(userId) {
    resetUserForm();
    document.getElementById('user-id').value = userId;
    document.getElementById('modalTitle').innerHTML = '<i class="fa-solid fa-user-pen" style="margin-right: 10px; color: #667eea;"></i>Edit User';
    document.getElementById('pwd-note').textContent = '(Leave blank to keep current)';
    document.getElementById('user-password').required = false;
    document.getElementById('email-note').textContent = 'Email cannot be changed.';
    document.getElementById('user-email').readOnly = true;

    try {
        const token = Session.getToken();
        const response = await fetch(`${USER_API_ENDPOINT}/${userId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) throw new Error('Failed to fetch user details');
        
        const user = await response.json();
        populateUserForm(user);

    } catch (error) {
        console.error('Error loading user details:', error);
        // Fallback to mock data if API fails
        const mockUser = getMockUsers().find(u => u.id == userId) || {};
        populateUserForm(mockUser);
    }

    openModal('user-modal');
}

/**
 * Fills the form with user data
 */
function populateUserForm(user) {
    document.getElementById('user-first-name').value = user.firstName || '';
    document.getElementById('user-last-name').value = user.lastName || '';
    document.getElementById('user-email').value = user.email || '';
    document.getElementById('user-phone').value = user.phone || '';
    document.getElementById('user-active').checked = user.isActive !== false;

    // Check appropriate role checkboxes
    if (Array.isArray(user.roles)) {
        document.querySelectorAll('input[name="roles"]').forEach(cb => {
            cb.checked = user.roles.includes(cb.value);
        });
    }
}

/**
 * Resets the user form to initial state
 */
function resetUserForm() {
    const form = document.getElementById('userForm');
    if (form) form.reset();
    
    document.getElementById('user-id').value = '';
    document.getElementById('user-message').style.display = 'none';
    document.getElementById('user-email').readOnly = false;
}

/**
 * Opens the Delete Confirmation modal
 */
export function openDeleteUser(userId, userName) {
    document.getElementById('delete-user-id').value = userId;
    const messageEl = document.querySelector('#delete-confirm-modal .modal-body p:first-of-type');
    if (messageEl) {
        messageEl.innerHTML = `Are you sure you want to deactivate <strong>${userName}</strong>?`;
    }
    openModal('delete-confirm-modal');
}

/**
 * Generic modal opener with animation support
 */
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('active');
        modal.style.display = 'flex';
    }
}

/**
 * Generic modal closer with animation support
 */
function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('active');
        setTimeout(() => { modal.style.display = 'none'; }, 300);
    }
}

// ===============================
// 5. FORM SUBMISSION & API CALLS
// ===============================

async function handleUserFormSubmit(e) {
    e.preventDefault();

    const token = Session.getToken();
    if (!token) {
        showMessage('user-message', 'No authentication token found.', 'error');
        return;
    }

    const id = document.getElementById('user-id').value;
    const formData = {
        first_name: document.getElementById('user-first-name').value.trim(),
        last_name: document.getElementById('user-last-name').value.trim(),
        email: document.getElementById('user-email').value.trim(),
        phone: document.getElementById('user-phone').value.trim(),
        roles: Array.from(document.querySelectorAll('input[name="roles"]:checked')).map(cb => cb.value),
        active: document.getElementById('user-active').checked
    };

    const password = document.getElementById('user-password').value;
    if (password) formData.password = password;

    // Validation
    if (formData.roles.length === 0) {
        showMessage('user-message', 'Please select at least one role.', 'error');
        return;
    }

    showMessage('user-message', 'Saving changes...', 'loading');

    try {
        const url = id ? `${USER_API_ENDPOINT}/${id}` : USER_API_ENDPOINT;
        const method = id ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`API Error (${response.status}): ${errorText}`);
        }

        showMessage('user-message', 'User saved successfully!', 'success');
        setTimeout(() => {
            closeModal('user-modal');
            loadUsers();
        }, 1000);

    } catch (err) {
        showMessage('user-message', err.message || 'Failed to save user', 'error');
    }
}

/**
 * Handles user deletion (soft delete)
 */
async function handleUserDelete() {
    const userId = document.getElementById('delete-user-id').value;
    if (!userId) return;

    const token = Session.getToken();
    if (!token) {
        showMessage('user-message', 'No authentication token found.', 'error');
        return;
    }

    try {
        const response = await fetch(`${USER_API_ENDPOINT}/${userId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`API Error (${response.status}): ${errorText}`);
        }

        closeModal('delete-confirm-modal');
        showMessage('user-message', 'User deactivated successfully!', 'success');
        setTimeout(() => loadUsers(), 1000);

    } catch (err) {
        console.error('Delete failed:', err);
        alert('Failed to deactivate user. Please try again.');
    }
}

// ===============================
// 6. EVENT LISTENERS SETUP
// ===============================

function setupEventListeners() {
    // ---- PROFILE MODAL EVENTS ----
    document.getElementById('miniUserName')?.addEventListener('click', openProfileModal);
    document.getElementById('close-profile-modal-btn')?.addEventListener('click', closeProfileModal);
    document.getElementById('cancel-profile-btn')?.addEventListener('click', closeProfileModal);

    // Modal close buttons
    bindClickListeners(['close-modal-btn', 'cancel-user-btn'], () => closeModal('user-modal'));
    bindClickListeners(['close-delete-btn', 'cancel-delete-btn'], () => closeModal('delete-confirm-modal'));

    // Form submissions
    const userForm = document.getElementById('userForm');
    if (userForm) userForm.addEventListener('submit', handleUserFormSubmit);

    const deleteConfirmBtn = document.getElementById('confirm-delete-btn');
    if (deleteConfirmBtn) deleteConfirmBtn.addEventListener('click', handleUserDelete);

    // Sidebar toggle
    const sidebar = document.getElementById('sidebar');
    const toggleBtn = document.getElementById('sidebar-toggle');
    const body = document.body;
    
    if (sidebar && toggleBtn && body) {
        // Load saved state
        const isCollapsed = localStorage.getItem('sidebar_collapsed') === 'true';
        if (isCollapsed) {
            body.classList.add('sidebar-collapsed');
            sidebar.classList.add('collapsed');
            const icon = toggleBtn.querySelector('.toggle-icon');
            if (icon) {
                icon.classList.remove('fa-chevron-left');
                icon.classList.add('fa-chevron-right');
            }
        }

        toggleBtn.addEventListener('click', () => {
            body.classList.toggle('sidebar-collapsed');
            sidebar.classList.toggle('collapsed');
            const isNowCollapsed = body.classList.contains('sidebar-collapsed');
            const icon = toggleBtn.querySelector('.toggle-icon');
            
            if (icon) {
                icon.classList.toggle('fa-chevron-left', !isNowCollapsed);
                icon.classList.toggle('fa-chevron-right', isNowCollapsed);
            }
            
            localStorage.setItem('sidebar_collapsed', isNowCollapsed);
        });
    }

    // Logout button
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) logoutBtn.addEventListener('click', () => Session.logout());

    // Search input (placeholder for future filtering)
    const searchInput = document.getElementById('userSearch');
    if (searchInput) {
        searchInput.addEventListener('keyup', debounce((e) => {
            console.log('Searching users:', e.target.value);
            // TODO: Implement real-time filter
        }, 300));
    }

    // Filter dropdowns (placeholders for future filtering)
    const roleFilter = document.getElementById('roleFilter');
    const statusFilter = document.getElementById('statusFilter');
    if (roleFilter) roleFilter.addEventListener('change', () => console.log('Role filter:', roleFilter.value));
    if (statusFilter) statusFilter.addEventListener('change', () => console.log('Status filter:', statusFilter.value));
}

/**
 * Helper to bind click listeners by ID
 */
function bindClickListeners(ids, handler) {
    ids.forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener('click', handler);
    });
}

/**
 * Utility: Show messages/alerts
 */
function showMessage(elementId, message, type) {
    const el = document.getElementById(elementId);
    if (!el) return;

    el.style.display = 'block';
    el.className = `alert alert-${type}`;
    el.textContent = message;
}

/**
 * Utility: Debounce function for search/filter inputs
 */
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
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


// ===============================
// 7. EXPORT FOR WINDOW ACCESS
// ===============================
// Needed because HTML uses onclick handlers on buttons

window.openUserModal = openUserModal;
window.openEditUser = openEditUser;
window.openDeleteUser = openDeleteUser;
window.closeModal = closeModal;

// ===============================
// INITIALIZATION
// ===============================

// Initialize the page when this module loads
initUsersPage();

export default {
    initUsersPage,
    loadUsers,
    openUserModal,
    openEditUser,
    openDeleteUser
};
