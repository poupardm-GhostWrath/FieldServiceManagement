// public/js/features/invoice-viewer.js
import { Session } from '../auth/session.js';

export function loadInvoices() {
    const container = document.getElementById('invoice-list-container');
    if (!container) return;

    const token = Session.getToken();

    fetch('/api/my-invoices', {
        headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => {
        renderInvoices(data.invoices);
    })
    .catch(err => {
        console.error('Invoice error:', err);
        container.innerHTML = '<p class="error-text">Unable to load invoices.</p>';
    });
}

function renderInvoices(invoices) {
    const container = document.getElementById('invoice-list-container');

    if (!invoices || invoices.length === 0) {
        container.innerHTML = '<p>No invoices available at this time.</p>';
        return;
    }

    const html = invoices.map(inv => `
        <div class="card invoice-card">
            <div class="invoice-header">
                <h3>Invoice #${inv.invoiceNumber}</h3>
                <small>${new Date(inv.date).toLocaleDateString()}</small>
            </div>
            <div class="invoice-details">
                <p><strong>Amount:</strong> $${(inv.amount).toFixed(2)}</p>
                <p><strong>Status:</strong> 
                    <span class="status-badge ${inv.paid ? 'status-completed' : 'status-pending'}">
                        ${inv.paid ? 'Paid' : 'Unpaid'}
                    </span>
                </p>
            </div>
            <div class="card-footer">
                <a href="${inv.downloadUrl}" target="_blank" class="btn btn-primary btn-small">
                    Download PDF
                </a>
            </div>
        </div>
    `).join('');

    container.innerHTML = html;
}

// Initialize
document.addEventListener('DOMContentLoaded', loadInvoices);
