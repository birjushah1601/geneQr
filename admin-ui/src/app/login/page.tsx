"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import OTPInput from '@/components/auth/OTPInput';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081';

type LoginStep = 'identifier' | 'otp' | 'password';

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();

  const [step, setStep] = useState<LoginStep>('identifier');
  const [identifier, setIdentifier] = useState('');
  const [otp, setOtp] = useState('');
  const [password, setPassword] = useState('');
  const [usePassword, setUsePassword] = useState(true); // Default to password instead of OTP
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [otpSentTo, setOtpSentTo] = useState('');
  const [expiresIn, setExpiresIn] = useState(0);

  // Timer for OTP expiry
  const [timeRemaining, setTimeRemaining] = useState(0);

  React.useEffect(() => {
    if (timeRemaining > 0) {
      const timer = setTimeout(() => setTimeRemaining(timeRemaining - 1), 1000);
      return () => clearTimeout(timer);
    }
  }, [timeRemaining]);

  // Send OTP
  const handleSendOTP = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/v1/auth/send-otp`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ identifier }),
      });

      if (response.ok) {
        const data = await response.json();
        setOtpSentTo(data.sent_to);
        setExpiresIn(data.expires_in);
        setTimeRemaining(data.expires_in);
        setStep('otp');
      } else {
        const data = await response.json();
        setError(data.error?.message || 'Failed to send OTP');
      }
    } catch (err) {
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Verify OTP
  const handleVerifyOTP = async (e: React.FormEvent) => {
    e.preventDefault();
    if (otp.length !== 6) {
      setError('Please enter a complete 6-digit code');
      return;
    }

    setError('');
    setIsLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/v1/auth/verify-otp`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ identifier, code: otp }),
      });

      if (response.ok) {
        const data = await response.json();
        await login(data.access_token, data.refresh_token);
        router.push('/dashboard');
      } else {
        const data = await response.json();
        setError(data.error?.message || 'Invalid OTP');
        setOtp('');
      }
    } catch (err) {
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Login with password
  const handlePasswordLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/v1/auth/login-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ identifier, password }),
      });

      const data = await response.json();
      console.log('Login response:', { status: response.status, data });
      
      if (response.ok) {
        if (data.access_token && data.refresh_token) {
          await login(data.access_token, data.refresh_token);
          router.push('/dashboard');
        } else {
          console.error('Missing tokens in response:', data);
          setError('Invalid server response. Please check backend logs.');
        }
      } else {
        setError(data.error?.message || data.message || 'Invalid credentials');
      }
    } catch (err) {
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Resend OTP
  const handleResendOTP = async () => {
    setOtp('');
    await handleSendOTP(new Event('submit') as any);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-xl p-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              ABY-MED Platform
            </h1>
            <p className="text-gray-600">
              {step === 'identifier' && 'Sign in to your account'}
              {step === 'otp' && 'Enter verification code'}
              {step === 'password' && 'Enter your password'}
            </p>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          {/* Step 1: Enter Identifier */}
          {step === 'identifier' && (
            <form onSubmit={usePassword ? (e) => { e.preventDefault(); setStep('password'); } : handleSendOTP}>
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email or Phone Number
                </label>
                <input
                  type="text"
                  value={identifier}
                  onChange={(e) => setIdentifier(e.target.value)}
                  placeholder="user@example.com or +1234567890"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                  disabled={isLoading}
                />
              </div>

              <button
                type="submit"
                disabled={isLoading || !identifier}
                className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
              >
                {isLoading ? 'Loading...' : (usePassword ? 'Continue' : 'Send OTP')}
              </button>

              <div className="mt-6 text-center">
                <button
                  type="button"
                  onClick={() => setUsePassword(!usePassword)}
                  className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                >
                  {usePassword ? 'Use OTP instead' : 'Use password instead'}
                </button>
              </div>

              <div className="mt-6 text-center">
                <a href="/register" className="text-sm text-gray-600 hover:text-gray-900">
                  Don't have an account? <span className="text-blue-600 font-medium">Sign up</span>
                </a>
              </div>
            </form>
          )}

          {/* Step 2: Enter OTP */}
          {step === 'otp' && (
            <form onSubmit={handleVerifyOTP}>
              <div className="mb-6">
                <p className="text-sm text-gray-600 text-center mb-4">
                  Code sent to <span className="font-semibold">{otpSentTo}</span>
                </p>

                <OTPInput
                  value={otp}
                  onChange={setOtp}
                  disabled={isLoading}
                />

                {timeRemaining > 0 && (
                  <p className="text-xs text-gray-500 text-center mt-3">
                    Code expires in {Math.floor(timeRemaining / 60)}:{(timeRemaining % 60).toString().padStart(2, '0')}
                  </p>
                )}
              </div>

              <button
                type="submit"
                disabled={isLoading || otp.length !== 6}
                className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
              >
                {isLoading ? 'Verifying...' : 'Verify & Login'}
              </button>

              <div className="mt-6 flex justify-between text-sm">
                <button
                  type="button"
                  onClick={() => setStep('identifier')}
                  className="text-gray-600 hover:text-gray-900"
                  disabled={isLoading}
                >
                  ← Back
                </button>
                <button
                  type="button"
                  onClick={handleResendOTP}
                  className="text-blue-600 hover:text-blue-700 font-medium"
                  disabled={isLoading || timeRemaining > expiresIn - 60}
                >
                  Resend code
                </button>
              </div>

              <div className="mt-6 text-center">
                <button
                  type="button"
                  onClick={() => { setStep('password'); setUsePassword(true); }}
                  className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                  disabled={isLoading}
                >
                  Use password instead
                </button>
              </div>
            </form>
          )}

          {/* Step 3: Enter Password */}
          {step === 'password' && (
            <form onSubmit={handlePasswordLogin}>
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Password
                </label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                  disabled={isLoading}
                />
              </div>

              <button
                type="submit"
                disabled={isLoading || !password}
                className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
              >
                {isLoading ? 'Logging in...' : 'Login'}
              </button>

              <div className="mt-6 flex justify-between text-sm">
                <button
                  type="button"
                  onClick={() => setStep('identifier')}
                  className="text-gray-600 hover:text-gray-900"
                  disabled={isLoading}
                >
                  ← Back
                </button>
                <a href="/forgot-password" className="text-blue-600 hover:text-blue-700 font-medium">
                  Forgot password?
                </a>
              </div>

              <div className="mt-6 text-center">
                <button
                  type="button"
                  onClick={() => { setStep('identifier'); setUsePassword(false); }}
                  className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                  disabled={isLoading}
                >
                  Use OTP instead
                </button>
              </div>
            </form>
          )}
        </div>

        {/* Footer */}
        <p className="text-center text-sm text-gray-600 mt-8">
          Secure authentication powered by ABY-MED
        </p>
      </div>
    </div>
  );
}
