// Configuration
const API_BASE = ''; 

// Elements
const sidebar = document.getElementById('sidebar');
const sidebarToggle = document.getElementById('sidebarToggle');
const appLayout = document.getElementById('appLayout');
const loadingOverlay = document.getElementById('loadingOverlay');
const userNameEl = document.getElementById('userName');
const welcomeNameEl = document.getElementById('welcomeName');
const userRoleEl = document.getElementById('userRole');
const logoutBtn = document.getElementById('logoutBtn');

// Toggle Functionality
sidebarToggle.addEventListener('click', () => {
  // Toggle 'collapsed' class on appLayout and sidebar
  appLayout.classList.toggle('collapsed');
  sidebar.classList.toggle('collapsed');
  
  // Optional: Save preference to localStorage
  const isCollapsed = appLayout.classList.contains('collapsed');
  localStorage.setItem('sidebarCollapsed', isCollapsed);
});

// Restore Preference on Load
document.addEventListener('DOMContentLoaded', async () => {
  const isCollapsed = localStorage.getItem('sidebarCollapsed') === 'true';
  if (isCollapsed) {
    appLayout.classList.add('collapsed');
    sidebar.classList.add('collapsed');
  }

  // Check Auth
  const token = localStorage.getItem('authToken');
  if (!token) {
    window.location.href = 'index.html';
    return;
  }

  try {
    // Fetch User Profile
    const res = await fetch('/users/me', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (!res.ok) throw new Error('Auth failed');
    
    const data = await res.json();
    const name = `${data.user.first_name || data.user.firstName} ${data.user.last_name || data.user.lastName}`;
    
    userNameEl.textContent = name.split(' ')[0];
    welcomeNameEl.textContent = name.split(' ')[0];
    userRoleEl.textContent = data.user.roles?.[0] || 'User';

    // Filter Nav Items (Simple version)
    filterNavByRole(data.user.roles);

  } catch (err) {
    console.error(err);
    localStorage.removeItem('authToken');
    window.location.href = 'index.html';
  } finally {
    loadingOverlay.style.display = 'none';
  }
});

function filterNavByRole(roles) {
  const navItems = document.querySelectorAll('.nav-item[data-roles]');
  navItems.forEach(item => {
    const allowed = item.getAttribute('data-roles').split(' ');
    const hasAccess = roles.some(r => allowed.includes(r.toLowerCase()));
    item.style.display = hasAccess ? '' : 'none';
  });
}

// Logout
logoutBtn.addEventListener('click', () => {
  if(confirm('Log out?')) {
    localStorage.removeItem('authToken');
    window.location.href = 'index.html';
  }
});
