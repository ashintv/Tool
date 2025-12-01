import React, { useState } from 'react'
import { Button } from "../components/ui/button"
import { Input } from "../components/ui/input"
import { Label } from "../components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card"

const Signup = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  })
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
    // Clear messages when user starts typing
    if (error) setError('')
    if (success) setSuccess('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')
    setSuccess('')

    // Client-side validation
    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match')
      setIsLoading(false)
      return
    }

    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters long')
      setIsLoading(false)
      return
    }

    try {
      const response = await fetch('/api/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: formData.username,
          email: formData.email,
          password: formData.password
        }),
      })

      const data = await response.json()

      if (response.ok) {
        setSuccess('Account created successfully! You can now log in.')
        // Reset form
        setFormData({
          username: '',
          email: '',
          password: '',
          confirmPassword: ''
        })
        // Redirect to login after 2 seconds
        setTimeout(() => {
          window.location.href = '/login'
        }, 2000)
      } else {
        setError(data.error || 'Registration failed')
      }
    } catch (err) {
      setError('Network error. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex items-center bg-white px-20 justify-center  py-12 rounded-xl">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-light text-gray-900 mb-2">Create Account</h1>
          <p className="text-gray-600">
            Enter your information to get started
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="border border-gray-300 text-gray-700 px-4 py-3 text-sm">
              {error}
            </div>
          )}

          {success && (
            <div className="border border-gray-300 text-gray-700 px-4 py-3 text-sm">
              {success}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="username" className="text-gray-700 font-normal">Username</Label>
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
            <Label htmlFor="email" className="text-gray-700 font-normal">Email</Label>
            <Input
              id="email"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              required
              placeholder="Enter your email"
              disabled={isLoading}
              className="border-gray-300 focus:border-gray-500"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password" className="text-gray-700 font-normal">Password</Label>
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

          <div className="space-y-2">
            <Label htmlFor="confirmPassword" className="text-gray-700 font-normal">Confirm Password</Label>
            <Input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
              placeholder="Confirm your password"
              disabled={isLoading}
              className="border-gray-300 focus:border-gray-500"
            />
          </div>

          <Button
            type="submit"
            className="w-full bg-gray-900 hover:bg-gray-800 text-white"
            disabled={isLoading}
          >
            {isLoading ? 'Creating Account...' : 'Create Account'}
          </Button>

          <div className="text-center text-sm">
            <span className="text-gray-600">Already have an account? </span>
            <a href="/login" className="text-gray-900 hover:underline">
              Sign in
            </a>
          </div>
        </form>
      </div>
    </div>
  )
}

export default Signup
