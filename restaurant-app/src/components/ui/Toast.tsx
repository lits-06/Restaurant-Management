import React, { useEffect } from 'react';

interface ToastProps {
  message: string;
  onClose: () => void;
  duration?: number;
}

const Toast: React.FC<ToastProps> = ({ message, onClose, duration = 4000 }) => {
  useEffect(() => {
    const t = setTimeout(onClose, duration);
    return () => clearTimeout(t);
  }, [message, onClose, duration]);

  return (
    <div className="fixed top-6 right-6 z-50 flex items-center gap-3 bg-error text-on-error px-5 py-3 rounded-xl shadow-lg animate-in slide-in-from-top-4 fade-in duration-300 max-w-sm w-max">
      <span className="material-symbols-outlined text-xl flex-shrink-0">error</span>
      <span className="text-sm font-medium leading-snug">{message}</span>
      <button
        type="button"
        onClick={onClose}
        className="flex-shrink-0 opacity-80 hover:opacity-100 transition-opacity ml-1"
        aria-label="Close"
      >
        <span className="material-symbols-outlined text-lg">close</span>
      </button>
    </div>
  );
};

export default Toast;
