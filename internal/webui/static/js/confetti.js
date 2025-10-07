// Confetti Celebration System for Sloth Runner
// Lightweight, performant confetti animations for success celebrations

class ConfettiCelebration {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.particles = [];
        this.animationId = null;
        this.isActive = false;
    }

    init() {
        if (this.canvas) return; // Already initialized

        this.canvas = document.createElement('canvas');
        this.canvas.id = 'confetti-canvas';
        this.canvas.style.cssText = `
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            pointer-events: none;
            z-index: 99999;
        `;
        document.body.appendChild(this.canvas);

        this.ctx = this.canvas.getContext('2d');
        this.resize();

        window.addEventListener('resize', () => this.resize());
    }

    resize() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }

    createParticle(x, y) {
        const colors = [
            '#4F46E5', // Primary
            '#10B981', // Success
            '#F59E0B', // Warning
            '#3B82F6', // Info
            '#EC4899', // Pink
            '#8B5CF6', // Purple
            '#14B8A6', // Teal
            '#F97316', // Orange
        ];

        return {
            x: x || Math.random() * this.canvas.width,
            y: y || -10,
            size: Math.random() * 8 + 4,
            speedX: Math.random() * 6 - 3,
            speedY: Math.random() * 3 + 2,
            color: colors[Math.floor(Math.random() * colors.length)],
            rotation: Math.random() * 360,
            rotationSpeed: Math.random() * 10 - 5,
            gravity: 0.3,
            opacity: 1,
            shape: Math.random() > 0.5 ? 'circle' : 'square'
        };
    }

    update() {
        if (!this.isActive) return;

        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

        // Update and draw particles
        for (let i = this.particles.length - 1; i >= 0; i--) {
            const p = this.particles[i];

            // Update position
            p.x += p.speedX;
            p.y += p.speedY;
            p.speedY += p.gravity;
            p.rotation += p.rotationSpeed;

            // Fade out as it falls
            if (p.y > this.canvas.height * 0.7) {
                p.opacity -= 0.02;
            }

            // Remove if off screen or fully transparent
            if (p.y > this.canvas.height + 50 || p.opacity <= 0) {
                this.particles.splice(i, 1);
                continue;
            }

            // Draw particle
            this.ctx.save();
            this.ctx.translate(p.x, p.y);
            this.ctx.rotate((p.rotation * Math.PI) / 180);
            this.ctx.globalAlpha = p.opacity;

            if (p.shape === 'circle') {
                this.ctx.beginPath();
                this.ctx.arc(0, 0, p.size / 2, 0, Math.PI * 2);
                this.ctx.fillStyle = p.color;
                this.ctx.fill();
            } else {
                this.ctx.fillStyle = p.color;
                this.ctx.fillRect(-p.size / 2, -p.size / 2, p.size, p.size);
            }

            this.ctx.restore();
        }

        // Continue animation if particles remain
        if (this.particles.length > 0) {
            this.animationId = requestAnimationFrame(() => this.update());
        } else {
            this.stop();
        }
    }

    /**
     * Trigger confetti burst
     * @param {Object} options - Confetti options
     */
    burst(options = {}) {
        const {
            particleCount = 100,
            origin = { x: 0.5, y: 0.5 },
            spread = 360,
            startVelocity = 45,
            decay = 0.9,
            scalar = 1
        } = options;

        this.init();

        const originX = origin.x * this.canvas.width;
        const originY = origin.y * this.canvas.height;

        for (let i = 0; i < particleCount; i++) {
            const angle = (spread / particleCount) * i - spread / 2;
            const velocity = startVelocity * (0.5 + Math.random() * 0.5);

            const particle = this.createParticle(originX, originY);
            particle.speedX = Math.cos((angle * Math.PI) / 180) * velocity * scalar;
            particle.speedY = Math.sin((angle * Math.PI) / 180) * velocity * scalar;
            particle.gravity = 0.5;

            this.particles.push(particle);
        }

        if (!this.isActive) {
            this.isActive = true;
            this.update();
        }
    }

    /**
     * Continuous confetti rain
     * @param {number} duration - Duration in milliseconds
     * @param {Object} options - Rain options
     */
    rain(duration = 3000, options = {}) {
        const {
            particleRate = 10,
            colors = null
        } = options;

        this.init();
        this.isActive = true;

        const interval = setInterval(() => {
            for (let i = 0; i < particleRate; i++) {
                this.particles.push(this.createParticle());
            }
        }, 100);

        setTimeout(() => {
            clearInterval(interval);
        }, duration);

        this.update();
    }

    /**
     * Confetti cannon from specific position
     * @param {Object} options - Cannon options
     */
    cannon(options = {}) {
        const {
            position = { x: 0.5, y: 1 },
            particleCount = 150,
            angle = 270,
            spread = 45
        } = options;

        this.init();

        const originX = position.x * this.canvas.width;
        const originY = position.y * this.canvas.height;

        for (let i = 0; i < particleCount; i++) {
            const spreadAngle = (Math.random() - 0.5) * spread;
            const finalAngle = angle + spreadAngle;
            const velocity = Math.random() * 15 + 10;

            const particle = this.createParticle(originX, originY);
            particle.speedX = Math.cos((finalAngle * Math.PI) / 180) * velocity;
            particle.speedY = Math.sin((finalAngle * Math.PI) / 180) * velocity;
            particle.gravity = 0.4;

            this.particles.push(particle);
        }

        if (!this.isActive) {
            this.isActive = true;
            this.update();
        }
    }

    /**
     * Firework explosion effect
     * @param {Object} options - Firework options
     */
    firework(options = {}) {
        const {
            position = { x: Math.random(), y: Math.random() * 0.5 },
            particleCount = 50,
            colors = null
        } = options;

        this.init();

        const originX = position.x * this.canvas.width;
        const originY = position.y * this.canvas.height;

        for (let i = 0; i < particleCount; i++) {
            const angle = (360 / particleCount) * i;
            const velocity = Math.random() * 10 + 8;

            const particle = this.createParticle(originX, originY);
            particle.speedX = Math.cos((angle * Math.PI) / 180) * velocity;
            particle.speedY = Math.sin((angle * Math.PI) / 180) * velocity;
            particle.gravity = 0.2;
            particle.size = Math.random() * 4 + 2;

            this.particles.push(particle);
        }

        if (!this.isActive) {
            this.isActive = true;
            this.update();
        }
    }

    /**
     * Multiple fireworks
     * @param {number} count - Number of fireworks
     * @param {number} delay - Delay between fireworks (ms)
     */
    fireworks(count = 5, delay = 200) {
        for (let i = 0; i < count; i++) {
            setTimeout(() => {
                this.firework({
                    position: {
                        x: Math.random() * 0.6 + 0.2,
                        y: Math.random() * 0.4 + 0.1
                    }
                });
            }, i * delay);
        }
    }

    /**
     * Success celebration - combination of effects
     */
    celebrate() {
        // Central burst
        this.burst({
            particleCount: 150,
            spread: 180,
            origin: { x: 0.5, y: 0.5 }
        });

        // Side cannons
        setTimeout(() => {
            this.cannon({
                position: { x: 0, y: 1 },
                particleCount: 100,
                angle: 270,
                spread: 60
            });
        }, 100);

        setTimeout(() => {
            this.cannon({
                position: { x: 1, y: 1 },
                particleCount: 100,
                angle: 270,
                spread: 60
            });
        }, 200);

        // Top fireworks
        setTimeout(() => {
            this.fireworks(3, 150);
        }, 400);
    }

    /**
     * Workflow success celebration
     */
    workflowSuccess() {
        this.burst({
            particleCount: 100,
            spread: 120,
            origin: { x: 0.5, y: 0.6 }
        });

        setTimeout(() => {
            this.rain(2000, { particleRate: 5 });
        }, 500);
    }

    /**
     * Agent connected celebration
     */
    agentConnected() {
        this.burst({
            particleCount: 50,
            spread: 90,
            origin: { x: 0.5, y: 0.3 }
        });
    }

    /**
     * Task completed celebration
     */
    taskCompleted() {
        this.firework({
            position: { x: 0.5, y: 0.4 },
            particleCount: 40
        });
    }

    stop() {
        this.isActive = false;
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }
        this.particles = [];
        if (this.ctx) {
            this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        }
    }

    destroy() {
        this.stop();
        if (this.canvas && this.canvas.parentNode) {
            this.canvas.parentNode.removeChild(this.canvas);
        }
        this.canvas = null;
        this.ctx = null;
    }
}

// Create global instance
const confetti = new ConfettiCelebration();

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ConfettiCelebration;
}
