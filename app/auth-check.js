// Simple auth check before page renders
(function authCheck() {
  const token = localStorage.getItem('authToken');
  
  if (!token) {
    // Redirect immediately without waiting for DOM ready
    window.location.href = 'index.html';
    return false;
  }
  
  return true;
})();
