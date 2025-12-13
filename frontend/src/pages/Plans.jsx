import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { 
  Route, 
  Plus, 
  Trash2,
  Calendar,
  Eye,
  Play,
  Loader2,
  DollarSign,
  MapPin
} from 'lucide-react'
import Modal from '../components/Modal'
import api from '../api'

export default function Plans() {
  const [plans, setPlans] = useState([])
  const [warehouses, setWarehouses] = useState([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    start_date: '',
    end_date: '',
    warehouse_id: '',
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = () => {
    setLoading(true)
    Promise.all([
      api.get('/plans'),
      api.get('/warehouses')
    ])
      .then(([plansRes, warehousesRes]) => {
        setPlans(plansRes.data.data || [])
        setWarehouses(warehousesRes.data.data || [])
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const openModal = () => {
    const today = new Date().toISOString().split('T')[0]
    const nextWeek = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
    setFormData({
      name: '',
      start_date: today,
      end_date: nextWeek,
      warehouse_id: warehouses[0]?.id?.toString() || '',
    })
    setModalOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    const warehouseId = parseInt(formData.warehouse_id)
    if (isNaN(warehouseId)) {
      alert('Please select a warehouse')
      return
    }

    const data = {
      name: formData.name,
      start_date: formData.start_date,
      end_date: formData.end_date,
      warehouse_id: warehouseId,
    }
    
    try {
      await api.post('/plans', data)
      setModalOpen(false)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'An error occurred')
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this plan?')) return
    try {
      await api.delete(`/plans/${id}`)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to delete')
    }
  }

  const getStatusBadge = (status) => {
    switch (status) {
      case 'optimized': return 'badge-success'
      case 'optimizing': return 'badge-warning'
      case 'executed': return 'badge-info'
      default: return 'badge-info'
    }
  }

  const getWarehouseName = (id) => {
    return warehouses.find(w => w.id === id)?.name || 'Unknown'
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-500" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-display font-bold">Delivery Plans</h1>
          <p className="text-dark-400 mt-1">Create and optimize delivery routes</p>
        </div>
        <button onClick={openModal} className="btn btn-primary">
          <Plus className="w-5 h-5" />
          New Plan
        </button>
      </div>

      {plans.length === 0 ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card text-center py-16"
        >
          <Route className="w-16 h-16 mx-auto text-dark-500 mb-4" />
          <h3 className="text-xl font-semibold mb-2">No delivery plans yet</h3>
          <p className="text-dark-400 mb-6">Create a plan to start optimizing your routes</p>
          <button onClick={openModal} className="btn btn-primary">
            <Plus className="w-5 h-5" />
            Create First Plan
          </button>
        </motion.div>
      ) : (
        <div className="grid gap-4">
          {plans.map((plan, i) => (
            <motion.div
              key={plan.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              className="card card-hover"
            >
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div className="flex items-center gap-4">
                  <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-accent to-accent-dark flex items-center justify-center shadow-lg shadow-accent/25">
                    <Route className="w-7 h-7 text-white" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <h3 className="font-semibold text-lg">{plan.name}</h3>
                      <span className={`badge ${getStatusBadge(plan.status)}`}>
                        {plan.status}
                      </span>
                    </div>
                    <div className="flex items-center gap-4 mt-1 text-sm text-dark-400">
                      <span className="flex items-center gap-1">
                        <Calendar className="w-4 h-4" />
                        {new Date(plan.start_date).toLocaleDateString()} - {new Date(plan.end_date).toLocaleDateString()}
                      </span>
                      <span>Warehouse: {getWarehouseName(plan.warehouse_id)}</span>
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-6">
                  {plan.status === 'optimized' && (
                    <div className="flex gap-6">
                      <div className="text-right">
                        <p className="text-xs text-dark-400 flex items-center gap-1">
                          <MapPin className="w-3 h-3" /> Distance
                        </p>
                        <p className="font-mono font-semibold">{plan.total_distance?.toFixed(1) || 0} km</p>
                      </div>
                      <div className="text-right">
                        <p className="text-xs text-dark-400 flex items-center gap-1">
                          <DollarSign className="w-3 h-3" /> Cost
                        </p>
                        <p className="font-mono font-semibold">${plan.total_cost?.toFixed(2) || '0.00'}</p>
                      </div>
                    </div>
                  )}
                  
                  <div className="flex gap-2">
                    <Link
                      to={`/plans/${plan.id}`}
                      className="btn btn-secondary"
                    >
                      <Eye className="w-4 h-4" />
                      View
                    </Link>
                    <button
                      onClick={() => handleDelete(plan.id)}
                      className="p-2 rounded-lg hover:bg-red-500/10 text-dark-400 hover:text-red-400 transition-colors"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      )}

      <Modal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title="Create New Plan"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Plan Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Weekly Delivery Plan"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Warehouse</label>
            <select
              value={formData.warehouse_id}
              onChange={(e) => setFormData({ ...formData, warehouse_id: e.target.value })}
              required
            >
              <option value="">Select a warehouse</option>
              {warehouses.map(w => (
                <option key={w.id} value={w.id}>{w.name}</option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Start Date</label>
              <input
                type="date"
                value={formData.start_date}
                onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">End Date</label>
              <input
                type="date"
                value={formData.end_date}
                onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                required
              />
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <button type="button" onClick={() => setModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1">
              Create Plan
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

