// Main JavaScript functionality

// Tab switching
document.addEventListener('DOMContentLoaded', () => {
    // Handle tabs
    const tabs = document.querySelectorAll('.tab');
    const tabContents = document.querySelectorAll('.tab-content');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const targetTab = tab.dataset.tab;

            // Remove active class from all tabs and contents
            tabs.forEach(t => t.classList.remove('active'));
            tabContents.forEach(tc => tc.classList.remove('active'));

            // Add active class to clicked tab and corresponding content
            tab.classList.add('active');
            document.getElementById(`${targetTab}-tab`).classList.add('active');

            // Update preview if switching to preview tab
            if (targetTab === 'preview') {
                updatePreview();
            }
        });
    });

    // Handle form submissions with hashcash
    const createForm = document.getElementById('create-form');
    if (createForm) {
        createForm.addEventListener('submit', handleFormSubmit);
    }

    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', handleFormSubmit);
    }

    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleFormSubmit);
    }

    // Live preview for content textarea
    const contentTextarea = document.getElementById('content');
    if (contentTextarea) {
        contentTextarea.addEventListener('input', debounce(updatePreview, 500));
    }
});

// Update markdown preview
async function updatePreview() {
    const content = document.getElementById('content');
    const preview = document.getElementById('preview-content');
    
    if (!content || !preview) return;

    const text = content.value;
    if (!text.trim()) {
        preview.innerHTML = '<p class="text-muted">Enter content in the Text tab to see preview</p>';
        return;
    }

    // Simple markdown rendering (client-side)
    // For production, you might want to use a proper markdown library
    preview.innerHTML = renderMarkdown(text);
}

// Simple markdown renderer (basic implementation)
function renderMarkdown(text) {
    let html = text;

    // Escape HTML
    html = html.replace(/&/g, '&amp;')
               .replace(/</g, '&lt;')
               .replace(/>/g, '&gt;');

    // Headers
    html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>');
    html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>');
    html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>');

    // Bold
    html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');

    // Italic
    html = html.replace(/\*(.*?)\*/g, '<em>$1</em>');

    // Code (inline)
    html = html.replace(/`(.*?)`/g, '<code>$1</code>');

    // Links
    html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>');

    // Line breaks
    html = html.replace(/\n/g, '<br>');

    return html;
}

// Handle form submission with hashcash
async function handleFormSubmit(e) {
    e.preventDefault();
    
    const form = e.target;
    const submitBtn = form.querySelector('#submit-btn');
    const btnText = submitBtn.querySelector('.btn-text');
    const btnLoading = submitBtn.querySelector('.btn-loading');
    const hashcashInput = form.querySelector('#hashcash');

    // Disable submit button
    submitBtn.disabled = true;
    if (btnText && btnLoading) {
        btnText.style.display = 'none';
        btnLoading.style.display = 'inline';
    }

    try {
        // Determine resource based on form action
        let resource = 'create';
        if (form.action.includes('/register')) {
            resource = 'register';
        } else if (form.action.includes('/login')) {
            resource = 'login';
        }

        // Get hashcash challenge
        const response = await fetch(`/hashcash?resource=${resource}`);
        const data = await response.json();
        const challenge = data.challenge;

        // Solve hashcash
        const worker = new HashcashWorker(20);
        const solution = await worker.solve(challenge);

        if (!solution) {
            throw new Error('Failed to solve hashcash');
        }

        // Set hashcash value
        hashcashInput.value = solution;

        // Submit form
        if (form.id === 'create-form') {
            // Handle create form specially to show modal
            const formData = new FormData(form);
            const response = await fetch(form.action, {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                const result = await response.json();
                showSuccessModal(result.url, result.edit_code);
            } else {
                const errorText = await response.text();
                alert('Error: ' + errorText);
            }
        } else {
            // Regular form submission
            form.submit();
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred: ' + error.message);
    } finally {
        // Re-enable submit button
        submitBtn.disabled = false;
        if (btnText && btnLoading) {
            btnText.style.display = 'inline';
            btnLoading.style.display = 'none';
        }
    }
}

// Show success modal
function showSuccessModal(url, editCode) {
    const modal = document.getElementById('success-modal');
    const urlLink = document.getElementById('entry-url');
    const editCodeDisplay = document.getElementById('edit-code-display');

    urlLink.href = '/' + url;
    urlLink.textContent = window.location.origin + '/' + url;
    editCodeDisplay.textContent = editCode;

    modal.style.display = 'flex';
}

// Debounce helper
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

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    // Ctrl/Cmd + Enter to submit form
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        const form = document.querySelector('form');
        if (form) {
            form.dispatchEvent(new Event('submit', { cancelable: true, bubbles: true }));
        }
    }

    // Ctrl/Cmd + P to toggle preview
    if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
        e.preventDefault();
        const previewTab = document.querySelector('[data-tab="preview"]');
        if (previewTab) {
            previewTab.click();
        }
    }
});
