/**
 * GoShort Analytics Visualization Engine
 * Renders interactive Chart.js charts and metrics breakdowns.
 */
class AnalyticsRenderer {
  constructor() {
    this.chartInstance = null;
  }

  // Render or Update the Device Chart
  renderDeviceChart(devicesData = {}) {
    const canvas = document.getElementById('deviceChart');
    if (!canvas) return;
    const ctx = canvas.getContext('2d');

    if (this.chartInstance) {
      this.chartInstance.destroy();
      this.chartInstance = null;
    }

    const labels = Object.keys(devicesData);
    const dataValues = Object.values(devicesData);

    const legendContainer = document.getElementById('device-legend');
    if (legendContainer) {
      legendContainer.innerHTML = '';
    }

    if (labels.length === 0 || dataValues.reduce((a, b) => a + b, 0) === 0) {
      // Empty state
      if (legendContainer) {
        legendContainer.innerHTML = '<div class="text-center py-6 text-slate-500 text-xs italic">No device click data recorded yet.</div>';
      }
      return;
    }

    const colorPalette = [
      { bg: '#6366f1', border: '#4f46e5', label: 'Desktop' },
      { bg: '#a855f7', border: '#9333ea', label: 'Mobile' },
      { bg: '#ec4899', border: '#db2777', label: 'Tablet' },
      { bg: '#3b82f6', border: '#2563eb', label: 'Bot / Other' },
      { bg: '#10b981', border: '#059669', label: 'Other' },
    ];

    const bgColors = labels.map((_, i) => colorPalette[i % colorPalette.length].bg);
    const borderColors = labels.map((_, i) => colorPalette[i % colorPalette.length].border);

    this.chartInstance = new Chart(ctx, {
      type: 'doughnut',
      data: {
        labels: labels.map(l => l.charAt(0).toUpperCase() + l.slice(1)),
        datasets: [{
          data: dataValues,
          backgroundColor: bgColors,
          borderColor: borderColors,
          borderWidth: 2,
          hoverOffset: 4
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        cutout: '70%',
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            backgroundColor: '#0f172a',
            titleColor: '#ffffff',
            bodyColor: '#cbd5e1',
            borderColor: '#334155',
            borderWidth: 1,
            padding: 10,
            boxPadding: 4,
            callbacks: {
              label: function(context) {
                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                const val = context.parsed;
                const pct = total > 0 ? Math.round((val / total) * 100) : 0;
                return ` ${context.label}: ${val} (${pct}%)`;
              }
            }
          }
        }
      }
    });

    // Populate Custom Legend
    if (legendContainer) {
      const total = dataValues.reduce((a, b) => a + b, 0);
      labels.forEach((label, idx) => {
        const val = dataValues[idx];
        const pct = total > 0 ? Math.round((val / total) * 100) : 0;
        const color = bgColors[idx];

        const item = document.createElement('div');
        item.className = 'flex items-center justify-between text-xs text-slate-300';
        item.innerHTML = `
          <div class="flex items-center gap-2">
            <span class="w-2.5 h-2.5 rounded-full" style="background-color: ${color}"></span>
            <span class="capitalize">${label}</span>
          </div>
          <div class="flex items-center gap-2 font-mono">
            <span class="text-white font-semibold">${val}</span>
            <span class="text-slate-500">(${pct}%)</span>
          </div>
        `;
        legendContainer.appendChild(item);
      });
    }
  }

  // Render Countries List and Percentage Progress Bars
  renderCountriesList(countriesData = {}) {
    const container = document.getElementById('countries-list');
    if (!container) return;

    container.innerHTML = '';

    const entries = Object.entries(countriesData).sort((a, b) => b[1] - a[1]);
    const total = entries.reduce((acc, [, count]) => acc + count, 0);

    if (entries.length === 0 || total === 0) {
      container.innerHTML = '<div class="text-center py-6 text-slate-500 text-xs italic">No geographic click data recorded yet.</div>';
      return;
    }

    // Helper for flag emojis from ISO 2-letter country codes
    const getFlagEmoji = (countryCode) => {
      if (!countryCode || countryCode.length !== 2 || countryCode === 'Unknown') return '🌐';
      const codePoints = countryCode
        .toUpperCase()
        .split('')
        .map(char => 127397 + char.charCodeAt(0));
      return String.fromCodePoint(...codePoints);
    };

    entries.forEach(([country, count]) => {
      const pct = total > 0 ? Math.round((count / total) * 100) : 0;
      const flag = getFlagEmoji(country);

      const row = document.createElement('div');
      row.className = 'p-2 rounded-lg bg-slate-900/60 border border-slate-800 text-xs space-y-1';
      row.innerHTML = `
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 font-medium text-slate-200">
            <span class="text-sm">${flag}</span>
            <span class="font-mono font-semibold">${country}</span>
          </div>
          <div class="flex items-center gap-2 font-mono">
            <span class="text-white font-bold">${count}</span>
            <span class="text-slate-500 text-[11px]">(${pct}%)</span>
          </div>
        </div>
        <div class="w-full bg-slate-800 rounded-full h-1.5 overflow-hidden">
          <div class="bg-gradient-to-r from-purple-500 to-indigo-500 h-1.5 rounded-full transition-all duration-500" style="width: ${pct}%"></div>
        </div>
      `;
      container.appendChild(row);
    });
  }

  // Populate All Modal Analytics Details
  populate(urlId, shortCode, data) {
    document.getElementById('analytics-url-id').textContent = urlId;
    document.getElementById('analytics-short-code').textContent = `/${shortCode}`;
    
    const totalClicks = data.total_clicks || 0;
    const uniqueVisitors = data.unique_visitors || 0;
    
    document.getElementById('analytics-total-clicks').textContent = totalClicks.toLocaleString();
    document.getElementById('analytics-unique-visitors').textContent = uniqueVisitors.toLocaleString();

    let repeatRate = '0%';
    if (totalClicks > 0 && uniqueVisitors > 0) {
      const rate = Math.round(((totalClicks - uniqueVisitors) / totalClicks) * 100);
      repeatRate = `${Math.max(0, rate)}%`;
    }
    document.getElementById('analytics-repeat-rate').textContent = repeatRate;

    this.renderDeviceChart(data.devices || {});
    this.renderCountriesList(data.countries || {});
  }
}

window.analyticsRenderer = new AnalyticsRenderer();
