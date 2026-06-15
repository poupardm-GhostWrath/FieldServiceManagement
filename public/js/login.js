// login.js - Handle Login & Multi-Role Redirect
async function handleLogin(email, password) {
    try {
        const response = await fetch('/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            throw new Error('Invalid credentials');
        }

        const data = await response.json();
        
        // Store auth token and ALL roles
        localStorage.setItem('token', data.token);
        localStorage.setItem('userRoles', JSON.stringify(data.user.roles)); // Array of roles
        localStorage.setItem('userId', data.user.id);
        localStorage.setItem('userName', `${data.user.firstName} ${data.user.lastName}`);

        // MULTI-ROLE REDIRECT LOGIC
        // Priority order: Admin > Dispatcher > Technician > Customer
        const roles = data.user.roles;
        
        if (roles.includes('admin')) {
            window.location.href = '/dashboard/admin';
        } else if (roles.includes('dispatcher')) {
            window.location.href = '/dashboard/dispatcher';
        } else if (roles.includes('technician')) {
            window.location.href = '/dashboard/technician';
        } else if (roles.includes('customer')) {
            window.location.href = '/dashboard/customer';
        } else {
            console.error('No valid roles found:', roles);
            window.location.href = '/';
        }

    } catch (error) {
        alert('Login failed: ' + error.message);
    }
}
