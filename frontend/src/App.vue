<script>
import axios from 'axios'

const API_URL = 'http://4.144.133.123:8080/api'

export default {
  data() {
    return {
      currentPage: 'products',
      products: [],
      orders: [],
      users: [],
      loading: false
    }
  },
  watch: {
    currentPage(newPage) {
      this.loadData(newPage)
    }
  },
  mounted() {
    this.loadData('products')
  },
  methods: {
    async loadData(page) {
      this.loading = true
      try {
        if (page === 'products') {
          const res = await axios.get(`${API_URL}/products`)
          this.products = res.data || []
        } else if (page === 'orders') {
          const res = await axios.get(`${API_URL}/orders`)
          this.orders = res.data || []
        } else if (page === 'users') {
          const res = await axios.get(`${API_URL}/users`)
          this.users = res.data || []
        }
      } catch (error) {
        console.error(`Error loading ${page}:`, error)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>