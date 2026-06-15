// public/js/auth/guards.js
import { Session } from './session.js';

export const Guards = {
    // Enforce login on page load
    enforceLogin() {
        if (!Session.isAuthenticated()) {
            window.location.href = '/login.html';
        }
    },

    // Ensure user has at least one of these roles, else redirect to home or error page
    requireRoles(allowedRoles, fallbackUrl = '/') {
        if (!Session.hasAnyRole(allowedRoles)) {
            console.warn(`Access denied. User lacks roles: ${allowedRoles.join(', ')}`);
            // Optional: Show alert or redirect
            alert('You do not have permission to access this dashboard.');
            window.location.href = fallbackUrl;
        }
    },

    // Hide an element by ID if user lacks permission
    hideIfNoRole(elementId, requiredRole) {
        const el = document.getElementById(elementId);
        if (el && !Session.hasRole(requiredRole)) {
            el.style.display = 'none';
        }
    },

    // Show element ONLY if user has specific role
    showIfRole(elementId, requiredRole) {
        const el = document.getElementById(elementId);
        if (el && !Session.hasRole(requiredRole)) {
            el.style.display = 'none';
        } else if (el) {
            el.style.display = 'block'; // or whatever default
        }
    },

    // Dynamic Sidebar Menu Builder
    // Takes a list of menu items and filters them based on roles
    buildSidebarMenu(menuItems) {
        const container = document.getElementById('sidebar-menu');
        if (!container) return;

        const myRoles = Session.getRoles();
        
        const filteredItems = menuItems.filter(item => 
            !item.requiredRole || myRoles.includes(item.requiredRole)
        );

        container.innerHTML = filteredItems.map(item => `
            <a href="${item.url}" class="sidebar-link">
                <span class="icon">${item.icon}</span>
                <span>${item.label}</span>
            </a>
        `).join('');
    }
};
