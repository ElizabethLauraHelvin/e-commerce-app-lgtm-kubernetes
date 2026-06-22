<template>
  <div id="app">
    <!-- Login/Register Page -->
    <div v-if="!user" class="auth-page">
      <div class="auth-container">
        <div class="auth-brand">
          <h1>🛒 ShopHub</h1>
          <p>Your one-stop online shop</p>
        </div>
        <div class="auth-form">
          <div class="auth-tabs">
            <button :class="{ active: authMode === 'login' }" @click="authMode = 'login'">Login</button>
            <button :class="{ active: authMode === 'register' }" @click="authMode = 'register'">Register</button>
          </div>
          <form @submit.prevent="handleAuth">
            <input v-if="authMode === 'register'" v-model="authForm.name" type="text" placeholder="Full Name" required />
            <input v-model="authForm.email" type="email" placeholder="Email" required />
            <input v-model="authForm.password" type="password" placeholder="Password" required />
            <button type="submit" class="btn-primary" :disabled="authLoading">
              {{ authLoading ? 'Please wait...' : (authMode === 'login' ? 'Sign In' : 'Create Account') }}
            </button>
          </form>
          <p v-if="authError" class="auth-error">{{ authError }}</p>
          <p class="auth-hint">Demo: demo@shop.com / demo123</p>
        </div>
      </div>
    </div>

    <!-- Main App -->
    <div v-else class="shop">
      <!-- Top Bar -->
      <header class="topbar">
        <div class="topbar-left">
          <h1 @click="currentPage = 'products'">🛒 ShopHub</h1>
        </div>
        <div class="topbar-center">
          <input v-model="searchQuery" type="text" placeholder="Search products..." class="search-input" />
        </div>
        <div class="topbar-right">
          <button class="cart-btn" @click="currentPage = 'cart'">
            🛒 <span v-if="cart.length" class="cart-badge">{{ cartItemCount }}</span>
          </button>
          <button class="orders-btn" @click="currentPage = 'orders'">📦 Orders</button>
          <div class="user-menu">
            <span>Hi, {{ user.name }}</span>
            <button @click="logout" class="logout-btn">Logout</button>
          </div>
        </div>
      </header>

      <!-- Products Page -->
      <main v-if="currentPage === 'products'" class="main-content">
        <div class="categories">
          <button :class="{ active: selectedCategory === 'All' }" @click="selectedCategory = 'All'">All</button>
          <button v-for="cat in categories" :key="cat" :class="{ active: selectedCategory === cat }" @click="selectedCategory = cat">{{ cat }}</button>
        </div>
        <div class="products-grid">
          <div v-for="product in filteredProducts" :key="product.id" class="product-card">
            <div class="product-image">
              <img :src="product.image" :alt="product.name" @error="handleImageError" />
            </div>
            <div class="product-info">
              <span class="product-category">{{ product.category }}</span>
              <h3>{{ product.name }}</h3>
              <div class="product-footer">
                <span class="product-price">${{ product.price.toFixed(2) }}</span>
                <button class="add-cart-btn" @click="addToCart(product)" :disabled="product.stock === 0">
                  {{ product.stock === 0 ? 'Out of Stock' : '+ Cart' }}
                </button>
              </div>
              <span class="product-stock">{{ product.stock }} in stock</span>
            </div>
          </div>
        </div>
      </main>

      <!-- Cart Page -->
      <main v-if="currentPage === 'cart'" class="main-content">
        <h2>Shopping Cart</h2>
        <div v-if="cart.length === 0" class="empty-state">
          <p>Your cart is empty</p>
          <button class="btn-primary" @click="currentPage = 'products'">Continue Shopping</button>
        </div>
        <div v-else class="cart-layout">
          <div class="cart-items">
            <div v-for="item in cart" :key="item.product.id" class="cart-item">
              <img :src="item.product.image" :alt="item.product.name" @error="handleImageError" />
              <div class="cart-item-info">
                <h4>{{ item.product.name }}</h4>
                <p class="cart-item-price">${{ item.product.price.toFixed(2) }}</p>
              </div>
              <div class="cart-item-qty">
                <button @click="updateQty(item, -1)">-</button>
                <span>{{ item.qty }}</span>
                <button @click="updateQty(item, 1)">+</button>
              </div>
              <span class="cart-item-total">${{ (item.product.price * item.qty).toFixed(2) }}</span>
              <button class="remove-btn" @click="removeFromCart(item)">✕</button>
            </div>
          </div>
          <div class="cart-summary">
            <h3>Order Summary</h3>
            <div class="summary-row"><span>Subtotal</span><span>${{ cartTotal.toFixed(2) }}</span></div>
            <div class="summary-row"><span>Shipping</span><span>Free</span></div>
            <div class="summary-row total"><span>Total</span><span>${{ cartTotal.toFixed(2) }}</span></div>
            <div class="payment-method">
              <label>Payment Method</label>
              <select v-model="paymentMethod">
                <option value="credit_card">Credit Card</option>
                <option value="bank_transfer">Bank Transfer</option>
                <option value="e_wallet">E-Wallet</option>
              </select>
            </div>
            <button class="btn-checkout" @click="checkout" :disabled="checkoutLoading">
              {{ checkoutLoading ? 'Processing...' : 'Checkout' }}
            </button>
          </div>
        </div>
      </main>

      <!-- Orders Page -->
      <main v-if="currentPage === 'orders'" class="main-content">
        <h2>My Orders</h2>
        <div v-if="orders.length === 0" class="empty-state">
          <p>No orders yet</p>
          <button class="btn-primary" @click="currentPage = 'products'">Start Shopping</button>
        </div>
        <div v-else class="orders-list">
          <div v-for="order in orders" :key="order.id" class="order-card">
            <div class="order-header">
              <span class="order-id">{{ order.id }}</span>
              <span :class="['order-status', order.status]">{{ order.status }}</span>
            </div>
            <div class="order-items">
              <span v-for="item in order.items" :key="item.product_id" class="order-item-tag">
                {{ item.name }} x{{ item.quantity }}
              </span>
            </div>
            <div class="order-footer">
              <span class="order-total">${{ order.total.toFixed(2) }}</span>
              <span class="order-date">{{ formatDate(order.created_at) }}</span>
            </div>
          </div>
        </div>
      </main>

      <!-- Notification -->
      <div v-if="notification" :class="['toast', notification.type]">
        {{ notification.message }}
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { initFaro, getFaro } from './faro'

