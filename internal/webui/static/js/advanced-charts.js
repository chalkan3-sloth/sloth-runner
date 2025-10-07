// Advanced Charts Library for Sloth Runner
// Provides Gauge, Heatmap, Sankey, and other advanced visualizations using Chart.js plugins

class AdvancedCharts {
    constructor() {
        this.charts = new Map();
        this.defaultColors = {
            primary: '#4F46E5',
            success: '#10B981',
            warning: '#F59E0B',
            danger: '#EF4444',
            info: '#3B82F6'
        };
    }

    /**
     * Create a Gauge Chart
     * Perfect for displaying percentages, progress, or single metrics
     */
    createGauge(canvasId, options = {}) {
        const {
            value = 0,
            max = 100,
            min = 0,
            label = '',
            unit = '%',
            color = null,
            thresholds = null,
            animated = true
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        // Calculate percentage
        const percentage = ((value - min) / (max - min)) * 100;

        // Determine color based on thresholds or use default
        let gaugeColor = color || this.defaultColors.primary;
        if (thresholds) {
            if (percentage >= thresholds.excellent) {
                gaugeColor = this.defaultColors.success;
            } else if (percentage >= thresholds.good) {
                gaugeColor = this.defaultColors.info;
            } else if (percentage >= thresholds.warning) {
                gaugeColor = this.defaultColors.warning;
            } else {
                gaugeColor = this.defaultColors.danger;
            }
        }

        const config = {
            type: 'doughnut',
            data: {
                datasets: [{
                    data: [value, max - value],
                    backgroundColor: [gaugeColor, 'rgba(200, 200, 200, 0.2)'],
                    borderWidth: 0,
                    circumference: 180,
                    rotation: 270
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                cutout: '75%',
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        enabled: false
                    }
                },
                animation: animated ? {
                    animateRotate: true,
                    animateScale: true,
                    duration: 1000,
                    easing: 'easeOutCubic'
                } : false
            },
            plugins: [{
                id: 'gaugeText',
                afterDraw: (chart) => {
                    const ctx = chart.ctx;
                    const centerX = chart.chartArea.left + (chart.chartArea.right - chart.chartArea.left) / 2;
                    const centerY = chart.chartArea.top + (chart.chartArea.bottom - chart.chartArea.top) / 2 + 20;

                    ctx.save();
                    ctx.textAlign = 'center';
                    ctx.textBaseline = 'middle';

                    // Draw value
                    ctx.font = 'bold 32px sans-serif';
                    ctx.fillStyle = gaugeColor;
                    ctx.fillText(`${value}${unit}`, centerX, centerY);

                    // Draw label
                    if (label) {
                        ctx.font = '14px sans-serif';
                        ctx.fillStyle = '#666';
                        ctx.fillText(label, centerX, centerY + 30);
                    }

                    ctx.restore();
                }
            }]
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Create a Radial Progress Chart
     * Similar to gauge but full circle
     */
    createRadialProgress(canvasId, options = {}) {
        const {
            value = 0,
            max = 100,
            label = '',
            colors = [this.defaultColors.primary, this.defaultColors.info]
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        const percentage = (value / max) * 100;

        const config = {
            type: 'doughnut',
            data: {
                datasets: [{
                    data: [value, max - value],
                    backgroundColor: colors,
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                cutout: '80%',
                plugins: {
                    legend: { display: false },
                    tooltip: { enabled: false }
                },
                animation: {
                    duration: 1000,
                    easing: 'easeInOutQuart'
                }
            },
            plugins: [{
                id: 'radialText',
                afterDraw: (chart) => {
                    const ctx = chart.ctx;
                    const centerX = chart.chartArea.left + (chart.chartArea.right - chart.chartArea.left) / 2;
                    const centerY = chart.chartArea.top + (chart.chartArea.bottom - chart.chartArea.top) / 2;

                    ctx.save();
                    ctx.textAlign = 'center';
                    ctx.textBaseline = 'middle';

                    ctx.font = 'bold 36px sans-serif';
                    ctx.fillStyle = colors[0];
                    ctx.fillText(`${percentage.toFixed(0)}%`, centerX, centerY - 10);

                    if (label) {
                        ctx.font = '14px sans-serif';
                        ctx.fillStyle = '#666';
                        ctx.fillText(label, centerX, centerY + 25);
                    }

                    ctx.restore();
                }
            }]
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Create a Matrix/Heatmap Chart
     * Great for showing correlation or activity patterns
     */
    createHeatmap(canvasId, data, options = {}) {
        const {
            xLabels = [],
            yLabels = [],
            colorScheme = 'blue', // blue, green, red, purple
            showValues = true
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        const colorSchemes = {
            blue: ['#EFF6FF', '#DBEAFE', '#BFDBFE', '#93C5FD', '#60A5FA', '#3B82F6', '#2563EB', '#1D4ED8'],
            green: ['#ECFDF5', '#D1FAE5', '#A7F3D0', '#6EE7B7', '#34D399', '#10B981', '#059669', '#047857'],
            red: ['#FEF2F2', '#FEE2E2', '#FECACA', '#FCA5A5', '#F87171', '#EF4444', '#DC2626', '#B91C1C'],
            purple: ['#F5F3FF', '#EDE9FE', '#DDD6FE', '#C4B5FD', '#A78BFA', '#8B5CF6', '#7C3AED', '#6D28D9']
        };

        const colors = colorSchemes[colorScheme] || colorSchemes.blue;

        // Flatten data for Chart.js
        const chartData = [];
        data.forEach((row, y) => {
            row.forEach((value, x) => {
                chartData.push({
                    x: xLabels[x] || x,
                    y: yLabels[y] || y,
                    v: value
                });
            });
        });

        // Find max value for color scaling
        const maxValue = Math.max(...data.flat());

        const config = {
            type: 'scatter',
            data: {
                datasets: [{
                    data: chartData,
                    backgroundColor: function(context) {
                        const value = context.raw.v;
                        const index = Math.floor((value / maxValue) * (colors.length - 1));
                        return colors[index];
                    },
                    borderColor: 'rgba(255, 255, 255, 0.5)',
                    borderWidth: 2,
                    pointStyle: 'rect',
                    pointRadius: 20
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        type: 'category',
                        labels: xLabels,
                        grid: { display: false }
                    },
                    y: {
                        type: 'category',
                        labels: yLabels,
                        grid: { display: false }
                    }
                },
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `Value: ${context.raw.v}`;
                            }
                        }
                    }
                }
            }
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Create a Sparkline Chart
     * Small, simple line chart without axes
     */
    createSparkline(canvasId, data, options = {}) {
        const {
            color = this.defaultColors.primary,
            fillColor = null,
            showPoints = false
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        const config = {
            type: 'line',
            data: {
                labels: data.map((_, i) => i),
                datasets: [{
                    data: data,
                    borderColor: color,
                    backgroundColor: fillColor || `${color}22`,
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4,
                    pointRadius: showPoints ? 3 : 0,
                    pointHoverRadius: showPoints ? 5 : 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: { display: false },
                    y: { display: false }
                },
                plugins: {
                    legend: { display: false },
                    tooltip: { enabled: showPoints }
                },
                elements: {
                    point: { hoverRadius: 4 }
                }
            }
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Create a Polar Area Chart
     * Good for comparing multiple categories
     */
    createPolarArea(canvasId, data, options = {}) {
        const {
            labels = [],
            colors = [
                this.defaultColors.primary,
                this.defaultColors.success,
                this.defaultColors.warning,
                this.defaultColors.danger,
                this.defaultColors.info
            ]
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        const config = {
            type: 'polarArea',
            data: {
                labels: labels,
                datasets: [{
                    data: data,
                    backgroundColor: colors.map(c => `${c}CC`),
                    borderColor: colors,
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            padding: 15,
                            font: { size: 12 }
                        }
                    }
                },
                scales: {
                    r: {
                        ticks: { display: false },
                        grid: { color: 'rgba(0, 0, 0, 0.05)' }
                    }
                }
            }
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Create a Radar Chart
     * Perfect for showing multi-dimensional data
     */
    createRadar(canvasId, data, options = {}) {
        const {
            labels = [],
            datasets = [],
            colors = [this.defaultColors.primary, this.defaultColors.success]
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        const chartDatasets = datasets.map((dataset, index) => ({
            label: dataset.label || `Dataset ${index + 1}`,
            data: dataset.data,
            backgroundColor: `${colors[index % colors.length]}33`,
            borderColor: colors[index % colors.length],
            borderWidth: 2,
            pointBackgroundColor: colors[index % colors.length],
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: colors[index % colors.length]
        }));

        const config = {
            type: 'radar',
            data: {
                labels: labels,
                datasets: chartDatasets
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    r: {
                        beginAtZero: true,
                        ticks: { display: false },
                        grid: { color: 'rgba(0, 0, 0, 0.05)' }
                    }
                },
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            padding: 15,
                            font: { size: 12 }
                        }
                    }
                }
            }
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Create a Stacked Area Chart
     * Great for showing composition over time
     */
    createStackedArea(canvasId, data, options = {}) {
        const {
            labels = [],
            datasets = [],
            colors = [
                this.defaultColors.primary,
                this.defaultColors.success,
                this.defaultColors.warning
            ]
        } = options;

        const canvas = document.getElementById(canvasId);
        if (!canvas) return null;

        const chartDatasets = datasets.map((dataset, index) => ({
            label: dataset.label,
            data: dataset.data,
            backgroundColor: `${colors[index % colors.length]}99`,
            borderColor: colors[index % colors.length],
            borderWidth: 2,
            fill: true,
            tension: 0.4
        }));

        const config = {
            type: 'line',
            data: {
                labels: labels,
                datasets: chartDatasets
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: {
                    mode: 'index',
                    intersect: false
                },
                scales: {
                    y: {
                        stacked: true,
                        beginAtZero: true,
                        grid: { color: 'rgba(0, 0, 0, 0.05)' }
                    },
                    x: {
                        grid: { display: false }
                    }
                },
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            padding: 15,
                            font: { size: 12 }
                        }
                    },
                    tooltip: {
                        mode: 'index',
                        intersect: false
                    }
                }
            }
        };

        const chart = new Chart(canvas, config);
        this.charts.set(canvasId, chart);
        return chart;
    }

    /**
     * Update chart data
     */
    updateChart(canvasId, newData) {
        const chart = this.charts.get(canvasId);
        if (!chart) return;

        if (Array.isArray(newData)) {
            chart.data.datasets[0].data = newData;
        } else if (newData.datasets) {
            chart.data = newData;
        }

        chart.update();
    }

    /**
     * Destroy a chart
     */
    destroyChart(canvasId) {
        const chart = this.charts.get(canvasId);
        if (chart) {
            chart.destroy();
            this.charts.delete(canvasId);
        }
    }

    /**
     * Destroy all charts
     */
    destroyAll() {
        this.charts.forEach(chart => chart.destroy());
        this.charts.clear();
    }

    /**
     * Get chart instance
     */
    getChart(canvasId) {
        return this.charts.get(canvasId);
    }
}

// Create global instance
const advancedCharts = new AdvancedCharts();

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = AdvancedCharts;
}

// Export globally
window.advancedCharts = advancedCharts;
