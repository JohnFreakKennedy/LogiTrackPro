import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { 
  Warehouse, 
  Users, 
  Truck, 
  Route,
  TrendingUp,
  Package,
  MapPin,
  DollarSign,
  ArrowRight,
  Calendar,
  Loader2
} from 'lucide-react'
import api from '../api'

export default function Dashboard() {
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/analytics/dashboard')
      .then(res => setData(res.data.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-500" />
      </div>
    )
  }

  const stats = [
    { 
      label: 'Warehouses', 
      value: data?.total_warehouses || 0, 
      icon: Warehouse, 
      color: 'primary',
      link: '/warehouses'
    },
    { 
      label: 'Customers', 
      value: data?.total_customers || 0, 
      icon: Users, 
      color: 'blue',
      link: '/customers'
    },
    { 
      label: 'Vehicles', 
      value: data?.total_vehicles || 0, 
      icon: Truck, 
      color: 'purple',
      link: '/vehicles'
    },
    { 
      label: 'Active Plans', 
      value: data?.active_plans || 0, 
      icon: Route, 
      color: 'accent',
      link: '/plans'
    },
  ]

  const metrics = [
    { 
      label: 'Total Deliveries', 
      value: data?.total_deliveries || 0,
      icon: Package,
    },
    { 
      label: 'Distance (km)', 
      value: (data?.total_distance_km || 0).toFixed(1),
      icon: MapPin,
    },
    { 
      label: 'Total Cost', 
      value: `$${(data?.total_cost || 0).toFixed(2)}`,
      icon: DollarSign,
    },
  ]

  const getColorClasses = (color) => {
    const colors = {
      primary: 'from-primary-500 to-primary-600 shadow-primary-500/25',
      blue: 'from-blue-500 to-blue-600 shadow-blue-500/25',
      purple: 'from-purple-500 to-purple-600 shadow-purple-500/25',
      accent: 'from-accent to-accent-dark shadow-accent/25',
    }
    return colors[color] || colors.primary
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-display font-bold">Dashboard</h1>
        <p className="text-dark-400 mt-1">Overview of your logistics operations</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <motion.div
            key={stat.label}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.1 }}
          >
            <Link
              to={stat.link}
              className="card card-hover block group"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-dark-400 text-sm">{stat.label}</p>
                  <p className="text-3xl font-display font-bold mt-1">{stat.value}</p>
                </div>
                <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${getColorClasses(stat.color)} flex items-center justify-center shadow-lg`}>
                  <stat.icon className="w-6 h-6 text-white" />
                </div>
              </div>
              <div className="flex items-center gap-1 mt-4 text-sm text-dark-400 group-hover:text-primary-400 transition-colors">
                <span>View all</span>
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </div>
            </Link>
          </motion.div>
        ))}
      </div>

      {/* Metrics and Recent Plans */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Performance Metrics */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="card"
        >
          <div className="flex items-center gap-2 mb-6">
            <TrendingUp className="w-5 h-5 text-primary-400" />
            <h2 className="text-lg font-semibold">Performance Metrics</h2>
          </div>
          <div className="space-y-4">
            {metrics.map((metric) => (
              <div key={metric.label} className="flex items-center justify-between p-4 bg-dark-800/50 rounded-xl">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-dark-700 flex items-center justify-center">
                    <metric.icon className="w-5 h-5 text-dark-300" />
                  </div>
                  <span className="text-dark-300">{metric.label}</span>
                </div>
                <span className="font-mono font-semibold text-lg">{metric.value}</span>
              </div>
            ))}
          </div>
        </motion.div>

        {/* Recent Plans */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="lg:col-span-2 card"
        >
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-2">
              <Calendar className="w-5 h-5 text-primary-400" />
              <h2 className="text-lg font-semibold">Recent Plans</h2>
            </div>
            <Link to="/plans" className="text-sm text-primary-400 hover:text-primary-300 flex items-center gap-1">
              View all <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
          
          {data?.recent_plans?.length > 0 ? (
            <div className="space-y-3">
              {data.recent_plans.map((plan) => (
                <Link
                  key={plan.id}
                  to={`/plans/${plan.id}`}
                  className="flex items-center justify-between p-4 bg-dark-800/50 rounded-xl hover:bg-dark-800 transition-colors group"
                >
                  <div>
                    <p className="font-medium group-hover:text-primary-400 transition-colors">
                      {plan.name}
                    </p>
                    <p className="text-sm text-dark-400">
                      {new Date(plan.start_date).toLocaleDateString()} - {new Date(plan.end_date).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <p className="text-sm text-dark-400">Cost</p>
                      <p className="font-mono font-semibold">${plan.total_cost?.toFixed(2) || '0.00'}</p>
                    </div>
                    <span className={`badge ${
                      plan.status === 'optimized' ? 'badge-success' :
                      plan.status === 'optimizing' ? 'badge-warning' :
                      plan.status === 'executed' ? 'badge-info' : 'badge-info'
                    }`}>
                      {plan.status}
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="text-center py-12 text-dark-400">
              <Route className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>No plans created yet</p>
              <Link to="/plans" className="text-primary-400 hover:text-primary-300 mt-2 inline-block">
                Create your first plan
              </Link>
            </div>
          )}
        </motion.div>
      </div>

      {/* Quick Actions */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.6 }}
        className="card"
      >
        <h2 className="text-lg font-semibold mb-4">Quick Actions</h2>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <Link to="/warehouses" className="btn btn-secondary">
            <Warehouse className="w-4 h-4" />
            Add Warehouse
          </Link>
          <Link to="/customers" className="btn btn-secondary">
            <Users className="w-4 h-4" />
            Add Customer
          </Link>
          <Link to="/vehicles" className="btn btn-secondary">
            <Truck className="w-4 h-4" />
            Add Vehicle
          </Link>
          <Link to="/plans" className="btn btn-accent">
            <Route className="w-4 h-4" />
            New Plan
          </Link>
        </div>
      </motion.div>
    </div>
  )
}

