/**
 * GoShort Application Controller & State Engine
 */
class App {
  constructor() {
    this.currentTab = 'home';
    this.currentPage = 1;
    this.pageSize = 10;
    this.totalPages = 1;
    this.totalLinks = 0;
    this.urls = [];
    this.filteredUrls = [];
    this.searchQuery = '';
    this.filterStatus = 'all';
    this.latestCreatedUrl = null;
    this.activeAnalyticsUrlId = null;
    this.activeAnalyticsShortCode = null;
    this.qrCodeInstance = null;

    // Initialize when DOM is ready
    document.addEventListener('DOMContentLoaded', () => this.init());
  }

  init() {
    // Initialize Lucide Icons
    if (window.lucide) {
      window.lucide.createIcons();
    }

    // Check Authentication & Update UI
    this.updateAuthUI();

    // Check API Health
    this.checkHealth(false);
    setInterval(() => this.checkHealth(false), 30000); // Check every 30s

    // Auto-focus on home input
    const urlInput = document.getElementById('original-url-input');
    if (urlInput) urlInput.focus();
  }

  // ==========================================
  // TOAST NOTIFICATION SYSTEM
  // ==========================================
  showToast(message, type = 'info', duration = 3500) {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = 'toast flex items-center gap-3 p-3.5 rounded-xl border shadow-xl backdrop-blur-md text-xs font-medium';

    let iconName = 'info';
    let typeClasses = 'bg-slate-900/95 border-slate-700 text-slate-200';

    if (type === 'success') {
      iconName = 'check-circle-2';
      typeClasses = 'bg-emerald-950/90 border-emerald-500/40 text-emerald-200 shadow-emerald-900/20';
    } else if (type === 'error') {
      iconName = 'alert-circle';
      typeClasses = 'bg-rose-950/90 border-rose-500/40 text-rose-200 shadow-rose-900/20';
    } else if (type === 'warning') {
      iconName = 'alert-triangle';
      typeClasses = 'bg-amber-950/90 border-amber-500/40 text-amber-200 shadow-amber-900/20';
    }

    toast.className += ` ${typeClasses}`;
    toast.innerHTML = `
      <i data-lucide="${iconName}" class="w-4 h-4 shrink-0"></i>
      <span class="flex-1">${message}</span>
      <button onclick="this.parentElement.remove()" class="p-1 text-slate-400 hover:text-white transition-colors">
        <i data-lucide="x" class="w-3.5 h-3.5"></i>
      </button>
    `;

    container.appendChild(toast);
    if (window.lucide) window.lucide.createIcons({ root: toast });

    // Trigger animation
    requestAnimationFrame(() => {
      toast.classList.add('show');
    });

    // Auto dismiss
    setTimeout(() => {
      toast.classList.remove('show');
      setTimeout(() => toast.remove(), 300);
    }, duration);
  }

