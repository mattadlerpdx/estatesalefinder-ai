'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { MagnifyingGlassIcon, FunnelIcon, MapPinIcon, CalendarIcon } from '@heroicons/react/24/outline'

interface Listing {
  id: string
  title: string
  description?: string
  address: string
  city: string
  state: string
  zip_code: string
  start_date: string
  end_date: string
  thumbnail_url?: string
  images?: Array<{
    image_url: string
    is_primary: boolean
  }>
  is_scraped: boolean
  source?: {
    name: string
    url: string
  }
  sale_type?: string
  status?: string
  view_count?: number
  featured?: boolean
}

export default function SalesPage() {
  const [sales, setSales] = useState<Listing[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [filters, setFilters] = useState({
    city: '',
    state: '',
    saleType: '',
    featured: false
  })
  const [showFilters, setShowFilters] = useState(false)

  useEffect(() => {
    fetchSales()
  }, [])

  const fetchSales = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      if (filters.city) params.append('city', filters.city)
      if (filters.state) params.append('state', filters.state)
      if (filters.saleType) params.append('sale_type', filters.saleType)
      if (filters.featured) params.append('featured', 'true')

      const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
      const url = `${backendUrl}/api/sales${params.toString() ? '?' + params.toString() : ''}`
      const response = await fetch(url)

      if (!response.ok) {
        throw new Error('Failed to fetch sales')
      }

      const data = await response.json()
      setSales(data.sales || [])
    } catch (error) {
      console.error('Error fetching sales:', error)
      setSales([])
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    fetchSales()
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    })
  }

  const getSaleTypeBadge = (saleType: string) => {
    const types: Record<string, { bg: string, text: string }> = {
      estate_sale: { bg: 'bg-blue-100', text: 'text-blue-800' },
      moving_sale: { bg: 'bg-green-100', text: 'text-green-800' },
      auction: { bg: 'bg-purple-100', text: 'text-purple-800' },
      garage_sale: { bg: 'bg-yellow-100', text: 'text-yellow-800' }
    }
    const style = types[saleType] || { bg: 'bg-gray-100', text: 'text-gray-800' }
    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
        {saleType.replace('_', ' ').toUpperCase()}
      </span>
    )
  }

  const filteredSales = (Array.isArray(sales) ? sales : []).filter(sale => {
    if (!searchQuery) return true
    const query = searchQuery.toLowerCase()
    return (
      sale.title.toLowerCase().includes(query) ||
      sale.city.toLowerCase().includes(query) ||
      sale.state.toLowerCase().includes(query) ||
      (sale.description && sale.description.toLowerCase().includes(query))
    )
  })

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">
      {/* Header */}
      <div className="bg-white dark:bg-gray-800 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Estate Sales Near You</h1>
          <p className="mt-2 text-gray-600 dark:text-gray-300">Discover amazing deals on estate sales across the country</p>
        </div>
      </div>

      {/* Search and Filters */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
          <form onSubmit={handleSearch} className="space-y-4">
            {/* Search Bar */}
            <div className="flex gap-4">
              <div className="flex-1 relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type="text"
                  placeholder="Search by title, city, or description..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="block w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md leading-5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>
              <button
                type="button"
                onClick={() => setShowFilters(!showFilters)}
                className="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
              >
                <FunnelIcon className="h-5 w-5 mr-2" />
                Filters
              </button>
            </div>

            {/* Filter Options */}
            {showFilters && (
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 pt-4 border-t dark:border-gray-700">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">City</label>
                  <input
                    type="text"
                    value={filters.city}
                    onChange={(e) => setFilters({ ...filters, city: e.target.value })}
                    placeholder="Enter city"
                    className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">State</label>
                  <input
                    type="text"
                    value={filters.state}
                    onChange={(e) => setFilters({ ...filters, state: e.target.value })}
                    placeholder="e.g., OR"
                    className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Sale Type</label>
                  <select
                    value={filters.saleType}
                    onChange={(e) => setFilters({ ...filters, saleType: e.target.value })}
                    className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  >
                    <option value="">All Types</option>
                    <option value="estate_sale">Estate Sale</option>
                    <option value="moving_sale">Moving Sale</option>
                    <option value="auction">Auction</option>
                    <option value="garage_sale">Garage Sale</option>
                  </select>
                </div>
                <div className="flex items-end">
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      checked={filters.featured}
                      onChange={(e) => setFilters({ ...filters, featured: e.target.checked })}
                      className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 rounded"
                    />
                    <span className="ml-2 text-sm text-gray-700 dark:text-gray-300">Featured only</span>
                  </label>
                </div>
              </div>
            )}

            <div className="flex gap-2">
              <button
                type="submit"
                onClick={() => fetchSales()}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Apply Filters
              </button>
              {(filters.city || filters.state || filters.saleType || filters.featured) && (
                <button
                  type="button"
                  onClick={() => {
                    setFilters({ city: '', state: '', saleType: '', featured: false })
                    setTimeout(fetchSales, 100)
                  }}
                  className="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
                >
                  Clear Filters
                </button>
              )}
            </div>
          </form>
        </div>
      </div>

      {/* Sales Grid */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-12">
        {loading ? (
          <div className="flex justify-center items-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredSales.length === 0 ? (
          <div className="text-center py-12">
            <MagnifyingGlassIcon className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">No sales found</h3>
            <p className="mt-1 text-sm text-gray-500">
              {searchQuery || filters.city || filters.state || filters.saleType
                ? "Try adjusting your search or filters"
                : "Be the first to list a sale!"}
            </p>
            <div className="mt-6">
              <Link
                href="/dashboard"
                className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                List Your Sale
              </Link>
            </div>
          </div>
        ) : (
          <>
            <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
              Found {filteredSales.length} {filteredSales.length === 1 ? 'sale' : 'sales'}
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredSales.map((sale) => {
                const linkTarget = sale.is_scraped && sale.source ? sale.source.url : `/sales/${sale.id}`
                const linkProps = sale.is_scraped && sale.source
                  ? { target: "_blank", rel: "noopener noreferrer" }
                  : {}

                return (
                  <Link
                    key={sale.id}
                    href={linkTarget}
                    {...linkProps}
                    className="bg-white dark:bg-gray-800 rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden group"
                  >
                    {/* Image */}
                    <div className="aspect-w-16 aspect-h-9 bg-gray-200 relative h-48">
                      {sale.images && sale.images.length > 0 ? (
                        <img
                          src={sale.images.find(img => img.is_primary)?.image_url || sale.images[0].image_url}
                          alt={sale.title}
                          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
                        />
                      ) : sale.thumbnail_url ? (
                        <img
                          src={sale.thumbnail_url}
                          alt={sale.title}
                          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
                        />
                      ) : (
                        <div className="flex items-center justify-center h-full bg-gradient-to-br from-blue-100 to-blue-50 dark:from-gray-700 dark:to-gray-800">
                          <MapPinIcon className="h-16 w-16 text-blue-300 dark:text-gray-600" />
                        </div>
                      )}
                      {sale.featured && (
                        <div className="absolute top-2 right-2">
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-400 text-yellow-900">
                            FEATURED
                          </span>
                        </div>
                      )}
                      {sale.is_scraped && sale.source && (
                        <div className="absolute top-2 left-2">
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100">
                            via {sale.source.name}
                          </span>
                        </div>
                      )}
                    </div>

                    {/* Content */}
                    <div className="p-5">
                      <div className="mb-2">
                        {sale.sale_type && getSaleTypeBadge(sale.sale_type)}
                      </div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2 line-clamp-2 group-hover:text-blue-600 dark:group-hover:text-blue-400">
                        {sale.title}
                      </h3>
                      {sale.description && (
                        <p className="text-sm text-gray-600 dark:text-gray-300 mb-3 line-clamp-2">
                          {sale.description}
                        </p>
                      )}

                      <div className="space-y-2 text-sm text-gray-500 dark:text-gray-400">
                        <div className="flex items-start">
                          <MapPinIcon className="h-4 w-4 mr-2 flex-shrink-0 mt-0.5" />
                          <div className="flex-1 min-w-0">
                            {sale.address && sale.address !== 'TBA' && (
                              <div className="truncate text-gray-900 dark:text-gray-100">{sale.address}</div>
                            )}
                            <div className="truncate">
                              {sale.city}, {sale.state} {sale.zip_code}
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center">
                          <CalendarIcon className="h-4 w-4 mr-2 flex-shrink-0" />
                          <span className="truncate">
                            {formatDate(sale.start_date)}
                          </span>
                        </div>
                      </div>

                      <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                        {sale.is_scraped ? (
                          <span className="text-green-600 dark:text-green-400">Listed on {sale.source?.name}</span>
                        ) : (
                          <span>{sale.view_count || 0} views</span>
                        )}
                        <span className="text-blue-600 dark:text-blue-400 group-hover:text-blue-700 dark:group-hover:text-blue-300 font-medium">
                          {sale.is_scraped ? 'View on Site →' : 'View Details →'}
                        </span>
                      </div>
                    </div>
                  </Link>
                )
              })}
            </div>
          </>
        )}
      </div>
    </div>
  )
}