initFaro()

const API_URL = (window.__CONFIG__ && window.__CONFIG__.API_URL) || import.meta.env.VITE_API_URL || 'http://4.144.133.123:8080/api'

export default {
  name: 'App',
  data() {
    return {
      currentPage: 'products',
      user: null,
      authMode: 'login',
      authForm: { name: '', email: '', password: '' },
      authLoading: false,
      authError: '',
      products: [],
      cart: [],
      orders: [],
      searchQuery: '',
      selectedCategory: 'All',
      paymentMethod: 'credit_card',
      checkoutLoading: false,
      notification: null,
    }
  },
  computed: {
    categories() {
      return [...new Set(this.products.map(p => p.category))]
    },
    filteredProducts() {
      return this.products.filter(p => {
        const matchCategory = this.selectedCategory === 'All' || p.category === this.selectedCategory
        const matchSearch = p.name.toLowerCase().includes(this.searchQuery.toLowerCase())
        return matchCategory && matchSearch
      })
    },
    cartItemCount() {
      return this.cart.reduce((sum, item) => sum + item.qty, 0)
    },
    cartTotal() {
      return this.cart.reduce((sum, item) => sum + item.product.price * item.qty, 0)
    },
  },
  mounted() {
    const saved = localStorage.getItem('shopUser')
    if (saved) {
      this.user = JSON.parse(saved)
      this.loadProducts()
      this.loadOrders()
    }
  },
  methods: {
    async handleAuth() {
      this.authLoading = true
      this.authError = ''
      try {
        const endpoint = this.authMode === 'login' ? '/auth/login' : '/auth/register'
        const res = await axios.post(`${API_URL}${endpoint}`, this.authForm)
        this.user = res.data.user
        localStorage.setItem('shopUser', JSON.stringify(this.user))
        this.loadProducts()
        this.loadOrders()
        getFaro()?.api.pushEvent('user_auth', { type: this.authMode, email: this.user.email })
      } catch (err) {
        this.authError = err.response?.data?.error || 'Authentication failed'
      } finally {
        this.authLoading = false
      }
    },
    logout() {
      this.user = null
      this.cart = []
      this.orders = []
      localStorage.removeItem('shopUser')
    },
    async loadProducts() {
      try {
        const res = await axios.get(`${API_URL}/products`)
        this.products = res.data || []
      } catch (err) {
        this.showNotification('Failed to load products', 'error')
      }
    },
    async loadOrders() {
      try {
        const res = await axios.get(`${API_URL}/orders?user_id=${this.user.id}`)
        this.orders = (res.data || []).sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
      } catch (err) {
        console.error('Failed to load orders', err)
      }
    },
    addToCart(product) {
      const existing = this.cart.find(i => i.product.id === product.id)
      if (existing) {
        existing.qty++
      } else {
        this.cart.push({ product, qty: 1 })
      }
      this.showNotification(`${product.name} added to cart`, 'success')
      getFaro()?.api.pushEvent('add_to_cart', { product_id: product.id, product_name: product.name })
    },
    updateQty(item, delta) {
      item.qty += delta
      if (item.qty <= 0) this.removeFromCart(item)
    },
    removeFromCart(item) {
      this.cart = this.cart.filter(i => i.product.id !== item.product.id)
    },
    async checkout() {
      this.checkoutLoading = true
      try {
        const orderData = {
          user_id: String(this.user.id),
          items: this.cart.map(i => ({
            product_id: i.product.id,
            name: i.product.name,
            price: i.product.price,
            quantity: i.qty,
          })),
        }
        const orderRes = await axios.post(`${API_URL}/orders`, orderData)
        const order = orderRes.data

        await axios.post(`${API_URL}/payments`, {
          order_id: order.id,
          amount: order.total,
          method: this.paymentMethod,
        })

        await axios.patch(`${API_URL}/orders/${order.id}`, { status: 'paid' })

        getFaro()?.api.pushEvent('checkout_completed', { order_id: order.id, total: String(order.total) })
        this.cart = []
        this.showNotification('Order placed successfully!', 'success')
        this.currentPage = 'orders'
        this.loadOrders()
      } catch (err) {
        this.showNotification('Checkout failed: ' + (err.response?.data?.error || err.message), 'error')
      } finally {
        this.checkoutLoading = false
      }
    },
    showNotification(message, type) {
      this.notification = { message, type }
      setTimeout(() => { this.notification = null }, 3000)
    },
    formatDate(d) {
      return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
    },
    handleImageError(e) {
      e.target.src = 'https://placehold.co/200x200?text=Product'
    },
  },
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }

