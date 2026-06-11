import React, { useState, useEffect } from 'react';
import { authApi, usersApi } from '../services/api';
import { useAdminAuthStore, hasAdminAccess } from '../store/adminAuthStore';

interface LoginPageProps {
  onSuccess: () => void;
}

export default function LoginPage({ onSuccess }: LoginPageProps) {
    const [showPassword, setShowPassword] = useState(false);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const { setAuth } = useAdminAuthStore();

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

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            const res = await authApi.login(email, password);
            if (!res.user_id || !res.access_token || !res.refresh_token) {
                throw new Error('Invalid server response.');
            }
            const profileRes = await usersApi.getOne(res.user_id);
            const u = profileRes.user;
            const roles = u?.roles ?? [];
            if (!hasAdminAccess(roles)) {
                throw new Error('Your account does not have admin access.');
            }
            setAuth(
                {
                    user_id: res.user_id,
                    email: u?.email || email,
                    username: u?.username || '',
                    full_name: u?.full_name || '',
                    roles,
                },
                res.access_token,
                res.refresh_token
            );
            onSuccess();
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Sign in failed.');
        } finally {
            setLoading(false);
        }
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

                            {/* Error */}
                            {error && (
                                <p className="text-sm text-error bg-error/10 rounded-lg px-3 py-2">{error}</p>
                            )}

                            {/* Submit Button */}
                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full h-14 bg-primary-container hover:bg-[#c29d2b] active:scale-[0.98] text-on-primary-container font-label-sm text-[12px] uppercase tracking-widest rounded-lg transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2 disabled:opacity-60"
                            >
                                {loading ? (
                                    <span className="material-symbols-outlined animate-spin text-[18px]">progress_activity</span>
                                ) : (
                                    <>
                                        <span>Sign In</span>
                                        <span className="material-symbols-outlined text-[18px]">login</span>
                                    </>
                                )}
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