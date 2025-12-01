import React, { useState } from "react";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card";

const Login = () => {
  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
    // Clear error when user starts typing
    if (error) setError("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (response.ok) {
        // Store token in localStorage
        localStorage.setItem("token", data.token);
        localStorage.setItem("user_id", data.user_id);

        // Redirect to dashboard or home page
        window.location.href = "/dashboard";
      } else {
        setError(data.error || "Login failed");
      }
    } catch (err) {
      setError("Network error. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex items-center bg-white px-20 justify-center  py-12 rounded-xl">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-light text-gray-900 mb-2">Sign In</h1>
          <p className="text-gray-600">Enter your credentials to continue</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && <div className="border border-gray-300 text-gray-700 px-4 py-3 text-sm">{error}</div>}

          <div className="space-y-2">
            <Label htmlFor="username" className="text-gray-700 font-normal">
              Username
            </Label>
            <Input
              id="username"
              name="username"
              type="text"
              value={formData.username}
              onChange={handleChange}
              required
              placeholder="Enter your username"
              disabled={isLoading}
              className="border-gray-300 focus:border-gray-500"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password" className="text-gray-700 font-normal">
              Password
            </Label>
            <Input
              id="password"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleChange}
              required
              placeholder="Enter your password"
              disabled={isLoading}
              className="border-gray-300 focus:border-gray-500"
            />
          </div>

          <Button type="submit" className="w-full bg-gray-900 hover:bg-gray-800 text-white" disabled={isLoading}>
            {isLoading ? "Signing in..." : "Sign In"}
          </Button>

          <div className="text-center text-sm">
            <span className="text-gray-600">Don't have an account? </span>
            <a href="/signup" className="text-gray-900 hover:underline">
              Sign up
            </a>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