  // ==========================================
  // NAVIGATION & TAB SWITCHING
  // ==========================================
  switchTab(tabName) {
    this.currentTab = tabName;

    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(el => el.classList.add('hidden'));
    document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));

    // Show target tab
    const targetSection = document.getElementById(`tab-${tabName}`);
    if (targetSection) {
      targetSection.classList.remove('hidden');
    }

    const targetNav = document.getElementById(`nav-${tabName}`);
    if (targetNav) {
      targetNav.classList.add('active');
    }

    if (tabName === 'dashboard') {
      this.loadUrls(1);
    }

    if (window.lucide) window.lucide.createIcons();
  }

  // ==========================================
  // HEALTH CHECK
  // ==========================================
  async checkHealth(manual = false) {
    const healthPing = document.getElementById('health-ping');
    const healthDot = document.getElementById('health-dot');
    const healthText = document.getElementById('health-text');
    const statusVal = document.getElementById('health-status-value');
    const latencyVal = document.getElementById('health-latency');

    try {
      const res = await window.api.getHealth();
      if (res.status === 200 && res.data?.status === 'UP') {
        if (healthPing) healthPing.className = 'animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75';
        if (healthDot) healthDot.className = 'relative inline-flex rounded-full h-2 w-2 bg-emerald-500';
        if (healthText) healthText.textContent = `API Online (${res.latency}ms)`;
        if (statusVal) {
          statusVal.textContent = 'UP';
          statusVal.className = 'font-semibold text-emerald-400';
        }
        if (latencyVal) latencyVal.textContent = `${res.latency} ms`;
        if (manual) this.showToast(`Server is healthy (${res.latency}ms)`, 'success');
      }
    } catch (err) {
      if (healthPing) healthPing.className = 'hidden';
      if (healthDot) healthDot.className = 'relative inline-flex rounded-full h-2 w-2 bg-rose-500';
      if (healthText) healthText.textContent = 'API Offline';
      if (statusVal) {
        statusVal.textContent = 'DOWN';
        statusVal.className = 'font-semibold text-rose-400';
      }
      if (latencyVal) latencyVal.textContent = 'Timeout / Offline';
      if (manual) this.showToast(`Server health check failed: ${err.message}`, 'error');
    }
  }

  // ==========================================
  // AUTHENTICATION MANAGEMENT
  // ==========================================
  updateAuthUI() {
    const isAuth = window.api.isAuthenticated();
    const guestNav = document.getElementById('guest-nav');
    const userNav = document.getElementById('user-nav');
    const userEmailDisplay = document.getElementById('user-email-display');
    const dashboardGuestPrompt = document.getElementById('dashboard-guest-prompt');
    const dashboardContent = document.getElementById('dashboard-content');

    if (isAuth) {
      const user = window.api.getUser();
      if (guestNav) guestNav.classList.add('hidden');
      if (userNav) userNav.classList.remove('hidden');
      if (userEmailDisplay) userEmailDisplay.textContent = user?.email || 'User';
      if (dashboardGuestPrompt) dashboardGuestPrompt.classList.add('hidden');
      if (dashboardContent) dashboardContent.classList.remove('hidden');
    } else {
      if (guestNav) guestNav.classList.remove('hidden');
      if (userNav) userNav.classList.add('hidden');
      if (dashboardGuestPrompt) dashboardGuestPrompt.classList.remove('hidden');
      if (dashboardContent) dashboardContent.classList.add('hidden');
    }

    if (window.lucide) window.lucide.createIcons();
  }

  openAuthModal(mode = 'login') {
    const modal = document.getElementById('auth-modal');
    if (modal) modal.classList.remove('hidden');
    this.switchAuthMode(mode);
    if (window.lucide) window.lucide.createIcons();
  }

  closeAuthModal() {
    const modal = document.getElementById('auth-modal');
    if (modal) modal.classList.add('hidden');
  }

  switchAuthMode(mode) {
    const loginTab = document.getElementById('auth-tab-login');
    const registerTab = document.getElementById('auth-tab-register');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');

    if (mode === 'login') {
      loginTab.className = 'flex-1 pb-3 text-sm font-semibold border-b-2 border-indigo-500 text-white transition-colors';
      registerTab.className = 'flex-1 pb-3 text-sm font-semibold border-b-2 border-transparent text-slate-400 hover:text-slate-200 transition-colors';
      loginForm.classList.remove('hidden');
      registerForm.classList.add('hidden');
    } else {
      registerTab.className = 'flex-1 pb-3 text-sm font-semibold border-b-2 border-indigo-500 text-white transition-colors';
      loginTab.className = 'flex-1 pb-3 text-sm font-semibold border-b-2 border-transparent text-slate-400 hover:text-slate-200 transition-colors';
      registerForm.classList.remove('hidden');
      loginForm.classList.add('hidden');
    }
  }

  fillDemoCredentials() {
    const timestamp = Date.now().toString().slice(-4);
    const demoEmail = `demo${timestamp}@shortener.dev`;
    const demoPass = 'DemoPass123!';

    const loginEmail = document.getElementById('login-email');
    const loginPass = document.getElementById('login-password');
    const regEmail = document.getElementById('register-email');
    const regPass = document.getElementById('register-password');
    const regConf = document.getElementById('register-confirm-password');

    if (loginEmail) loginEmail.value = demoEmail;
    if (loginPass) loginPass.value = demoPass;
    if (regEmail) regEmail.value = demoEmail;
    if (regPass) {
      regPass.value = demoPass;
      this.validatePasswordInput();
    }
    if (regConf) regConf.value = demoPass;

    this.showToast('Test credentials filled! Password meets all validation rules.', 'info');
  }

  validatePasswordInput() {
    const pwd = document.getElementById('register-password')?.value || '';
    
    const ruleLen = pwd.length >= 8;
    const ruleUpper = /[A-Z]/.test(pwd);
    const ruleLower = /[a-z]/.test(pwd);
    const ruleNum = /[0-9]/.test(pwd);
    const ruleSpec = /[.!_]/.test(pwd);

    this.updateRuleIndicator('pwd-rule-len', ruleLen);
    this.updateRuleIndicator('pwd-rule-upper', ruleUpper);
    this.updateRuleIndicator('pwd-rule-lower', ruleLower);
    this.updateRuleIndicator('pwd-rule-num', ruleNum);
    this.updateRuleIndicator('pwd-rule-spec', ruleSpec);
  }

  updateRuleIndicator(elementId, isValid) {
    const el = document.getElementById(elementId);
    if (!el) return;
    if (isValid) {
      el.classList.add('rule-valid');
      el.querySelector('span').textContent = '✓';
    } else {
      el.classList.remove('rule-valid');
      el.querySelector('span').textContent = '•';
    }
  }

  async handleLogin(event) {
    event.preventDefault();
    const email = document.getElementById('login-email').value.trim();
    const password = document.getElementById('login-password').value;
    const submitBtn = document.getElementById('login-submit-btn');

    try {
      submitBtn.disabled = true;
      submitBtn.innerHTML = '<div class="inline-block animate-spin mr-2"><i data-lucide="loader-2" class="w-4 h-4"></i></div> Signing in...';
      if (window.lucide) window.lucide.createIcons();

      const res = await window.api.login(email, password);
      this.showToast('Signed in successfully!', 'success');
      this.closeAuthModal();
      this.updateAuthUI();

      if (this.currentTab === 'dashboard') {
        this.loadUrls(1);
      }
    } catch (err) {
      this.showToast(err.message || 'Login failed', 'error');
    } finally {
      submitBtn.disabled = false;
      submitBtn.innerHTML = '<span>Sign In</span><i data-lucide="arrow-right" class="w-4 h-4"></i>';
      if (window.lucide) window.lucide.createIcons();
    }
  }

  async handleRegister(event) {
    event.preventDefault();
    const email = document.getElementById('register-email').value.trim();
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;
    const submitBtn = document.getElementById('register-submit-btn');

    if (password !== confirmPassword) {
      this.showToast('Passwords do not match', 'error');
      return;
    }

    try {
      submitBtn.disabled = true;
      submitBtn.innerHTML = '<div class="inline-block animate-spin mr-2"><i data-lucide="loader-2" class="w-4 h-4"></i></div> Creating account...';
      if (window.lucide) window.lucide.createIcons();

      await window.api.register(email, password, confirmPassword);
      this.showToast('Account created! Logging in...', 'success');

      // Auto login
      await window.api.login(email, password);
      this.closeAuthModal();
      this.updateAuthUI();

      if (this.currentTab === 'dashboard') {
        this.loadUrls(1);
      }
    } catch (err) {
      this.showToast(err.message || 'Registration failed', 'error');
    } finally {
      submitBtn.disabled = false;
      submitBtn.innerHTML = '<span>Create Account</span><i data-lucide="check" class="w-4 h-4"></i>';
      if (window.lucide) window.lucide.createIcons();
    }
  }

  logout() {
    window.api.logout();
    this.updateAuthUI();
    this.showToast('You have been signed out.', 'info');
    if (this.currentTab === 'dashboard') {
      this.switchTab('home');
    }
  }

  // ==========================================
  // CREATE SHORT URL
  // ==========================================
  toggleAdvancedOptions() {
    const options = document.getElementById('advanced-options');
    const chevron = document.getElementById('advanced-chevron');
    if (!options) return;

    if (options.classList.contains('hidden')) {
      options.classList.remove('hidden');
      if (chevron) chevron.style.transform = 'rotate(180deg)';
    } else {
      options.classList.add('hidden');
      if (chevron) chevron.style.transform = 'rotate(0deg)';
    }
  }

  async handleCreateUrl(event) {
    event.preventDefault();

    if (!window.api.isAuthenticated()) {
      this.showToast('Please sign in or create an account to generate short links.', 'warning');
      this.openAuthModal('login');
      return;
    }

    const originalUrlInput = document.getElementById('original-url-input');
    const customAliasInput = document.getElementById('custom-alias-input');
    const expiresAtInput = document.getElementById('expires-at-input');
    const shortenBtn = document.getElementById('shorten-btn');

    let originalUrl = originalUrlInput.value.trim();
    if (!originalUrl.startsWith('http://') && !originalUrl.startsWith('https://')) {
      originalUrl = 'https://' + originalUrl;
    }

    const customAlias = customAliasInput?.value?.trim() || '';
    const expiresAt = expiresAtInput?.value || null;

    try {
      shortenBtn.disabled = true;
      shortenBtn.innerHTML = '<div class="inline-block animate-spin mr-2"><i data-lucide="loader-2" class="w-4 h-4"></i></div> Shortening...';
      if (window.lucide) window.lucide.createIcons();

      const res = await window.api.createUrl({
        originalUrl,
        customAlias,
        expiresAt
      });

      const urlObj = res.data;
      this.latestCreatedUrl = urlObj;
      this.renderCreatedResult(urlObj);
      this.showToast('Short link created successfully!', 'success');

      // Clear inputs
      originalUrlInput.value = '';
      if (customAliasInput) customAliasInput.value = '';
      if (expiresAtInput) expiresAtInput.value = '';

      // If dashboard was loaded, refresh it
      if (this.currentTab === 'dashboard') {
        this.loadUrls(1);
      }
    } catch (err) {
      this.showToast(err.message || 'Failed to create short URL', 'error');
    } finally {
      shortenBtn.disabled = false;
      shortenBtn.innerHTML = '<i data-lucide="scissors" class="w-4 h-4"></i><span>Shorten URL</span>';
      if (window.lucide) window.lucide.createIcons();
    }
  }

  renderCreatedResult(urlObj) {
    const container = document.getElementById('result-container');
    const shortLinkEl = document.getElementById('result-short-link');
    const origLinkEl = document.getElementById('result-original-link');
    const timeEl = document.getElementById('result-created-time');

    if (!container || !shortLinkEl || !origLinkEl) return;

    const shortCode = urlObj.ShortCode || urlObj.short_code;
    const fullShortUrl = `${window.api.baseUrl}/${shortCode}`;
    const origUrl = urlObj.OriginalURL || urlObj.original_url;

    shortLinkEl.textContent = fullShortUrl;
    shortLinkEl.href = fullShortUrl;
    origLinkEl.textContent = origUrl;
    if (timeEl) timeEl.textContent = new Date().toLocaleTimeString();

    container.classList.remove('hidden');
    if (window.lucide) window.lucide.createIcons();
  }

  copyResultLink() {
    if (!this.latestCreatedUrl) return;
    const shortCode = this.latestCreatedUrl.ShortCode || this.latestCreatedUrl.short_code;
    const fullUrl = `${window.api.baseUrl}/${shortCode}`;
    this.copyToClipboard(fullUrl);

    const btnText = document.getElementById('result-copy-text');
    if (btnText) {
      btnText.textContent = 'Copied!';
      setTimeout(() => {
        btnText.textContent = 'Copy Link';
      }, 2000);
    }
  }

  copyToClipboard(text) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(() => {
        this.showToast('Copied to clipboard!', 'success');
      }).catch(() => {
        this.fallbackCopy(text);
      });
    } else {
      this.fallbackCopy(text);
    }
  }

  fallbackCopy(text) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand('copy');
      this.showToast('Copied to clipboard!', 'success');
    } catch {
      this.showToast('Failed to copy', 'error');
    }
    document.body.removeChild(textarea);
  }

  // ==========================================
  // DASHBOARD LINKS LIST & PAGINATION
  // ==========================================
  async loadUrls(page = 1) {
    if (!window.api.isAuthenticated()) {
      this.updateAuthUI();
      return;
    }

    this.currentPage = page;
    const tbody = document.getElementById('links-tbody');
    if (tbody) {
      tbody.innerHTML = `
        <tr>
          <td colspan="5" class="text-center py-12 text-slate-500">
            <div class="inline-block animate-spin mr-2"><i data-lucide="loader-2" class="w-5 h-5"></i></div>
            Fetching links...
          </td>
        </tr>
      `;
      if (window.lucide) window.lucide.createIcons();
    }

    try {
      const res = await window.api.getUrls(this.currentPage, this.pageSize);
      const data = res.data?.data || [];
      const pagination = res.data?.pagination || {};

      this.urls = Array.isArray(data) ? data : [];
      this.totalLinks = pagination.total || this.urls.length;
      this.totalPages = pagination.total_pages || Math.ceil(this.totalLinks / this.pageSize) || 1;

      this.applyFilterAndSearch();
      this.updateMetricsRibbon();
      this.renderPagination();
    } catch (err) {
      if (err.status === 401) {
        window.api.logout();
        this.updateAuthUI();
        this.showToast('Session expired. Please sign in again.', 'warning');
      } else {
        if (tbody) {
          tbody.innerHTML = `
            <tr>
              <td colspan="5" class="text-center py-8 text-rose-400">
                Failed to load links: ${err.message}
              </td>
            </tr>
          `;
        }
      }
    }
  }

  updateMetricsRibbon() {
    const totalEl = document.getElementById('stat-total-links');
    const activeEl = document.getElementById('stat-active-links');
    const inactiveEl = document.getElementById('stat-inactive-links');
    const pageEl = document.getElementById('stat-current-page');

    const activeCount = this.urls.filter(u => {
      const isAct = (u.IsActive !== undefined) ? u.IsActive : u.is_active;
      const isDel = (u.IsDeleted !== undefined) ? u.IsDeleted : u.is_deleted;
      return isAct && !isDel;
    }).length;

    const inactiveCount = this.urls.length - activeCount;

    if (totalEl) totalEl.textContent = this.totalLinks.toLocaleString();
    if (activeEl) activeEl.textContent = activeCount.toLocaleString();
    if (inactiveEl) inactiveEl.textContent = inactiveCount.toLocaleString();
    if (pageEl) pageEl.textContent = `Page ${this.currentPage} / ${this.totalPages || 1}`;
  }

  handleSearch(event) {
    this.searchQuery = event.target.value.toLowerCase().trim();
    this.applyFilterAndSearch();
  }

  handleFilterChange() {
    const sel = document.getElementById('filter-status');
    this.filterStatus = sel ? sel.value : 'all';
    this.applyFilterAndSearch();
  }

  handlePageSizeChange() {
    const sel = document.getElementById('page-size');
    this.pageSize = sel ? parseInt(sel.value, 10) : 10;
    this.loadUrls(1);
  }

  applyFilterAndSearch() {
    let list = [...this.urls];

    // Status filter
    if (this.filterStatus === 'active') {
      list = list.filter(u => (u.IsActive !== undefined ? u.IsActive : u.is_active));
    } else if (this.filterStatus === 'inactive') {
      list = list.filter(u => !(u.IsActive !== undefined ? u.IsActive : u.is_active));
    }

    // Search filter
    if (this.searchQuery) {
      list = list.filter(u => {
        const code = (u.ShortCode || u.short_code || '').toLowerCase();
        const orig = (u.OriginalURL || u.original_url || '').toLowerCase();
        return code.includes(this.searchQuery) || orig.includes(this.searchQuery);
      });
    }

    this.filteredUrls = list;
    this.renderTable();
  }

  renderTable() {
    const tbody = document.getElementById('links-tbody');
    if (!tbody) return;

    if (this.filteredUrls.length === 0) {
      tbody.innerHTML = `
        <tr>
          <td colspan="5" class="text-center py-12 text-slate-500">
            <i data-lucide="inbox" class="w-8 h-8 mx-auto mb-2 opacity-50"></i>
            <p class="text-sm font-medium">No links found matching your filters.</p>
          </td>
        </tr>
      `;
      if (window.lucide) window.lucide.createIcons();
      return;
    }

    tbody.innerHTML = this.filteredUrls.map(item => {
      const id = item.ID || item.id;
      const code = item.ShortCode || item.short_code;
      const original = item.OriginalURL || item.original_url;
      const isActive = item.IsActive !== undefined ? item.IsActive : item.is_active;
      const expiresAt = item.ExpiresAt || item.expires_at;
      const fullShortUrl = `${window.api.baseUrl}/${code}`;

      let isExpired = false;
      let expiryLabel = 'Never';
      if (expiresAt) {
        const expDate = new Date(expiresAt);
        if (!isNaN(expDate.getTime())) {
          isExpired = expDate < new Date();
          expiryLabel = expDate.toLocaleDateString() + ' ' + expDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
      }

      let statusBadge = '<span class="px-2 py-0.5 rounded text-[11px] font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">Active</span>';
      if (isExpired) {
        statusBadge = '<span class="px-2 py-0.5 rounded text-[11px] font-semibold bg-rose-500/10 text-rose-400 border border-rose-500/20">Expired</span>';
      } else if (!isActive) {
        statusBadge = '<span class="px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-500/10 text-amber-400 border border-amber-500/20">Inactive</span>';
      }

      return `
        <tr class="link-row border-b border-slate-800/60 hover:bg-slate-800/30 transition-colors">
          <td class="py-3.5 px-4 sm:px-6">
            <div class="flex items-center gap-2">
              <span class="font-mono font-semibold text-indigo-400 hover:text-indigo-300 transition-colors">
                /${code}
              </span>
              <button onclick="app.copyToClipboard('${fullShortUrl}')" title="Copy Short Link" class="p-1 rounded text-slate-400 hover:text-white hover:bg-slate-800 transition-colors">
                <i data-lucide="copy" class="w-3.5 h-3.5"></i>
              </button>
              <a href="${fullShortUrl}" target="_blank" title="Test Redirect" class="p-1 rounded text-slate-400 hover:text-white hover:bg-slate-800 transition-colors">
                <i data-lucide="external-link" class="w-3.5 h-3.5"></i>
              </a>
            </div>
          </td>
          
          <td class="py-3.5 px-4 hidden md:table-cell max-w-xs">
            <div class="truncate text-xs text-slate-400" title="${original}">
              ${original}
            </div>
          </td>

          <td class="py-3.5 px-4">
            ${statusBadge}
          </td>

          <td class="py-3.5 px-4 hidden lg:table-cell text-xs text-slate-400 font-mono">
            ${expiryLabel}
          </td>

          <td class="py-3.5 px-4 text-right">
            <div class="flex items-center justify-end gap-1">
              
              <!-- Analytics -->
              <button
                onclick="app.openAnalyticsModal(${id}, '${code}')"
                title="View Analytics"
                class="p-1.5 rounded-lg text-slate-400 hover:text-indigo-400 hover:bg-indigo-500/10 transition-colors"
              >
                <i data-lucide="bar-chart-3" class="w-4 h-4"></i>
              </button>

              <!-- QR Code -->
              <button
                onclick='app.showQrModal(${JSON.stringify(item)})'
                title="Generate QR Code"
                class="p-1.5 rounded-lg text-slate-400 hover:text-purple-400 hover:bg-purple-500/10 transition-colors"
              >
                <i data-lucide="qr-code" class="w-4 h-4"></i>
              </button>

              <!-- Toggle Active -->
              <button
                onclick="app.toggleUrlActive(${id}, ${isActive})"
                title="${isActive ? 'Deactivate link' : 'Activate link'}"
                class="p-1.5 rounded-lg text-slate-400 hover:text-amber-400 hover:bg-amber-500/10 transition-colors"
              >
                <i data-lucide="${isActive ? 'pause' : 'play'}" class="w-4 h-4"></i>
              </button>

              <!-- Edit -->
              <button
                onclick='app.openEditModal(${JSON.stringify(item)})'
                title="Edit URL"
                class="p-1.5 rounded-lg text-slate-400 hover:text-sky-400 hover:bg-sky-500/10 transition-colors"
              >
                <i data-lucide="edit-3" class="w-4 h-4"></i>
              </button>

              <!-- Delete -->
              <button
                onclick="app.confirmDeleteUrl(${id}, '${code}')"
                title="Delete URL"
                class="p-1.5 rounded-lg text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 transition-colors"
              >
                <i data-lucide="trash-2" class="w-4 h-4"></i>
              </button>

            </div>
          </td>
        </tr>
      `;
    }).join('');

    if (window.lucide) window.lucide.createIcons();
  }

  renderPagination() {
    const infoEl = document.getElementById('pagination-info');
    const numbersEl = document.getElementById('page-numbers');
    const prevBtn = document.getElementById('prev-page-btn');
    const nextBtn = document.getElementById('next-page-btn');

    const start = (this.currentPage - 1) * this.pageSize + 1;
    const end = Math.min(this.currentPage * this.pageSize, this.totalLinks);

    if (infoEl) {
      infoEl.innerHTML = `Showing <span class="font-semibold text-white">${this.totalLinks > 0 ? start : 0}</span> to <span class="font-semibold text-white">${end}</span> of <span class="font-semibold text-white">${this.totalLinks}</span> entries`;
    }

    if (prevBtn) prevBtn.disabled = this.currentPage <= 1;
    if (nextBtn) nextBtn.disabled = this.currentPage >= this.totalPages;

    if (numbersEl) {
      numbersEl.innerHTML = '';
      for (let p = 1; p <= this.totalPages; p++) {
        if (p === 1 || p === this.totalPages || (p >= this.currentPage - 1 && p <= this.currentPage + 1)) {
          const btn = document.createElement('button');
          btn.className = `w-8 h-8 rounded-lg text-xs font-semibold transition-colors ${p === this.currentPage ? 'bg-indigo-600 text-white shadow-sm' : 'bg-slate-800 hover:bg-slate-700 text-slate-300'}`;
          btn.textContent = p;
          btn.onclick = () => this.changePage(p);
          numbersEl.appendChild(btn);
        }
      }
    }
  }

  changePage(page) {
    if (page < 1 || page > this.totalPages || page === this.currentPage) return;
    this.loadUrls(page);
  }

  // ==========================================
  // TOGGLE ACTIVE, EDIT & DELETE ACTIONS
  // ==========================================
  async toggleUrlActive(id, currentActive) {
    try {
      if (currentActive) {
        await window.api.deactivateUrl(id);
        this.showToast('Link deactivated', 'info');
      } else {
        await window.api.activateUrl(id);
        this.showToast('Link activated', 'success');
      }
      this.loadUrls(this.currentPage);
    } catch (err) {
      this.showToast(err.message || 'Failed to toggle status', 'error');
    }
  }

  async confirmDeleteUrl(id, code) {
    if (!confirm(`Are you sure you want to delete short URL /${code}?`)) return;

    try {
      await window.api.deleteUrl(id);
      this.showToast(`Deleted /${code}`, 'success');
      this.loadUrls(this.currentPage);
    } catch (err) {
      this.showToast(err.message || 'Failed to delete URL', 'error');
    }
  }

  openEditModal(item) {
    const modal = document.getElementById('edit-modal');
    const idInput = document.getElementById('edit-id');
    const codeEl = document.getElementById('edit-modal-code');
    const origInput = document.getElementById('edit-original-url');
    const expInput = document.getElementById('edit-expires-at');
    const activeCheck = document.getElementById('edit-is-active');

    const id = item.ID || item.id;
    const code = item.ShortCode || item.short_code;
    const orig = item.OriginalURL || item.original_url;
    const isActive = item.IsActive !== undefined ? item.IsActive : item.is_active;
    const expiresAt = item.ExpiresAt || item.expires_at;

    if (idInput) idInput.value = id;
    if (codeEl) codeEl.textContent = `/${code}`;
    if (origInput) origInput.value = orig;
    if (activeCheck) activeCheck.checked = !!isActive;

    if (expInput) {
      if (expiresAt) {
        const d = new Date(expiresAt);
        if (!isNaN(d.getTime())) {
          expInput.value = d.toISOString().slice(0, 16);
        } else {
          expInput.value = '';
        }
      } else {
        expInput.value = '';
      }
    }

    if (modal) modal.classList.remove('hidden');
    if (window.lucide) window.lucide.createIcons();
  }

  closeEditModal() {
    const modal = document.getElementById('edit-modal');
    if (modal) modal.classList.add('hidden');
  }

  async handleUpdateUrl(event) {
    event.preventDefault();
    const id = document.getElementById('edit-id').value;
    const originalUrl = document.getElementById('edit-original-url').value.trim();
    const expiresAt = document.getElementById('edit-expires-at').value || null;
    const isActive = document.getElementById('edit-is-active').checked;

    try {
      await window.api.updateUrl(id, {
        originalUrl,
        expiresAt,
        isActive
      });

      this.showToast('Link updated successfully', 'success');
      this.closeEditModal();
      this.loadUrls(this.currentPage);
    } catch (err) {
      this.showToast(err.message || 'Failed to update link', 'error');
    }
  }

  // ==========================================
  // ANALYTICS MODAL
  // ==========================================
  async openAnalyticsModal(urlId, shortCode) {
    this.activeAnalyticsUrlId = urlId;
    this.activeAnalyticsShortCode = shortCode;

    const modal = document.getElementById('analytics-modal');
    if (modal) modal.classList.remove('hidden');

    await this.refreshCurrentAnalytics();
  }

  closeAnalyticsModal() {
    const modal = document.getElementById('analytics-modal');
    if (modal) modal.classList.add('hidden');
  }

  async refreshCurrentAnalytics() {
    if (!this.activeAnalyticsUrlId) return;

    try {
      const res = await window.api.getAnalytics(this.activeAnalyticsUrlId);
      const data = res.data || {};
      window.analyticsRenderer.populate(this.activeAnalyticsUrlId, this.activeAnalyticsShortCode, data);
      if (window.lucide) window.lucide.createIcons();
    } catch (err) {
      this.showToast(err.message || 'Failed to load analytics', 'error');
    }
  }

  // ==========================================
  // QR CODE MODAL
  // ==========================================
  showQrModal(item) {
    if (!item) return;
    const code = item.ShortCode || item.short_code;
    const fullUrl = `${window.api.baseUrl}/${code}`;

    const modal = document.getElementById('qr-modal');
    const urlEl = document.getElementById('qr-short-url');
    const renderEl = document.getElementById('qrcode-render');

    if (urlEl) urlEl.textContent = fullUrl;

    if (renderEl) {
      renderEl.innerHTML = '';
      if (window.QRCode) {
        this.qrCodeInstance = new window.QRCode(renderEl, {
          text: fullUrl,
          width: 180,
          height: 180,
          colorDark: '#0f172a',
          colorLight: '#ffffff',
          correctLevel: window.QRCode.CorrectLevel.H
        });
      }
    }

    if (modal) modal.classList.remove('hidden');
    if (window.lucide) window.lucide.createIcons();
  }

  closeQrModal() {
    const modal = document.getElementById('qr-modal');
    if (modal) modal.classList.add('hidden');
  }

  downloadQrCode() {
    const qrImg = document.querySelector('#qrcode-render img');
    const qrCanvas = document.querySelector('#qrcode-render canvas');
    let src = '';

    if (qrImg && qrImg.src) {
      src = qrImg.src;
    } else if (qrCanvas) {
      src = qrCanvas.toDataURL('image/png');
    }

    if (!src) {
      this.showToast('QR code not available for download', 'error');
      return;
    }

    const a = document.createElement('a');
    a.href = src;
    a.download = `qrcode-${this.activeAnalyticsShortCode || 'link'}.png`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    this.showToast('Downloaded QR Code', 'success');
  }

  // ==========================================
  // INTERACTIVE API TESTER
  // ==========================================
  prefillApiTester(method, path, body, useAuth) {
    this.switchTab('api');

    const methodEl = document.getElementById('api-method');
    const pathEl = document.getElementById('api-path');
    const bodyEl = document.getElementById('api-body');
    const authEl = document.getElementById('api-use-auth');

    if (methodEl) methodEl.value = method;
    if (pathEl) pathEl.value = path;
    if (bodyEl) bodyEl.value = body || '';
    if (authEl) authEl.checked = useAuth;

    const consoleEl = document.getElementById('interactive-console');
    if (consoleEl) {
      consoleEl.scrollIntoView({ behavior: 'smooth' });
    }
  }

  async executeApiTest(event) {
    event.preventDefault();
    const method = document.getElementById('api-method').value;
    const path = document.getElementById('api-path').value.trim();
    const rawBody = document.getElementById('api-body').value.trim();
    const useAuth = document.getElementById('api-use-auth').checked;
    const sendBtn = document.getElementById('api-send-btn');
    const respContainer = document.getElementById('api-response-container');
    const statusEl = document.getElementById('api-response-status');
    const timeEl = document.getElementById('api-response-time');
    const bodyEl = document.getElementById('api-response-body');

    let body = null;
    if (['POST', 'PUT', 'PATCH'].includes(method) && rawBody) {
      try {
        body = JSON.parse(rawBody);
      } catch {
        this.showToast('Invalid JSON in Request Body', 'error');
        return;
      }
    }

    try {
      sendBtn.disabled = true;
      sendBtn.innerHTML = '<div class="inline-block animate-spin mr-2"><i data-lucide="loader-2" class="w-4 h-4"></i></div> Sending...';
      if (window.lucide) window.lucide.createIcons();

      const res = await window.api.request(path, {
        method,
        body,
        useAuth
      });

      respContainer.classList.remove('hidden');
      statusEl.textContent = `${res.status} ${res.statusText || 'OK'}`;
      statusEl.className = 'px-2 py-0.5 rounded font-mono font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20';
      timeEl.textContent = `${res.latency} ms`;
      bodyEl.textContent = JSON.stringify(res.data, null, 2);
    } catch (err) {
      respContainer.classList.remove('hidden');
      statusEl.textContent = `${err.status || 'ERR'} ${err.status ? 'FAILED' : 'NETWORK ERROR'}`;
      statusEl.className = 'px-2 py-0.5 rounded font-mono font-bold bg-rose-500/10 text-rose-400 border border-rose-500/20';
      timeEl.textContent = `${err.latency || 0} ms`;
      bodyEl.textContent = JSON.stringify(err.data || { error: err.message }, null, 2);
    } finally {
      sendBtn.disabled = false;
      sendBtn.innerHTML = '<i data-lucide="send" class="w-4 h-4"></i> Send Request';
      if (window.lucide) window.lucide.createIcons();
    }
  }
}

// Instantiate global app
window.app = new App();
