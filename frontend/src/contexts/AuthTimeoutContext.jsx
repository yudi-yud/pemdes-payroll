import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { useAuth } from './AuthContext';

const AuthTimeoutContext = createContext(null);

export const useAuthTimeout = () => useContext(AuthTimeoutContext);

const IDLE_TIMEOUT = 60 * 60 * 1000; // 1 hour
const WARNING_TIME = 5 * 60 * 1000; // 5 minutes before timeout

export const AuthTimeoutProvider = ({ children }) => {
  const { logout } = useAuth();
  const timeoutRef = useRef(null);
  const warningTimeoutRef = useRef(null);
  const [showWarning, setShowWarning] = useState(false);
  const [timeLeft, setTimeLeft] = useState(0);

  const handleLogout = () => {
    logout();
    setShowWarning(false);
    window.location.href = '/login';
  };

  const resetTimeout = () => {
    // Clear existing timeouts
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    if (warningTimeoutRef.current) {
      clearTimeout(warningTimeoutRef.current);
    }
    setShowWarning(false);

    // Set warning timeout (5 minutes before logout)
    warningTimeoutRef.current = setTimeout(() => {
      setShowWarning(true);
      setTimeLeft(WARNING_TIME / 1000); // 5 minutes in seconds
    }, IDLE_TIMEOUT - WARNING_TIME);

    // Set logout timeout
    timeoutRef.current = setTimeout(() => {
      handleLogout();
    }, IDLE_TIMEOUT);
  };

  // Countdown for warning
  useEffect(() => {
    let interval;
    if (showWarning && timeLeft > 0) {
      interval = setInterval(() => {
        setTimeLeft((prev) => {
          if (prev <= 1) {
            handleLogout();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [showWarning, timeLeft]);

  useEffect(() => {
    // Events to track user activity
    const events = [
      'mousedown',
      'mousemove',
      'keypress',
      'scroll',
      'touchstart',
      'click',
      'keydown',
    ];

    // Add event listeners
    events.forEach((event) => {
      document.addEventListener(event, resetTimeout);
    });

    // Initial timeout setup
    resetTimeout();

    // Cleanup
    return () => {
      events.forEach((event) => {
        document.removeEventListener(event, resetTimeout);
      });
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      if (warningTimeoutRef.current) {
        clearTimeout(warningTimeoutRef.current);
      }
    };
  }, []);

  const extendSession = () => {
    resetTimeout();
  };

  return (
    <AuthTimeoutContext.Provider value={{ showWarning, timeLeft, extendSession }}>
      {children}
      {showWarning && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl p-6 max-w-md w-full mx-4">
            <div className="text-center">
              <div className="text-5xl mb-4">⏰</div>
              <h3 className="text-xl font-bold text-gray-800 mb-2">Sesi Akan Berakhir</h3>
              <p className="text-gray-600 mb-4">
                Anda akan otomatis logout dalam <span className="font-bold text-red-600">{Math.floor(timeLeft / 60)}:{(timeLeft % 60).toString().padStart(2, '0')}</span> menit karena tidak ada aktivitas.
              </p>
              <button
                onClick={extendSession}
                className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors font-medium"
              >
                Perpanjang Sesi
              </button>
            </div>
          </div>
        </div>
      )}
    </AuthTimeoutContext.Provider>
  );
};
