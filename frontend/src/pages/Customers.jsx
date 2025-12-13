import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { 
  Users, 
  Plus, 
  Pencil, 
  Trash2, 
  MapPin,
  TrendingUp,
  Loader2
} from 'lucide-react'
import Modal from '../components/Modal'
import api from '../api'

export default function Customers() {
  const [customers, setCustomers] = useState([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    address: '',
    latitude: '',
    longitude: '',
    demand_rate: '',
    max_inventory: '',
    current_inventory: '',
    min_inventory: '',
    holding_cost: '',
    priority: '1',
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = () => {
    setLoading(true)
    api.get('/customers')
      .then(res => setCustomers(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const openModal = (item = null) => {
    if (item) {
      setEditingItem(item)
      setFormData({
        name: item.name,
        address: item.address || '',
        latitude: item.latitude.toString(),
        longitude: item.longitude.toString(),
        demand_rate: item.demand_rate.toString(),
        max_inventory: item.max_inventory.toString(),
        current_inventory: item.current_inventory.toString(),
        min_inventory: item.min_inventory.toString(),
        holding_cost: item.holding_cost.toString(),
        priority: item.priority.toString(),
      })
    } else {
      setEditingItem(null)
      setFormData({
        name: '',
        address: '',
        latitude: '',
        longitude: '',
        demand_rate: '10',
        max_inventory: '100',
        current_inventory: '50',
        min_inventory: '10',
        holding_cost: '0.3',
        priority: '1',
      })
    }
    setModalOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const data = {
      name: formData.name,
      address: formData.address,
      latitude: parseFloat(formData.latitude),
      longitude: parseFloat(formData.longitude),
      demand_rate: parseFloat(formData.demand_rate),
      max_inventory: parseFloat(formData.max_inventory),
      current_inventory: parseFloat(formData.current_inventory),
      min_inventory: parseFloat(formData.min_inventory),
      holding_cost: parseFloat(formData.holding_cost),
      priority: parseInt(formData.priority),
    }
    
    try {
      if (editingItem) {
        await api.put(`/customers/${editingItem.id}`, data)
      } else {
        await api.post('/customers', data)
      }
      setModalOpen(false)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'An error occurred')
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this customer?')) return
    try {
      await api.delete(`/customers/${id}`)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to delete')
    }
  }

  const getPriorityBadge = (priority) => {
    if (priority >= 3) return 'badge-danger'
    if (priority === 2) return 'badge-warning'
    return 'badge-info'
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
          <h1 className="text-3xl font-display font-bold">Customers</h1>
          <p className="text-dark-400 mt-1">Manage customer locations and demands</p>
        </div>
        <button onClick={() => openModal()} className="btn btn-primary">
          <Plus className="w-5 h-5" />
          Add Customer
        </button>
      </div>

      {customers.length === 0 ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card text-center py-16"
        >
          <Users className="w-16 h-16 mx-auto text-dark-500 mb-4" />
          <h3 className="text-xl font-semibold mb-2">No customers yet</h3>
          <p className="text-dark-400 mb-6">Add your first customer to start planning deliveries</p>
          <button onClick={() => openModal()} className="btn btn-primary">
            <Plus className="w-5 h-5" />
            Add Customer
          </button>
        </motion.div>
      ) : (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Customer</th>
                <th>Location</th>
                <th>Demand Rate</th>
                <th>Inventory</th>
                <th>Priority</th>
                <th className="text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {customers.map((customer, i) => (
                <motion.tr
                  key={customer.id}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: i * 0.03 }}
                >
                  <td>
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-blue-500/10 flex items-center justify-center">
                        <Users className="w-5 h-5 text-blue-400" />
                      </div>
                      <div>
                        <p className="font-medium">{customer.name}</p>
                        {customer.address && (
                          <p className="text-sm text-dark-400">{customer.address}</p>
                        )}
                      </div>
                    </div>
                  </td>
                  <td>
                    <div className="flex items-center gap-1 text-dark-300">
                      <MapPin className="w-4 h-4" />
                      <span className="font-mono text-sm">
                        {customer.latitude.toFixed(4)}, {customer.longitude.toFixed(4)}
                      </span>
                    </div>
                  </td>
                  <td>
                    <div className="flex items-center gap-1">
                      <TrendingUp className="w-4 h-4 text-accent" />
                      <span className="font-mono">{customer.demand_rate}/day</span>
                    </div>
                  </td>
                  <td>
                    <div>
                      <p className="font-mono">
                        {customer.current_inventory} / {customer.max_inventory}
                      </p>
                      <div className="w-20 h-1.5 bg-dark-700 rounded-full mt-1">
                        <div
                          className="h-full bg-primary-500 rounded-full"
                          style={{ width: `${(customer.current_inventory / customer.max_inventory) * 100}%` }}
                        />
                      </div>
                    </div>
                  </td>
                  <td>
                    <span className={`badge ${getPriorityBadge(customer.priority)}`}>
                      P{customer.priority}
                    </span>
                  </td>
                  <td>
                    <div className="flex justify-end gap-1">
                      <button
                        onClick={() => openModal(customer)}
                        className="p-2 rounded-lg hover:bg-dark-700 text-dark-400 hover:text-white"
                      >
                        <Pencil className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleDelete(customer.id)}
                        className="p-2 rounded-lg hover:bg-red-500/10 text-dark-400 hover:text-red-400"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Modal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editingItem ? 'Edit Customer' : 'New Customer'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Acme Corp"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Address</label>
            <input
              type="text"
              value={formData.address}
              onChange={(e) => setFormData({ ...formData, address: e.target.value })}
              placeholder="456 Customer Ave, City"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Latitude</label>
              <input
                type="number"
                step="any"
                value={formData.latitude}
                onChange={(e) => setFormData({ ...formData, latitude: e.target.value })}
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Longitude</label>
              <input
                type="number"
                step="any"
                value={formData.longitude}
                onChange={(e) => setFormData({ ...formData, longitude: e.target.value })}
                required
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Demand Rate (per day)</label>
              <input
                type="number"
                step="any"
                value={formData.demand_rate}
                onChange={(e) => setFormData({ ...formData, demand_rate: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Priority (1-5)</label>
              <select
                value={formData.priority}
                onChange={(e) => setFormData({ ...formData, priority: e.target.value })}
              >
                <option value="1">1 - Low</option>
                <option value="2">2 - Medium</option>
                <option value="3">3 - High</option>
                <option value="4">4 - Critical</option>
                <option value="5">5 - Emergency</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Min Inventory</label>
              <input
                type="number"
                step="any"
                value={formData.min_inventory}
                onChange={(e) => setFormData({ ...formData, min_inventory: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Current</label>
              <input
                type="number"
                step="any"
                value={formData.current_inventory}
                onChange={(e) => setFormData({ ...formData, current_inventory: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Max Inventory</label>
              <input
                type="number"
                step="any"
                value={formData.max_inventory}
                onChange={(e) => setFormData({ ...formData, max_inventory: e.target.value })}
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Holding Cost ($/unit)</label>
            <input
              type="number"
              step="any"
              value={formData.holding_cost}
              onChange={(e) => setFormData({ ...formData, holding_cost: e.target.value })}
            />
          </div>

          <div className="flex gap-3 pt-4">
            <button type="button" onClick={() => setModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1">
              {editingItem ? 'Save Changes' : 'Create Customer'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

