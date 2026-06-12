import React, { useState } from 'react';
import { authApi, usersApi } from '../api/gateway.api';
import { useAuthStore } from '../store/authStore';

interface LoginPageProps {
  onSuccess: () => void;
}

type Tab = 'login' | 'register';

const LoginPage: React.FC<LoginPageProps> = ({ onSuccess }) => {
  const [tab, setTab] = useState<Tab>('login');

  // Login state
  const [loginEmail, setLoginEmail] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [loginLoading, setLoginLoading] = useState(false);
  const [loginError, setLoginError] = useState('');

  // Register state
  const [regEmail, setRegEmail] = useState('');
  const [regUsername, setRegUsername] = useState('');
  const [regFullName, setRegFullName] = useState('');
  const [regPhone, setRegPhone] = useState('');
  const [regPassword, setRegPassword] = useState('');
  const [regConfirm, setRegConfirm] = useState('');
  const [regLoading, setRegLoading] = useState(false);
  const [regError, setRegError] = useState('');
  const [regSuccess, setRegSuccess] = useState('');

  const { setAuth } = useAuthStore();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoginError('');
    setLoginLoading(true);
    try {
      const res = await authApi.login(loginEmail, loginPassword);
      const userId = res.user_id;
      if (!userId || !res.access_token || !res.refresh_token) {
        throw new Error('Invalid server response.');
      }
      const profileRes = await usersApi.getOne(userId);
      const u = profileRes.user;
      setAuth(
        {
          user_id: userId,
          email: u?.email || loginEmail,
          username: u?.username || '',
          full_name: u?.full_name || '',
          phone: u?.phone || '',
          roles: u?.roles || [],
        },
        res.access_token,
        res.refresh_token
      );
      onSuccess();
    } catch (err) {
      setLoginError(err instanceof Error ? err.message : 'Sign in failed.');
    } finally {
      setLoginLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setRegError('');
    setRegSuccess('');
    if (regPassword !== regConfirm) {
      setRegError('Passwords do not match.');
      return;
    }
    if (regPassword.length < 8) {
      setRegError('Password must be at least 8 characters.');
      return;
    }
    setRegLoading(true);
    try {
      await authApi.register({
        email: regEmail,
        username: regUsername,
        full_name: regFullName,
        phone: regPhone,
        password: regPassword,
      });
      setRegSuccess('Registration successful! Please sign in.');
      setLoginEmail(regEmail);
      setTimeout(() => setTab('login'), 1200);
    } catch (err) {
      setRegError(err instanceof Error ? err.message : 'Registration failed.');
    } finally {
      setRegLoading(false);
    }
  };

  return (
    <div className="flex-1 flex items-center justify-center px-4 py-16 bg-background">
      <div className="w-full max-w-md">
        {/* Card */}
        <div className="bg-surface-container-lowest rounded-2xl shadow-[0_8px_32px_rgba(0,0,0,0.08)] border border-outline-variant/30 overflow-hidden">
          {/* Brand */}
          <div className="text-center pt-10 pb-6 px-8 border-b border-outline-variant/30">
            <div className="inline-flex items-center justify-center w-14 h-14 bg-primary-container/20 rounded-full mb-4">
              <span className="material-symbols-outlined text-primary text-3xl" style={{ fontVariationSettings: "'FILL' 1" }}>
                restaurant_menu
              </span>
            </div>
            <h1 className="font-headline-md text-2xl text-on-surface font-bold">LuxeBistro</h1>
            <p className="text-label-sm text-on-surface-variant mt-1">Customer Account</p>
          </div>

          {/* Tabs */}
          <div className="flex border-b border-outline-variant/30">
            {(['login', 'register'] as Tab[]).map((t) => (
              <button
                key={t}
                type="button"
                onClick={() => { setTab(t); setLoginError(''); setRegError(''); setRegSuccess(''); }}
                className={`flex-1 py-3 text-sm font-semibold transition-all ${
                  tab === t
                    ? 'text-primary border-b-2 border-primary bg-primary/5'
                    : 'text-on-surface-variant hover:text-on-surface'
                }`}
              >
                {t === 'login' ? 'Sign In' : 'Register'}
              </button>
            ))}
          </div>

          <div className="p-8">
            {/* Login Form */}
            {tab === 'login' && (
              <form onSubmit={handleLogin} className="space-y-5">
                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="l-email">
                    Email
                  </label>
                  <div className="relative">
                    <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg">mail</span>
                    <input
                      id="l-email"
                      type="email"
                      required
                      placeholder="email@example.com"
                      value={loginEmail}
                      onChange={(e) => setLoginEmail(e.target.value)}
                      className="w-full h-12 pl-10 pr-4 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="l-password">
                    Password
                  </label>
                  <div className="relative">
                    <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg">lock</span>
                    <input
                      id="l-password"
                      type="password"
                      required
                      placeholder="••••••••"
                      value={loginPassword}
                      onChange={(e) => setLoginPassword(e.target.value)}
                      className="w-full h-12 pl-10 pr-4 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                    />
                  </div>
                </div>

                {loginError && (
                  <p className="text-sm text-error bg-error/10 rounded-lg px-3 py-2">{loginError}</p>
                )}

                <button
                  type="submit"
                  disabled={loginLoading}
                  className="w-full h-12 bg-primary text-on-primary rounded-lg font-semibold text-sm hover:opacity-90 active:scale-[0.98] transition-all disabled:opacity-60 flex items-center justify-center gap-2"
                >
                  {loginLoading ? (
                    <span className="material-symbols-outlined animate-spin text-lg">progress_activity</span>
                  ) : (
                    <>
                      Sign In
                      <span className="material-symbols-outlined text-lg">login</span>
                    </>
)}
                </button>

                <p className="text-center text-sm text-on-surface-variant">
                  Don't have an account?{' '}
                  <button type="button" onClick={() => setTab('register')} className="text-primary font-semibold hover:underline">
                    Register now
                  </button>
                </p>
              </form>
            )}

            {/* Register Form */}
            {tab === 'register' && (
              <form onSubmit={handleRegister} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="r-fullname">
                      Full Name
                    </label>
                    <input
                      id="r-fullname"
                      type="text"
                      required
                      placeholder="John Doe"
                      value={regFullName}
                      onChange={(e) => setRegFullName(e.target.value)}
                      className="w-full h-11 px-3 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                    />
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="r-username">
                      Username
                    </label>
                    <input
                      id="r-username"
                      type="text"
                      required
                      placeholder="johndoe"
                      value={regUsername}
                      onChange={(e) => setRegUsername(e.target.value)}
                      className="w-full h-11 px-3 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="r-email">
                    Email
                  </label>
                  <input
                    id="r-email"
                    type="email"
                    required
                    placeholder="email@example.com"
                    value={regEmail}
                    onChange={(e) => setRegEmail(e.target.value)}
                    className="w-full h-11 px-3 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="r-phone">
                    Phone Number
                  </label>
                  <input
                    id="r-phone"
                    type="tel"
                    placeholder="090 123 4567"
                    value={regPhone}
                    onChange={(e) => setRegPhone(e.target.value)}
                    className="w-full h-11 px-3 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="r-password">
                      Password
                    </label>
                    <input
                      id="r-password"
                      type="password"
                      required
                      placeholder="≥ 8 characters"
                      value={regPassword}
                      onChange={(e) => setRegPassword(e.target.value)}
                      className="w-full h-11 px-3 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                    />
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-on-surface-variant uppercase tracking-wider" htmlFor="r-confirm">
                      Confirm
                    </label>
                    <input
                      id="r-confirm"
                      type="password"
                      required
                      placeholder="Re-enter"
                      value={regConfirm}
                      onChange={(e) => setRegConfirm(e.target.value)}
                      className="w-full h-11 px-3 bg-surface border border-outline-variant rounded-lg text-sm outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
                    />
                  </div>
                </div>

                {regError && (
                  <p className="text-sm text-error bg-error/10 rounded-lg px-3 py-2">{regError}</p>
                )}
                {regSuccess && (
                  <p className="text-sm text-[#2e7d32] bg-[#2e7d32]/10 rounded-lg px-3 py-2">{regSuccess}</p>
                )}

                <button
                  type="submit"
                  disabled={regLoading}
                  className="w-full h-12 bg-primary text-on-primary rounded-lg font-semibold text-sm hover:opacity-90 active:scale-[0.98] transition-all disabled:opacity-60 flex items-center justify-center gap-2"
                >
                  {regLoading ? (
                    <span className="material-symbols-outlined animate-spin text-lg">progress_activity</span>
                  ) : (
                    <>
                      Create Account
                      <span className="material-symbols-outlined text-lg">person_add</span>
                    </>
)}
                </button>

                <p className="text-center text-sm text-on-surface-variant">
                  Already have an account?{' '}
                  <button type="button" onClick={() => setTab('login')} className="text-primary font-semibold hover:underline">
                    Sign In
                  </button>
                </p>
              </form>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
