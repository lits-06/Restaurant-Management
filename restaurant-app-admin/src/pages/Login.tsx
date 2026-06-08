import React, { useState, useEffect } from 'react';

export default function LoginPage() {
    const [showPassword, setShowPassword] = useState(false);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(false);

    // Hiệu ứng Parallax nền khi di chuột
    useEffect(() => {
        const handleMouseMove: EventListener = (event) => {
            const e = event as MouseEvent;
            const moveX = (e.clientX - window.innerWidth / 2) * 0.01;
            const moveY = (e.clientY - window.innerHeight / 2) * 0.01;
            document.body.style.backgroundPosition = `${moveX}px ${moveY}px`;
        };

        document.addEventListener('mousemove', handleMouseMove);
        return () => {
            document.removeEventListener('mousemove', handleMouseMove);
        };
    }, []);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        // Xử lý logic đăng nhập tại đây
        console.log({ email, password, rememberMe });
    };

    return (
        <div className="bg-surface text-on-surface min-h-screen flex flex-col font-body-md">
            <main className="flex-grow flex items-center justify-center px-margin-mobile md:px-0 py-12 relative overflow-hidden">
                {/* Decorative Ambient Element */}
                <div className="absolute top-[-10%] right-[-10%] w-96 h-96 bg-primary-container opacity-[0.03] rounded-full blur-3xl"></div>
                <div className="absolute bottom-[-10%] left-[-10%] w-96 h-96 bg-primary-container opacity-[0.03] rounded-full blur-3xl"></div>

                <div className="w-full max-w-[440px] z-10">
                    {/* Login Card */}
                    <div className="bg-surface-container-lowest login-card-shadow rounded-xl p-8 md:p-10 border border-outline-variant/30">
                        {/* Brand Anchor */}
                        <div className="text-center mb-10">
                            <div className="inline-flex items-center justify-center w-16 h-16 bg-surface-container-low rounded-full mb-4 border border-outline-variant/20">
                                <span className="material-symbols-outlined text-primary text-3xl" style={{ fontVariationSettings: "'FILL' 1" }}>
                                    restaurant_menu
                                </span>
                            </div>
                            <h1 className="font-headline-md text-[24px] text-on-surface tracking-tight">LuxeBistro</h1>
                            <p className="font-label-sm text-[12px] text-secondary uppercase tracking-widest mt-1">Admin Portal</p>
                        </div>

                        {/* Form */}
                        <form onSubmit={handleSubmit} className="space-y-6">
                            {/* Email Field */}
                            <div className="space-y-2">
                                <label className="font-label-sm text-[12px] text-on-surface-variant block px-1" htmlFor="email">
                                    Email Address
                                </label>
                                <div className="relative">
                                    <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-secondary text-[20px]">
                                        mail
                                    </span>
                                    <input
                                        type="email"
                                        id="email"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        placeholder="name@luxebistro.com"
                                        required
                                        className="w-full h-14 pl-12 pr-4 bg-surface-container-low border border-outline-variant rounded-lg font-body-md text-on-surface placeholder:text-secondary/50 outline-none transition-all input-focus-glow"
                                    />
                                </div>
                            </div>

                            {/* Password Field */}
                            <div className="space-y-2">
                                <label className="font-label-sm text-[12px] text-on-surface-variant block px-1" htmlFor="password">
                                    Password
                                </label>
                                <div className="relative">
                                    <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-secondary text-[20px]">
                                        lock
                                    </span>
                                    <input
                                        type={showPassword ? 'text' : 'password'}
                                        id="password"
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        placeholder="••••••••"
                                        required
                                        className="w-full h-14 pl-12 pr-12 bg-surface-container-low border border-outline-variant rounded-lg font-body-md text-on-surface placeholder:text-secondary/50 outline-none transition-all input-focus-glow"
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowPassword(!showPassword)}
                                        className="absolute right-4 top-1/2 -translate-y-1/2 text-secondary hover:text-primary transition-colors"
                                    >
                                        <span className="material-symbols-outlined text-[20px]">
                                            {showPassword ? 'visibility_off' : 'visibility'}
                                        </span>
                                    </button>
                                </div>
                            </div>

                            {/* Utilities */}
                            <div className="flex items-center justify-between py-1">
                                <label className="flex items-center cursor-pointer group">
                                    <input
                                        type="checkbox"
                                        checked={rememberMe}
                                        onChange={(e) => setRememberMe(e.target.checked)}
                                        className="w-5 h-5 rounded border-outline-variant text-primary focus:ring-primary/20 cursor-pointer"
                                    />
                                    <span className="ml-3 font-body-md text-[14px] text-on-surface-variant group-hover:text-on-surface transition-colors">
                                        Remember me
                                    </span>
                                </label>
                                <a className="font-label-sm text-[12px] text-primary hover:text-on-primary-fixed-variant transition-colors underline decoration-primary/30 underline-offset-4" href="#forgot">
                                    Forgot Password?
                                </a>
                            </div>

                            {/* Submit Button */}
                            <button
                                type="submit"
                                className="w-full h-14 bg-primary-container hover:bg-[#c29d2b] active:scale-[0.98] text-on-primary-container font-label-sm text-[12px] uppercase tracking-widest rounded-lg transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2"
                            >
                                <span>Sign In</span>
                                <span className="material-symbols-outlined text-[18px]">login</span>
                            </button>
                        </form>

                        {/* Security Note */}
                        <div className="mt-8 pt-8 border-t border-outline-variant/30 flex items-start gap-3 opacity-60">
                            <span className="material-symbols-outlined text-secondary text-[18px]">verified_user</span>
                            <p className="font-body-md text-[12px] leading-relaxed text-secondary">
                                This session is encrypted. By signing in, you agree to our security protocols and administrative data handling policies.
                            </p>
                        </div>
                    </div>

                    {/* Technical Assistance Support CTA */}
                    <div className="text-center mt-6">
                        <p className="font-body-md text-[14px] text-secondary">
                            Need technical assistance?{' '}
                            <a className="text-on-surface font-semibold hover:text-primary transition-colors" href="mailto:support@luxebistro.com">
                                Contact Support
                            </a>
                        </p>
                    </div>
                </div>
            </main>

            {/* Footer */}
            <footer className="bg-surface border-t border-outline-variant py-base">
                <div className="flex flex-col md:flex-row justify-between items-center w-full px-margin-desktop max-w-container-max mx-auto space-y-4 md:space-y-0">
                    <div className="font-label-sm text-[12px] font-bold text-on-surface">LuxeBistro Informed Hospitality</div>
                    <div className="text-secondary font-label-sm text-[12px] text-center md:text-left">
                        © 2024 LuxeBistro Informed Hospitality. All rights reserved.
                    </div>
                    <div className="flex gap-6">
                        <a className="font-label-sm text-[12px] text-secondary hover:text-primary transition-colors" href="#privacy">Privacy Policy</a>
                        <a className="font-label-sm text-[12px] text-secondary hover:text-primary transition-colors" href="#terms">Terms of Service</a>
                        <a className="font-label-sm text-[12px] text-secondary hover:text-primary transition-colors" href="#support">Support</a>
                    </div>
                </div>
            </footer>
        </div>
    );
}