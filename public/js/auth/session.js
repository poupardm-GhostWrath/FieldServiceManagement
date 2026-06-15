// public/js/auth/session.js

const AUTH_KEYS = {
    TOKEN: 'auth_token',
    ROLES: 'user_roles',
    USER_ID: 'user_id',
    NAME: 'user_name'
};

export const Session = {
    // Store data after successful login
    login(token, user) {
        localStorage.setItem(AUTH_KEYS.TOKEN, token);
        localStorage.setItem(AUTH_KEYS.ROLES, JSON.stringify(user.roles));
        localStorage.setItem(AUTH_KEYS.USER_ID, user.id);
        localStorage.setItem(AUTH_KEYS.NAME, `${user.first_name} ${user.last_name}`);
    },

    // Retrieve stored token
    getToken() {
        return localStorage.getItem(AUTH_KEYS.TOKEN);
    },

    // Retrieve stored roles (returns array or empty array)
    getRoles() {
        const roles = localStorage.getItem(AUTH_KEYS.ROLES);
        return roles ? JSON.parse(roles) : [];
    },

    // Check if user has specific role
    hasRole(role) {
        return this.getRoles().includes(role);
    },

    // Check if user has ANY of the provided roles
    hasAnyRole(rolesList) {
        const myRoles = this.getRoles();
        return rolesList.some(role => myRoles.includes(role));
    },

    // Clear all auth data (Logout)
    logout() {
        Object.values(AUTH_KEYS).forEach(key => localStorage.removeItem(key));
        window.location.href = '/'; // Redirect to home
    },

    // Check if currently logged in
    isAuthenticated() {
        return !!this.getToken();
    }
};
