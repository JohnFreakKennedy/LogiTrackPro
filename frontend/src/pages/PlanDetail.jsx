import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { 
  Route, 
  ArrowLeft,
  Play,
  Loader2,
  Calendar,
  Truck,
  MapPin,
  DollarSign,
  Clock,
  Package,
  CheckCircle
} from 'lucide-react'
import api from '../api'

export default function PlanDetail() {
  const { id } = useParams()
  const [plan, setPlan] = useState(null)
  const [loading, setLoading] = useState(true)
  const [optimizing, setOptimizing] = useState(false)

  useEffect(() => {
    loadPlan()
  }, [id])

  const loadPlan = () => {
    setLoading(true)
    api.get(`/plans/${id}`)
      .then(res => setPlan(res.data.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const handleOptimize = async () => {
    setOptimizing(true)
    try {
      const res = await api.post(`/plans/${id}/optimize`)
      setPlan(res.data.data)
    } catch (err) {
      alert(err.response?.data?.error || 'Optimization failed')
    } finally {
      setOptimizing(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-500" />
      </div>
    )
  }

  if (!plan) {
    return (
      <div className="card text-center py-16">
        <p className="text-dark-400">Plan not found</p>
        <Link to="/plans" className="btn btn-primary mt-4">
          <ArrowLeft className="w-5 h-5" />
          Back to Plans
        </Link>
      </div>
    )
  }

  const getStatusBadge = (status) => {
    switch (status) {
      case 'optimized': return 'badge-success'
      case 'optimizing': return 'badge-warning'
      case 'executed': return 'badge-info'
      default: return 'badge-info'
    }
  }

  // Group routes by day
  const routesByDay = plan.routes?.reduce((acc, route) => {
    if (!acc[route.day]) acc[route.day] = []
    acc[route.day].push(route)
    return acc
  }, {}) || {}

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link to="/plans" className="p-2 rounded-lg bg-dark-800 hover:bg-dark-700 transition-colors">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <h1 className="text-3xl font-display font-bold">{plan.name}</h1>
            <span className={`badge ${getStatusBadge(plan.status)}`}>
              {plan.status}
            </span>
          </div>
          <p className="text-dark-400 mt-1 flex items-center gap-2">
            <Calendar className="w-4 h-4" />
            {new Date(plan.start_date).toLocaleDateString()} - {new Date(plan.end_date).toLocaleDateString()}
          </p>
        </div>
        {plan.status !== 'optimized' && (
          <button
            onClick={handleOptimize}
            disabled={optimizing}
            className="btn btn-accent"
          >
            {optimizing ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Optimizing...
              </>
            ) : (
              <>
                <Play className="w-5 h-5" />
                Run Optimization
              </>
            )}
          </button>
        )}
      </div>

      {/* Summary Cards */}
      {plan.status === 'optimized' && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="card"
          >
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-xl bg-primary-500/10 flex items-center justify-center">
                <MapPin className="w-6 h-6 text-primary-400" />
              </div>
              <div>
                <p className="text-dark-400 text-sm">Total Distance</p>
                <p className="text-2xl font-display font-bold">{plan.total_distance?.toFixed(1)} km</p>
              </div>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="card"
          >
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center">
                <DollarSign className="w-6 h-6 text-accent" />
              </div>
              <div>
                <p className="text-dark-400 text-sm">Total Cost</p>
                <p className="text-2xl font-display font-bold">${plan.total_cost?.toFixed(2)}</p>
              </div>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="card"
          >
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-xl bg-purple-500/10 flex items-center justify-center">
                <Route className="w-6 h-6 text-purple-400" />
              </div>
              <div>
                <p className="text-dark-400 text-sm">Total Routes</p>
                <p className="text-2xl font-display font-bold">{plan.routes?.length || 0}</p>
              </div>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="card"
          >
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center">
                <Package className="w-6 h-6 text-blue-400" />
              </div>
              <div>
                <p className="text-dark-400 text-sm">Total Deliveries</p>
                <p className="text-2xl font-display font-bold">
                  {plan.routes?.reduce((sum, r) => sum + (r.stops?.length || 0), 0) || 0}
                </p>
              </div>
            </div>
          </motion.div>
        </div>
      )}

      {/* Routes */}
      {plan.status === 'optimized' && plan.routes?.length > 0 ? (
        <div className="space-y-6">
          {Object.entries(routesByDay).map(([day, routes]) => (
            <motion.div
              key={day}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="card"
            >
              <div className="flex items-center gap-2 mb-6">
                <Calendar className="w-5 h-5 text-primary-400" />
                <h2 className="text-xl font-display font-semibold">
                  Day {day} - {routes[0]?.date ? new Date(routes[0].date).toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' }) : ''}
                </h2>
              </div>

              <div className="space-y-4">
                {routes.map((route, routeIndex) => (
                  <div key={route.id} className="p-4 bg-dark-800/50 rounded-xl">
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-lg bg-purple-500/10 flex items-center justify-center">
                          <Truck className="w-5 h-5 text-purple-400" />
                        </div>
                        <div>
                          <p className="font-semibold">{route.vehicle?.name || `Vehicle ${route.vehicle_id}`}</p>
                          <p className="text-sm text-dark-400">
                            {route.stops?.length || 0} stops
                          </p>
                        </div>
                      </div>
                      <div className="flex gap-4 text-sm">
                        <span className="flex items-center gap-1 text-dark-300">
                          <MapPin className="w-4 h-4" />
                          {route.total_distance?.toFixed(1)} km
                        </span>
                        <span className="flex items-center gap-1 text-dark-300">
                          <Package className="w-4 h-4" />
                          {route.total_load?.toFixed(0)} units
                        </span>
                        <span className="flex items-center gap-1 text-accent">
                          <DollarSign className="w-4 h-4" />
                          ${route.total_cost?.toFixed(2)}
                        </span>
                      </div>
                    </div>

                    {/* Route Stops */}
                    <div className="relative ml-5 pl-6 border-l-2 border-dark-600 space-y-4">
                      {/* Start: Warehouse */}
                      <div className="relative">
                        <div className="absolute -left-[29px] w-4 h-4 rounded-full bg-primary-500 border-2 border-dark-800"></div>
                        <div className="text-sm">
                          <span className="text-primary-400 font-medium">Start</span>
                          <span className="text-dark-400 ml-2">Warehouse</span>
                        </div>
                      </div>

                      {route.stops?.map((stop, stopIndex) => (
                        <div key={stop.id} className="relative">
                          <div className="absolute -left-[29px] w-4 h-4 rounded-full bg-dark-600 border-2 border-dark-800 flex items-center justify-center">
                            <span className="text-[10px] font-bold">{stopIndex + 1}</span>
                          </div>
                          <div className="flex items-center justify-between bg-dark-700/50 rounded-lg p-3">
                            <div>
                              <p className="font-medium">{stop.customer?.name || `Customer ${stop.customer_id}`}</p>
                              {stop.customer?.address && (
                                <p className="text-sm text-dark-400">{stop.customer.address}</p>
                              )}
                            </div>
                            <div className="flex items-center gap-4 text-sm">
                              <span className="flex items-center gap-1 text-dark-300">
                                <Clock className="w-4 h-4" />
                                {stop.arrival_time}
                              </span>
                              <span className="flex items-center gap-1 text-primary-400">
                                <Package className="w-4 h-4" />
                                {stop.quantity?.toFixed(0)} units
                              </span>
                            </div>
                          </div>
                        </div>
                      ))}

                      {/* End: Return to Warehouse */}
                      <div className="relative">
                        <div className="absolute -left-[29px] w-4 h-4 rounded-full bg-green-500 border-2 border-dark-800">
                          <CheckCircle className="w-3 h-3 text-dark-800 absolute -top-[1px] -left-[1px]" />
                        </div>
                        <div className="text-sm">
                          <span className="text-green-400 font-medium">End</span>
                          <span className="text-dark-400 ml-2">Return to Warehouse</span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </motion.div>
          ))}
        </div>
      ) : plan.status !== 'optimized' ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card text-center py-16"
        >
          <Route className="w-16 h-16 mx-auto text-dark-500 mb-4" />
          <h3 className="text-xl font-semibold mb-2">Ready to Optimize</h3>
          <p className="text-dark-400 mb-6">
            Click the "Run Optimization" button to generate optimized delivery routes
          </p>
          <button onClick={handleOptimize} disabled={optimizing} className="btn btn-accent">
            {optimizing ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Optimizing...
              </>
            ) : (
              <>
                <Play className="w-5 h-5" />
                Run Optimization
              </>
            )}
          </button>
        </motion.div>
      ) : (
        <div className="card text-center py-16">
          <p className="text-dark-400">No routes generated</p>
        </div>
      )}
    </div>
  )
}

