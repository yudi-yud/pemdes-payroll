import { useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const IDLE_TIMEOUT = 60 * 60 * 1000; // 1 hour in milliseconds

export const useAuthTimeout = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const timeoutRef = useRef(null);
  const warningTimeoutRef = useRef(null);

  // Reset the timeout
  const resetTimeout = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    if (warningTimeoutRef.current) {
      clearTimeout(warningTimeoutRef.current);
    }

    // Set timeout for logout after IDLE_TIMEOUT
    timeoutRef.current = setTimeout(() => {
      handleLogout();
    }, IDLE_TIMEOUT);

    // Optional: Show warning 5 minutes before logout
    warningTimeoutRef.current = setTimeout(() => {
      showWarning();
    }, IDLE_TIMEOUT - 5 * 60 * 1000);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const showWarning = () => {
    // You could show a modal here instead
    console.warn('Anda akan logout dalam 5 menit karena tidak ada aktivitas');
  };

  useEffect(() => {
    // List of events to track user activity
    const events = [
      'mousedown',
      'mousemove',
      'keypress',
      'scroll',
      'touchstart',
      'click',
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

  return { resetTimeout };
};