/* Auth */
.auth-page { display: flex; justify-content: center; align-items: center; min-height: 100vh; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.auth-container { background: white; border-radius: 16px; overflow: hidden; width: 400px; box-shadow: 0 20px 60px rgba(0,0,0,0.3); }
.auth-brand { padding: 40px; text-align: center; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; }
.auth-brand h1 { font-size: 2em; }
.auth-brand p { opacity: 0.9; margin-top: 8px; }
.auth-form { padding: 30px; }
.auth-tabs { display: flex; margin-bottom: 20px; }
.auth-tabs button { flex: 1; padding: 10px; border: none; background: #f0f0f0; cursor: pointer; font-weight: 600; transition: all 0.2s; }
.auth-tabs button.active { background: #667eea; color: white; }
.auth-tabs button:first-child { border-radius: 8px 0 0 8px; }
.auth-tabs button:last-child { border-radius: 0 8px 8px 0; }
.auth-form form { display: flex; flex-direction: column; gap: 12px; }
.auth-form input { padding: 12px 16px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; }
.auth-form input:focus { outline: none; border-color: #667eea; }
.btn-primary { padding: 12px; background: #667eea; color: white; border: none; border-radius: 8px; font-size: 16px; font-weight: 600; cursor: pointer; }
.btn-primary:hover { background: #5a6fd6; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.auth-error { color: #e74c3c; text-align: center; margin-top: 12px; font-size: 14px; }
.auth-hint { color: #999; text-align: center; margin-top: 12px; font-size: 12px; }

/* Topbar */
.topbar { display: flex; align-items: center; padding: 12px 24px; background: white; box-shadow: 0 2px 8px rgba(0,0,0,0.08); position: sticky; top: 0; z-index: 100; }
.topbar-left h1 { font-size: 1.4em; cursor: pointer; color: #667eea; }
.topbar-center { flex: 1; max-width: 400px; margin: 0 24px; }
.search-input { width: 100%; padding: 10px 16px; border: 1px solid #e0e0e0; border-radius: 24px; font-size: 14px; }
.search-input:focus { outline: none; border-color: #667eea; }
.topbar-right { display: flex; align-items: center; gap: 16px; }
.cart-btn, .orders-btn { background: none; border: none; font-size: 16px; cursor: pointer; padding: 8px 12px; border-radius: 8px; position: relative; }
.cart-btn:hover, .orders-btn:hover { background: #f0f0f0; }
.cart-badge { position: absolute; top: 0; right: 0; background: #e74c3c; color: white; border-radius: 50%; width: 18px; height: 18px; font-size: 11px; display: flex; align-items: center; justify-content: center; }
.user-menu { display: flex; align-items: center; gap: 8px; font-size: 14px; }
.logout-btn { background: none; border: 1px solid #ddd; padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.logout-btn:hover { background: #fee; border-color: #e74c3c; color: #e74c3c; }

/* Main */
.main-content { max-width: 1200px; margin: 24px auto; padding: 0 24px; }
.main-content h2 { margin-bottom: 20px; font-size: 1.5em; }

/* Categories */
.categories { display: flex; gap: 8px; margin-bottom: 24px; flex-wrap: wrap; }
.categories button { padding: 8px 16px; border: 1px solid #e0e0e0; border-radius: 20px; background: white; cursor: pointer; font-size: 13px; transition: all 0.2s; }
.categories button.active { background: #667eea; color: white; border-color: #667eea; }
.categories button:hover { border-color: #667eea; }

/* Products Grid */
.products-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 20px; }
.product-card { background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.06); transition: transform 0.2s, box-shadow 0.2s; }
.product-card:hover { transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.12); }
.product-image { height: 180px; background: #f8f8f8; display: flex; align-items: center; justify-content: center; overflow: hidden; }
.product-image img { max-width: 80%; max-height: 80%; object-fit: contain; }
.product-info { padding: 16px; }
.product-category { font-size: 11px; color: #667eea; font-weight: 600; text-transform: uppercase; }
.product-info h3 { font-size: 14px; margin: 6px 0; color: #333; line-height: 1.3; }
.product-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 12px; }
.product-price { font-size: 18px; font-weight: 700; color: #e74c3c; }
.product-stock { font-size: 11px; color: #999; margin-top: 6px; display: block; }
.add-cart-btn { padding: 6px 14px; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 12px; font-weight: 600; }
.add-cart-btn:hover { background: #5a6fd6; }
.add-cart-btn:disabled { background: #ccc; cursor: not-allowed; }

/* Cart */
.cart-layout { display: grid; grid-template-columns: 1fr 350px; gap: 24px; }
.cart-items { display: flex; flex-direction: column; gap: 12px; }
.cart-item { display: flex; align-items: center; gap: 16px; background: white; padding: 16px; border-radius: 12px; }
.cart-item img { width: 60px; height: 60px; object-fit: contain; border-radius: 8px; background: #f8f8f8; }
.cart-item-info { flex: 1; }
.cart-item-info h4 { font-size: 14px; }
.cart-item-price { color: #666; font-size: 13px; }
.cart-item-qty { display: flex; align-items: center; gap: 8px; }
.cart-item-qty button { width: 28px; height: 28px; border: 1px solid #ddd; border-radius: 6px; background: white; cursor: pointer; font-size: 16px; }
.cart-item-qty button:hover { background: #f0f0f0; }
.cart-item-total { font-weight: 700; min-width: 80px; text-align: right; }
.remove-btn { background: none; border: none; color: #e74c3c; cursor: pointer; font-size: 18px; padding: 4px; }
.cart-summary { background: white; padding: 24px; border-radius: 12px; height: fit-content; position: sticky; top: 80px; }
.cart-summary h3 { margin-bottom: 16px; }
.summary-row { display: flex; justify-content: space-between; padding: 8px 0; font-size: 14px; color: #666; }
.summary-row.total { border-top: 2px solid #eee; margin-top: 8px; padding-top: 12px; font-size: 18px; font-weight: 700; color: #333; }
.payment-method { margin-top: 16px; }
.payment-method label { font-size: 13px; color: #666; display: block; margin-bottom: 6px; }
.payment-method select { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; }
.btn-checkout { width: 100%; padding: 14px; background: #27ae60; color: white; border: none; border-radius: 8px; font-size: 16px; font-weight: 700; cursor: pointer; margin-top: 16px; }
.btn-checkout:hover { background: #219a52; }
.btn-checkout:disabled { opacity: 0.6; cursor: not-allowed; }

/* Orders */
.orders-list { display: flex; flex-direction: column; gap: 12px; }
.order-card { background: white; padding: 20px; border-radius: 12px; }
.order-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.order-id { font-weight: 700; font-size: 14px; }
.order-status { padding: 4px 12px; border-radius: 12px; font-size: 12px; font-weight: 600; }
.order-status.pending { background: #fff3cd; color: #856404; }
.order-status.paid { background: #d4edda; color: #155724; }
.order-status.completed { background: #cce5ff; color: #004085; }
.order-items { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 12px; }
.order-item-tag { background: #f0f0f0; padding: 4px 10px; border-radius: 12px; font-size: 12px; }
.order-footer { display: flex; justify-content: space-between; align-items: center; }
.order-total { font-weight: 700; font-size: 16px; color: #e74c3c; }
.order-date { color: #999; font-size: 13px; }

/* Empty & Toast */
.empty-state { text-align: center; padding: 60px; color: #999; }
.empty-state p { margin-bottom: 16px; font-size: 16px; }
.toast { position: fixed; bottom: 24px; right: 24px; padding: 14px 24px; border-radius: 8px; color: white; font-weight: 500; z-index: 999; animation: slideIn 0.3s ease; }
.toast.success { background: #27ae60; }
.toast.error { background: #e74c3c; }
@keyframes slideIn { from { transform: translateX(100%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }

@media (max-width: 768px) {
  .topbar { flex-wrap: wrap; gap: 8px; }
  .topbar-center { order: 3; max-width: 100%; margin: 0; flex-basis: 100%; }
  .cart-layout { grid-template-columns: 1fr; }
  .products-grid { grid-template-columns: repeat(2, 1fr); }
}
</style>
