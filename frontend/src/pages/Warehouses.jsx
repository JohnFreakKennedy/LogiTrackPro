import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { 
  Warehouse, 
  Plus, 
  Pencil, 
  Trash2, 
  MapPin,
  Package,
  Loader2
} from 'lucide-react'
import Modal from '../components/Modal'
import api from '../api'

export default function Warehouses() {
  const [warehouses, setWarehouses] = useState([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    address: '',
    latitude: '',
    longitude: '',
    capacity: '',
    current_stock: '',
    holding_cost: '',
    replenishment_qty: '',
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = () => {
    setLoading(true)
    api.get('/warehouses')
      .then(res => setWarehouses(res.data.data || []))
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
        capacity: item.capacity.toString(),
        current_stock: item.current_stock.toString(),
        holding_cost: item.holding_cost.toString(),
        replenishment_qty: item.replenishment_qty.toString(),
      })
    } else {
      setEditingItem(null)
      setFormData({
        name: '',
        address: '',
        latitude: '',
        longitude: '',
        capacity: '1000',
        current_stock: '500',
        holding_cost: '0.5',
        replenishment_qty: '100',
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
      capacity: parseFloat(formData.capacity),
      current_stock: parseFloat(formData.current_stock),
      holding_cost: parseFloat(formData.holding_cost),
      replenishment_qty: parseFloat(formData.replenishment_qty),
    }
    
    try {
      if (editingItem) {
        await api.put(`/warehouses/${editingItem.id}`, data)
      } else {
        await api.post('/warehouses', data)
      }
      setModalOpen(false)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'An error occurred')
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this warehouse?')) return
    try {
      await api.delete(`/warehouses/${id}`)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to delete')
    }
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
          <h1 className="text-3xl font-display font-bold">Warehouses</h1>
          <p className="text-dark-400 mt-1">Manage your distribution centers</p>
        </div>
        <button onClick={() => openModal()} className="btn btn-primary">
          <Plus className="w-5 h-5" />
          Add Warehouse
        </button>
      </div>

      {warehouses.length === 0 ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card text-center py-16"
        >
          <Warehouse className="w-16 h-16 mx-auto text-dark-500 mb-4" />
          <h3 className="text-xl font-semibold mb-2">No warehouses yet</h3>
          <p className="text-dark-400 mb-6">Create your first warehouse to get started</p>
          <button onClick={() => openModal()} className="btn btn-primary">
            <Plus className="w-5 h-5" />
            Add Warehouse
          </button>
        </motion.div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {warehouses.map((warehouse, i) => (
            <motion.div
              key={warehouse.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              className="card card-hover"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center shadow-lg shadow-primary-500/25">
                    <Warehouse className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-lg">{warehouse.name}</h3>
                    <p className="text-sm text-dark-400 flex items-center gap-1">
                      <MapPin className="w-3 h-3" />
                      {warehouse.latitude.toFixed(4)}, {warehouse.longitude.toFixed(4)}
                    </p>
                  </div>
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => openModal(warehouse)}
                    className="p-2 rounded-lg hover:bg-dark-700 text-dark-400 hover:text-white transition-colors"
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDelete(warehouse.id)}
                    className="p-2 rounded-lg hover:bg-red-500/10 text-dark-400 hover:text-red-400 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              {warehouse.address && (
                <p className="text-sm text-dark-400 mb-4">{warehouse.address}</p>
              )}

              <div className="grid grid-cols-2 gap-3">
                <div className="p-3 bg-dark-800/50 rounded-lg">
                  <p className="text-xs text-dark-400">Capacity</p>
                  <p className="font-mono font-semibold">{warehouse.capacity.toLocaleString()}</p>
                </div>
                <div className="p-3 bg-dark-800/50 rounded-lg">
                  <p className="text-xs text-dark-400">Current Stock</p>
                  <p className="font-mono font-semibold">{warehouse.current_stock.toLocaleString()}</p>
                </div>
              </div>

              <div className="mt-3 flex items-center justify-between text-sm">
                <span className="text-dark-400">Utilization</span>
                <span className="font-mono">
                  {((warehouse.current_stock / warehouse.capacity) * 100).toFixed(1)}%
                </span>
              </div>
              <div className="mt-2 h-2 bg-dark-700 rounded-full overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-primary-500 to-primary-400 rounded-full transition-all"
                  style={{ width: `${Math.min((warehouse.current_stock / warehouse.capacity) * 100, 100)}%` }}
                />
              </div>
            </motion.div>
          ))}
        </div>
      )}

      <Modal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editingItem ? 'Edit Warehouse' : 'New Warehouse'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Main Distribution Center"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Address</label>
            <input
              type="text"
              value={formData.address}
              onChange={(e) => setFormData({ ...formData, address: e.target.value })}
              placeholder="123 Warehouse St, City"
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
                placeholder="40.7128"
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
                placeholder="-74.0060"
                required
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Capacity</label>
              <input
                type="number"
                step="any"
                value={formData.capacity}
                onChange={(e) => setFormData({ ...formData, capacity: e.target.value })}
                placeholder="1000"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Current Stock</label>
              <input
                type="number"
                step="any"
                value={formData.current_stock}
                onChange={(e) => setFormData({ ...formData, current_stock: e.target.value })}
                placeholder="500"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Holding Cost ($/unit)</label>
              <input
                type="number"
                step="any"
                value={formData.holding_cost}
                onChange={(e) => setFormData({ ...formData, holding_cost: e.target.value })}
                placeholder="0.5"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Replenishment Qty</label>
              <input
                type="number"
                step="any"
                value={formData.replenishment_qty}
                onChange={(e) => setFormData({ ...formData, replenishment_qty: e.target.value })}
                placeholder="100"
              />
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <button type="button" onClick={() => setModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1">
              {editingItem ? 'Save Changes' : 'Create Warehouse'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

