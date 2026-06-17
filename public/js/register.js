// register.js

document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('registerForm');
    const submitBtn = document.getElementById('submitBtn');

    // Regex Patterns
    const emailRegex = /[A-Za-z0-9]([A-Za-z0-9._%+-]*[A-Za-z0-9])?@[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?\.[A-Za-z]{2,}/;
    const passRegex = /[A-Za-z0-9@$!%*?&]{12,64}/;

    // Validation Functions
    function validateEmail(email) {
        return emailRegex.test(email);
    }

    function validatePassword(password) {
        if (!passRegex.test(password)) return false;
        
        // Check for specific character types
        const hasUpper = /[A-Z]/.test(password);
        const hasLower = /[a-z]/.test(password);
        const hasDigit = /[0-9]/.test(password);
        const hasSpecial = /[@$!%*?&]/.test(password);

        return hasUpper && hasLower && hasDigit && hasSpecial;
    }

    function showError(inputId, message) {
        const errorEl = document.getElementById(`error-${inputId}`);
        const successEl = document.getElementById(`success-${inputId}`);
        const inputEl = document.getElementById(inputId);

        if (errorEl) {
            errorEl.textContent = message;
            errorEl.style.display = 'block';
        }
        if (successEl) successEl.style.display = 'none';
        
        if (inputEl) {
            inputEl.classList.add('error');
            inputEl.classList.remove('success');
        }
    }

    function showSuccess(inputId) {
        const errorEl = document.getElementById(`error-${inputId}`);
        const successEl = document.getElementById(`success-${inputId}`);
        const inputEl = document.getElementById(inputId);

        if (errorEl) errorEl.style.display = 'none';
        if (successEl) successEl.style.display = 'block';
        
        if (inputEl) {
            inputEl.classList.remove('error');
            inputEl.classList.add('success');
        }
    }

    function clearErrors() {
        ['firstName', 'lastName', 'email', 'phone', 'password'].forEach(id => {
            const errorEl = document.getElementById(`error-${id}`);
            const successEl = document.getElementById(`success-${id}`);
            const inputEl = document.getElementById(id);
            
            if (errorEl) errorEl.style.display = 'none';
            if (successEl) successEl.style.display = 'none';
            if (inputEl) {
                inputEl.classList.remove('error');
                inputEl.classList.remove('success');
            }
        });
    }

    // Real-time validation listeners
    const inputs = ['firstName', 'lastName', 'email', 'password'];
    inputs.forEach(id => {
        document.getElementById(id).addEventListener('input', function() {
            if (this.value.trim() === '') {
                // Don't show error immediately on empty, just clear success
                showError(id, ''); 
                this.classList.remove('success');
                return;
            }

            if (id === 'email' && !validateEmail(this.value)) {
                showError(id, 'Please enter a valid email address.');
            } else if (id === 'email') {
                showSuccess(id);
            }

            if (id === 'password' && !validatePassword(this.value)) {
                showError(id, 'Password must be 12-64 chars with uppercase, lowercase, number, and special char.');
            } else if (id === 'password') {
                showSuccess(id);
            }
        });
    });

    // Form Submission
    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        clearErrors();

        const firstName = document.getElementById('firstName').value.trim();
        const lastName = document.getElementById('lastName').value.trim();
        const email = document.getElementById('email').value.trim();
        const phone = document.getElementById('phone').value.trim();
        const password = document.getElementById('password').value;

        let isValid = true;

        // Validate Required Fields
        if (!firstName) { showError('firstName', 'First Name is required'); isValid = false; }
        if (!lastName) { showError('lastName', 'Last Name is required'); isValid = false; }
        if (!email || !validateEmail(email)) { showError('email', 'Valid email is required'); isValid = false; }
        if (!password || !validatePassword(password)) { showError('password', 'Password does not meet requirements'); isValid = false; }

        if (!isValid) return;

        // Prepare Payload
        const payload = {
            first_name: firstName,
            last_name: lastName,
            email: email,
            password: password,
            phone: phone // Will be empty string or null if optional logic needs handling, but usually empty string works
        };

        submitBtn.disabled = true;
        submitBtn.textContent = 'Creating Account...';

        try {
            const response = await fetch('/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            });

            if (response.ok) {
                const data = await response.json();
                alert('Account created successfully! Redirecting to login...');
                window.location.href = 'login.html';
            } else {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || 'Registration failed. Please check your details.');
            }
        } catch (error) {
            alert('Error: ' + error.message);
            submitBtn.disabled = false;
            submitBtn.textContent = 'Create Account';
        }
    });
});
