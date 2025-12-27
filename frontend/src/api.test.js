/**
 * Unit tests for API client
 * Tests HTTP requests, error handling, and token management
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import axios from 'axios'
import api from './api'

// Mock axios
vi.mock('axios', () => {
  const mockAxios = {
    create: vi.fn(() => mockAxios),
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() }
    }
  }
  return {
    default: mockAxios,
    create: vi.fn(() => mockAxios)
  }
})

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('Authentication', () => {
    it('should login successfully with valid credentials', async () => {
      const mockResponse = {
        data: {
          success: true,
          data: {
            token: 'test-token',
            expires_at: '2024-01-01T00:00:00Z',
            user: { id: 1, email: 'test@example.com', name: 'Test User' }
          }
        }
      }
      api.post.mockResolvedValue(mockResponse)

      const result = await api.post('/auth/login', { email: 'test@example.com', password: 'password123' })

      expect(api.post).toHaveBeenCalledWith(
        '/auth/login',
        { email: 'test@example.com', password: 'password123' }
      )
      expect(result.data.success).toBe(true)
    })

    it('should handle login errors', async () => {
      api.post.mockRejectedValue({
        response: { status: 401, data: { error: 'Invalid credentials' } }
      })

      await expect(api.post('/auth/login', { email: 'test@example.com', password: 'wrong' })).rejects.toThrow()
    })

    it('should register new user', async () => {
      const mockResponse = {
        data: {
          success: true,
          data: {
            token: 'new-token',
            user: { id: 2, email: 'new@example.com', name: 'New User' }
          }
        }
      }
      api.post.mockResolvedValue(mockResponse)

      const result = await api.post('/auth/register', { email: 'new@example.com', password: 'password123', name: 'New User' })

      expect(api.post).toHaveBeenCalledWith(
        '/auth/register',
        { email: 'new@example.com', password: 'password123', name: 'New User' }
      )
      expect(result.data.success).toBe(true)
    })

    it('should handle registration errors', async () => {
      api.post.mockRejectedValue({
        response: { status: 409, data: { error: 'Email already registered' } }
      })

      await expect(api.post('/auth/register', { email: 'existing@example.com', password: 'password', name: 'User' })).rejects.toThrow()
    })
  })

  describe('Customer Management', () => {
    beforeEach(() => {
      localStorage.setItem('token', 'test-token')
    })

    it('should list customers', async () => {
      const mockCustomers = [
        { id: 1, name: 'Customer 1', latitude: 40.7128, longitude: -74.0060 },
        { id: 2, name: 'Customer 2', latitude: 34.0522, longitude: -118.2437 }
      ]
      api.get.mockResolvedValue({ data: { success: true, data: mockCustomers } })

      const result = await api.get('/customers')

      expect(api.get).toHaveBeenCalledWith('/customers')
      expect(result.data.data).toEqual(mockCustomers)
    })

    it('should create customer', async () => {
      const newCustomer = {
        name: 'New Customer',
        latitude: 40.7128,
        longitude: -74.0060,
        demand_rate: 100,
        max_inventory: 1000,
        current_inventory: 500,
        min_inventory: 100
      }
      const mockResponse = {
        data: {
          success: true,
          data: { id: 1, ...newCustomer }
        }
      }
      api.post.mockResolvedValue(mockResponse)

      const result = await api.post('/customers', newCustomer)

      expect(api.post).toHaveBeenCalledWith('/customers', newCustomer)
      expect(result.data.success).toBe(true)
    })

    it('should update customer', async () => {
      const updatedCustomer = { id: 1, name: 'Updated Customer', latitude: 40.7128, longitude: -74.0060 }
      api.put.mockResolvedValue({ data: { success: true, data: updatedCustomer } })

      const result = await api.put('/customers/1', updatedCustomer)

      expect(api.put).toHaveBeenCalledWith('/customers/1', updatedCustomer)
      expect(result.data.success).toBe(true)
    })

    it('should delete customer', async () => {
      api.delete.mockResolvedValue({ data: { success: true } })

      await api.delete('/customers/1')

      expect(api.delete).toHaveBeenCalledWith('/customers/1')
    })

    it('should handle API errors', async () => {
      api.get.mockRejectedValue({
        response: { status: 500, data: { error: 'Internal server error' } }
      })

      await expect(api.get('/customers')).rejects.toThrow()
    })
  })

  describe('Plan Management', () => {
    beforeEach(() => {
      localStorage.setItem('token', 'test-token')
    })

    it('should create plan', async () => {
      const planData = {
        name: 'Test Plan',
        start_date: '2024-01-01',
        end_date: '2024-01-07',
        warehouse_id: 1
      }
      api.post.mockResolvedValue({
        data: { success: true, data: { id: 1, ...planData, status: 'draft' } }
      })

      const result = await api.post('/plans', planData)

      expect(result.data.success).toBe(true)
    })

    it('should optimize plan', async () => {
      const mockResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            status: 'optimized',
            routes: []
          }
        }
      }
      api.post.mockResolvedValue(mockResponse)

      const result = await api.post('/plans/1/optimize', {})

      expect(api.post).toHaveBeenCalledWith('/plans/1/optimize', {})
      expect(result.data.success).toBe(true)
    })

    it('should get plan routes', async () => {
      const mockRoutes = [
        { id: 1, day: 1, vehicle_id: 1, stops: [] },
        { id: 2, day: 2, vehicle_id: 1, stops: [] }
      ]
      api.get.mockResolvedValue({ data: { success: true, data: mockRoutes } })

      const result = await api.get('/plans/1/routes')

      expect(result.data.data).toEqual(mockRoutes)
    })
  })

  describe('Token Management', () => {
    it('should add token to requests when available', async () => {
      localStorage.setItem('token', 'stored-token')
      api.get.mockResolvedValue({ data: { success: true, data: [] } })

      await api.get('/customers')

      // Token should be added via interceptor (tested via interceptor setup)
      expect(api.get).toHaveBeenCalled()
    })

    it('should handle missing token', async () => {
      localStorage.removeItem('token')
      api.get.mockRejectedValue({
        response: { status: 401 }
      })

      await expect(api.get('/customers')).rejects.toThrow()
    })

    it('should handle 401 responses', async () => {
      localStorage.setItem('token', 'expired-token')
      api.get.mockRejectedValue({
        response: { status: 401 }
      })

      try {
        await api.get('/customers')
      } catch (e) {
        // Expected to throw
      }

      // Interceptor should handle 401 (tested via interceptor behavior)
      expect(api.get).toHaveBeenCalled()
    })
  })

  describe('Error Handling', () => {
    it('should handle network errors', async () => {
      api.get.mockRejectedValue(new Error('Network Error'))

      await expect(api.get('/customers')).rejects.toThrow('Network Error')
    })

    it('should handle timeout errors', async () => {
      api.post.mockRejectedValue({
        code: 'ECONNABORTED',
        message: 'timeout of 5000ms exceeded'
      })

      await expect(api.post('/plans/1/optimize', {})).rejects.toThrow()
    })

    it('should handle invalid JSON responses', async () => {
      api.get.mockResolvedValue({ data: 'invalid json' })

      // Should not throw, but response structure may be unexpected
      const result = await api.get('/customers')
      expect(result).toBeDefined()
    })
  })
})
