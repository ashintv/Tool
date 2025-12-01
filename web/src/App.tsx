import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Button } from "./components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./components/ui/card";
import Login from "./pages/login";
import Signup from "./pages/signup";
import Dashboard from "./pages/dashboard";
import "./index.css";

// Landing Page Component
const LandingPage = () => {
  return (
    <div className=" bg-white flex items-center justify-center rounded-2xl py-12 px-4">
      <div className="max-w-2xl mx-auto text-center">
        <div className="mb-12">
          <h1 className="text-5xl font-light text-gray-900 mb-6">
            Tool
          </h1>
          <p className="text-lg text-gray-600 mb-12 leading-relaxed">
            A simple platform for managing your digital infrastructure.
            Please sign in or create an account to continue.
          </p>
        </div>

        <div className="space-y-8">
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button
              variant="outline"
              className="border-gray-300 text-gray-700 hover:bg-gray-50 px-8 py-3"
              onClick={() => window.location.href = '/login'}
            >
              Sign In
            </Button>
            <Button
              className="bg-gray-900 hover:bg-gray-800 text-white px-8 py-3"
              onClick={() => window.location.href = '/signup'}
            >
              Create Account
            </Button>
          </div>
        </div>

        <div className="mt-20 space-y-8 text-left max-w-lg mx-auto">
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Fast Setup</h3>
            <p className="text-gray-600 text-sm leading-relaxed">
              Get up and running in minutes with our straightforward interface
            </p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Secure</h3>
            <p className="text-gray-600 text-sm leading-relaxed">
              Your data is protected with industry-standard security practices
            </p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Reliable</h3>
            <p className="text-gray-600 text-sm leading-relaxed">
              Dependable infrastructure with consistent uptime and performance
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

// Protected Route Component
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const token = localStorage.getItem('token');
  return token ? <>{children}</> : <Navigate to="/login" />;
};

export function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        {/* Redirect any unknown routes to home */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;
