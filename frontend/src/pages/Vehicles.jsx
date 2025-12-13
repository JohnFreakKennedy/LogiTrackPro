import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { 
  Truck, 
  Plus, 
  Pencil, 
  Trash2,
  Gauge,
  DollarSign,
  MapPin,
  CheckCircle,
  XCircle,
  Loader2
} from 'lucide-react'
import Modal from '../components/Modal'
import api from '../api'

export default function Vehicles() {
  const [vehicles, setVehicles] = useState([])
  const [warehouses, setWarehouses] = useState([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    capacity: '',
    cost_per_km: '',
    fixed_cost: '',
    max_distance: '',
    available: true,
    warehouse_id: '',
  })

  useEffect(() => {
    Promise.all([
      api.get('/vehicles'),
      api.get('/warehouses')
    ])
      .then(([vehiclesRes, warehousesRes]) => {
        setVehicles(vehiclesRes.data.data || [])
        setWarehouses(warehousesRes.data.data || [])
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  const loadData = () => {
    api.get('/vehicles')
      .then(res => setVehicles(res.data.data || []))
      .catch(console.error)
  }

  const openModal = (item = null) => {
    if (item) {
      setEditingItem(item)
      setFormData({
        name: item.name,
        capacity: item.capacity.toString(),
        cost_per_km: item.cost_per_km.toString(),
        fixed_cost: item.fixed_cost.toString(),
        max_distance: item.max_distance.toString(),
        available: item.available,
        warehouse_id: item.warehouse_id?.toString() || '',
      })
    } else {
      setEditingItem(null)
      setFormData({
        name: '',
        capacity: '1000',
        cost_per_km: '0.5',
        fixed_cost: '50',
        max_distance: '300',
        available: true,
        warehouse_id: warehouses[0]?.id?.toString() || '',
      })
    }
    setModalOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    const warehouseId = parseInt(formData.warehouse_id)
    
    const data = {
      name: formData.name,
      capacity: parseFloat(formData.capacity),
      cost_per_km: parseFloat(formData.cost_per_km),
      fixed_cost: parseFloat(formData.fixed_cost),
      max_distance: parseFloat(formData.max_distance),
      available: formData.available,
      warehouse_id: isNaN(warehouseId) ? 0 : warehouseId,
    }
    
    try {
      if (editingItem) {
        await api.put(`/vehicles/${editingItem.id}`, data)
      } else {
        await api.post('/vehicles', data)
      }
      setModalOpen(false)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'An error occurred')
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this vehicle?')) return
    try {
      await api.delete(`/vehicles/${id}`)
      loadData()
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to delete')
    }
  }

  const getWarehouseName = (id) => {
    return warehouses.find(w => w.id === id)?.name || 'Unassigned'
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
          <h1 className="text-3xl font-display font-bold">Vehicles</h1>
          <p className="text-dark-400 mt-1">Manage your delivery fleet</p>
        </div>
        <button onClick={() => openModal()} className="btn btn-primary">
          <Plus className="w-5 h-5" />
          Add Vehicle
        </button>
      </div>

      {vehicles.length === 0 ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card text-center py-16"
        >
          <Truck className="w-16 h-16 mx-auto text-dark-500 mb-4" />
          <h3 className="text-xl font-semibold mb-2">No vehicles yet</h3>
          <p className="text-dark-400 mb-6">Add vehicles to enable route optimization</p>
          <button onClick={() => openModal()} className="btn btn-primary">
            <Plus className="w-5 h-5" />
            Add Vehicle
          </button>
        </motion.div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {vehicles.map((vehicle, i) => (
            <motion.div
              key={vehicle.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              className="card card-hover"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className={`w-12 h-12 rounded-xl flex items-center justify-center shadow-lg ${
                    vehicle.available 
                      ? 'bg-gradient-to-br from-purple-500 to-purple-600 shadow-purple-500/25'
                      : 'bg-gradient-to-br from-dark-600 to-dark-700'
                  }`}>
                    <Truck className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-lg">{vehicle.name}</h3>
                    <div className="flex items-center gap-1 text-sm">
                      {vehicle.available ? (
                        <><CheckCircle className="w-3 h-3 text-primary-400" /><span className="text-primary-400">Available</span></>
                      ) : (
                        <><XCircle className="w-3 h-3 text-dark-400" /><span className="text-dark-400">Unavailable</span></>
                      )}
                    </div>
                  </div>
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => openModal(vehicle)}
                    className="p-2 rounded-lg hover:bg-dark-700 text-dark-400 hover:text-white transition-colors"
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDelete(vehicle.id)}
                    className="p-2 rounded-lg hover:bg-red-500/10 text-dark-400 hover:text-red-400 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3 mb-4">
                <div className="p-3 bg-dark-800/50 rounded-lg">
                  <div className="flex items-center gap-1 text-dark-400 text-xs mb-1">
                    <Gauge className="w-3 h-3" />
                    Capacity
                  </div>
                  <p className="font-mono font-semibold">{vehicle.capacity.toLocaleString()}</p>
                </div>
                <div className="p-3 bg-dark-800/50 rounded-lg">
                  <div className="flex items-center gap-1 text-dark-400 text-xs mb-1">
                    <MapPin className="w-3 h-3" />
                    Max Distance
                  </div>
                  <p className="font-mono font-semibold">{vehicle.max_distance} km</p>
                </div>
              </div>

              <div className="flex items-center justify-between p-3 bg-dark-800/50 rounded-lg mb-3">
                <div className="flex items-center gap-1 text-dark-400 text-sm">
                  <DollarSign className="w-4 h-4" />
                  Cost
                </div>
                <div className="text-right">
                  <span className="font-mono">${vehicle.cost_per_km}/km</span>
                  <span className="text-dark-400 mx-1">+</span>
                  <span className="font-mono">${vehicle.fixed_cost} fixed</span>
                </div>
              </div>

              <div className="text-sm text-dark-400">
                Warehouse: <span className="text-dark-200">{getWarehouseName(vehicle.warehouse_id)}</span>
              </div>
            </motion.div>
          ))}
        </div>
      )}

      <Modal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editingItem ? 'Edit Vehicle' : 'New Vehicle'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Truck 001"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2 text-dark-300">Warehouse</label>
            <select
              value={formData.warehouse_id}
              onChange={(e) => setFormData({ ...formData, warehouse_id: e.target.value })}
            >
              <option value="">No warehouse assigned</option>
              {warehouses.map(w => (
                <option key={w.id} value={w.id}>{w.name}</option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Capacity (units)</label>
              <input
                type="number"
                step="any"
                value={formData.capacity}
                onChange={(e) => setFormData({ ...formData, capacity: e.target.value })}
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Max Distance (km)</label>
              <input
                type="number"
                step="any"
                value={formData.max_distance}
                onChange={(e) => setFormData({ ...formData, max_distance: e.target.value })}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Cost per km ($)</label>
              <input
                type="number"
                step="any"
                value={formData.cost_per_km}
                onChange={(e) => setFormData({ ...formData, cost_per_km: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2 text-dark-300">Fixed Cost ($)</label>
              <input
                type="number"
                step="any"
                value={formData.fixed_cost}
                onChange={(e) => setFormData({ ...formData, fixed_cost: e.target.value })}
              />
            </div>
          </div>

          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="available"
              checked={formData.available}
              onChange={(e) => setFormData({ ...formData, available: e.target.checked })}
              className="w-5 h-5 rounded border-dark-600 bg-dark-700 text-primary-500 focus:ring-primary-500"
            />
            <label htmlFor="available" className="text-sm text-dark-300">Vehicle is available for routing</label>
          </div>

          <div className="flex gap-3 pt-4">
            <button type="button" onClick={() => setModalOpen(false)} className="btn btn-secondary flex-1">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary flex-1">
              {editingItem ? 'Save Changes' : 'Create Vehicle'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

