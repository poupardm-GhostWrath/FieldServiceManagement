// public/js/features/job-history.js
import { Session } from '../auth/session.js';

const API_BASE = ''; // Relative path since we are on same domain

export function loadJobHistory() {
    const container = document.getElementById('job-history-container');
    if (!container) return;

    const token = Session.getToken();

    if (!token) {
        container.innerHTML = '<p>Error: Authentication token missing.</p>';
        return;
    }

    // Fetch jobs for the current user
    // Note: Backend needs a route like GET /customers/me/jobs or similar
    // For now, let's assume we call GET /users/me (if user is customer) 
    // OR a dedicated endpoint. Let's assume a mock endpoint for demo: /api/my-jobs
    
    fetch(`${API_BASE}/api/my-jobs`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (!response.ok) throw new Error('Failed to fetch jobs');
        return response.json();
    })
    .then(data => {
        renderJobs(data.jobs);
    })
    .catch(error => {
        console.error('Error loading jobs:', error);
        container.innerHTML = `<p class="error-text">Could not load jobs. Please try again later.</p>`;
    });
}

function renderJobs(jobs) {
    const container = document.getElementById('job-history-container');
    
    if (!jobs || jobs.length === 0) {
        container.innerHTML = '<p>No job history found.</p>';
        return;
    }

    const html = jobs.map(job => `
        <div class="card">
            <div class="card-header">
                <h3>${job.title || 'Job #${job.id}'}</h3>
                <span class="status-badge status-${job.status}">${formatStatus(job.status)}</span>
            </div>
            <div class="card-body">
                <p><strong>Date:</strong> ${new Date(job.createdAt).toLocaleDateString()}</p>
                <p><strong>Description:</strong> ${job.description || 'No description'}</p>
                <p><strong>Status:</strong> ${formatStatus(job.status)}</p>
            </div>
            <div class="card-footer">
                <button class="btn btn-small" onclick="alert('View details for Job ID: ${job.id}')">View Details</button>
            </div>
        </div>
    `).join('');

    container.innerHTML = html;
}

function formatStatus(status) {
    const map = {
        'pending': 'Pending',
        'in-progress': 'In Progress',
        'completed': 'Completed',
        'cancelled': 'Cancelled'
    };
    return map[status] || status.toUpperCase();
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', loadJobHistory);
