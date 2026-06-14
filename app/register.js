// Configuration
const API_BASE = '/auth';
const PASSWORD_REGEX = /^[A-Za-z0-9@$!%*?&]{12,64}$/;

// Element References
const registerForm = document.getElementById('registerForm');
const firstNameInput = document.getElementById('firstName');
const lastNameInput = document.getElementById('lastName');
const phoneInput = document.getElementById('phone');
const emailInput = document.getElementById('email');
const passwordInput = document.getElementById('password');
const confirmInput = document.getElementById('confirmPassword');
const submitBtn = document.getElementById('submitBtn');
const apiError = document.getElementById('apiError');
const successMsg = document.getElementById('successMsg');

// Toggle Password Visibility
document.querySelectorAll('.toggle-password').forEach(btn => {
  btn.addEventListener('click', function() {
    const input = this.parentElement.querySelector('input');
    const type = input.type === 'password' ? 'text' : 'password';
    input.type = type;
  });
});

// Validation Functions
function isValidEmail(email) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

function isValidPhone(phone) {
  const re = /^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$/;
  return re.test(phone.replace(/\s/g, ''));
}

function validateForm() {
  let isValid = true;
  clearErrors();

  // First Name
  if (!firstNameInput.value.trim()) {
    showError('firstNameError', 'First name is required');
    isValid = false;
  }

  // Last Name
  if (!lastNameInput.value.trim()) {
    showError('lastNameError', 'Last name is required');
    isValid = false;
  }

  // Email
  if (!emailInput.value.trim()) {
    showError('emailError', 'Email is required');
    isValid = false;
  } else if (!isValidEmail(emailInput.value)) {
    showError('emailError', 'Please enter a valid email');
    isValid = false;
  }

  // Password
  if (!passwordInput.value) {
    showError('passwordError', 'Password is required');
    isValid = false;
  } else if (passwordInput.value.length < 8) {
    showError('passwordError', 'Password must be at least 8 characters');
    isValid = false;
  }

  
  // Confirm Password
  if (!confirmInput.value) {
    showError('confirmPasswordError', 'Please confirm your password');
    isValid = false;
  } else if (confirmInput.value !== passwordInput.value) {
    showError('confirmPasswordError', 'Passwords do not match');
    isValid = false;
  }

  // Phone (optional but validate if provided)
  if (phoneInput.value.trim() && !isValidPhone(phoneInput.value)) {
    showError('phoneError', 'Please enter a valid phone number');
    isValid = false;
  }

  return isValid;
}

function showError(elementId, message) {
  const el = document.getElementById(elementId);
  if (el) el.textContent = message;
}

function clearErrors() {
  document.querySelectorAll('.error-message').forEach(el => el.textContent = '');
}

function showApiError(message) {
  apiError.querySelector('.alert-message').textContent = message;
  apiError.style.display = 'flex';
  setTimeout(() => { apiError.style.display = 'none'; }, 5000);
}

function showSuccess(message) {
  successMsg.querySelector('.alert-message').textContent = message;
  successMsg.style.display = 'flex';
  apiError.style.display = 'none';
  registerForm.reset();
}

// Handle Registration
async function handleRegister(e) {
  e.preventDefault();

  if (!validateForm()) return;

  // UI Loading State
  submitBtn.disabled = true;
  submitBtn.querySelector('.btn-text').style.display = 'none';
  submitBtn.querySelector('.btn-loader').style.display = 'inline';
  apiError.style.display = 'none';
  successMsg.style.display = 'none';

  try {
    const response = await fetch(`${API_BASE}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: emailInput.value.trim(),
        password: passwordInput.value,
        first_name: firstNameInput.value.trim(),
        last_name: lastNameInput.value.trim(),
        phone: phoneInput.value.trim(),
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Registration failed');
    }

    // Success handling
    showSuccess(data.message || 'Account created successfully! Redirecting to login...');
    
    // Redirect after delay
    setTimeout(() => {
      window.location.href = 'index.html';
    }, 2000);

  } catch (error) {
    console.error('Registration error:', error);
    showApiError(error.message || 'Unable to connect to server');
  } finally {
    submitBtn.disabled = false;
    submitBtn.querySelector('.btn-text').style.display = 'inline';
    submitBtn.querySelector('.btn-loader').style.display = 'none';
  }
}

// Real-time Validation
[emailInput, passwordInput, confirmInput].forEach(input => {
  input.addEventListener('input', () => {
    clearErrors();
    apiError.style.display = 'none';

    // Password
    if (input.id === 'password') {
      const val = input.value;
  
      // Check Length
      const lengthValid = val.length >= 12 && val.length <= 64;
      document.getElementById('req-length').className = lengthValid ? 'req-pass' : 'req-fail';
  
      // Check Characters
      const charsValid = PASSWORD_REGEX.test(val);
      document.getElementById('req-chars').className = charsValid ? 'req-pass' : 'req-fail';
  
      // Clear error if valid
      if (charsValid) {
        showError('passwordError', '');
      }
    }
    
    // Check password match dynamically
    if (input.id === 'confirmPassword' && confirmInput.value && passwordInput.value) {
      if (confirmInput.value !== passwordInput.value) {
        showError('confirmPasswordError', 'Passwords do not match');
      } else {
        showError('confirmPasswordError', '');
      }
    }
  });
});

phoneInput?.addEventListener('input', () => {
  apiError.style.display = 'none';
});

// Event Listeners
registerForm.addEventListener('submit', handleRegister);
