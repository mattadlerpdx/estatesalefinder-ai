'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import {
  MapPinIcon,
  CalendarIcon,
  ClockIcon,
  HeartIcon,
  ShareIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  EyeIcon
} from '@heroicons/react/24/outline'
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid'

interface Listing {
  id: number
  seller_id: number
  title: string
  description: string
  sale_type: string
  status: string
  address_line1: string
  address_line2?: string
  city: string
  state: string
  zip_code: string
  latitude?: number
  longitude?: number
  start_date: string
  end_date: string
  listing_tier: string
  view_count: number
  featured: boolean
  driving_directions?: string
  parking_info?: string
  created_at: string
  images?: Array<{
    id: number
    image_url: string
    is_primary: boolean
    display_order: number
  }>
}

export default function SaleDetailPage() {
  const params = useParams()
  const router = useRouter()
  const saleId = params?.id as string

  const [sale, setSale] = useState<Listing | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [currentImageIndex, setCurrentImageIndex] = useState(0)
  const [isSaved, setIsSaved] = useState(false)

  useEffect(() => {
    if (saleId) {
      fetchSale()
    }
  }, [saleId])

  const fetchSale = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/api/sales/${saleId}`)

      if (!response.ok) {
        if (response.status === 404) {
          setError('Sale not found')
        } else {
          setError('Failed to load sale')
        }
        return
      }

      const data = await response.json()
      setSale(data)
    } catch (err) {
      console.error('Error fetching sale:', err)
      setError('Failed to load sale')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      weekday: 'long',
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    })
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
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
      <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${style.bg} ${style.text}`}>
        {saleType.replace('_', ' ').toUpperCase()}
      </span>
    )
  }

  const nextImage = () => {
    if (sale?.images && sale.images.length > 0) {
      setCurrentImageIndex((prev) => (prev + 1) % sale.images!.length)
    }
  }

  const prevImage = () => {
    if (sale?.images && sale.images.length > 0) {
      setCurrentImageIndex((prev) => (prev - 1 + sale.images!.length) % sale.images!.length)
    }
  }

  const handleSave = () => {
    // TODO: Implement save functionality with authentication
    setIsSaved(!isSaved)
  }

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: sale?.title,
          text: `Check out this ${sale?.sale_type.replace('_', ' ')} in ${sale?.city}, ${sale?.state}`,
          url: window.location.href
        })
      } catch (err) {
        console.log('Error sharing:', err)
      }
    } else {
      // Fallback: Copy to clipboard
      navigator.clipboard.writeText(window.location.href)
      alert('Link copied to clipboard!')
    }
  }

  const getDirectionsUrl = () => {
    if (sale?.latitude && sale?.longitude) {
      return `https://www.google.com/maps/dir/?api=1&destination=${sale.latitude},${sale.longitude}`
    } else {
      const address = `${sale?.address_line1}, ${sale?.city}, ${sale?.state} ${sale?.zip_code}`
      return `https://www.google.com/maps/dir/?api=1&destination=${encodeURIComponent(address)}`
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error || !sale) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">{error || 'Sale not found'}</h2>
          <p className="text-gray-600 mb-6">The sale you're looking for doesn't exist or has been removed.</p>
          <Link
            href="/sales"
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
          >
            Browse All Sales
          </Link>
        </div>
      </div>
    )
  }

  const sortedImages = sale.images?.sort((a, b) => {
    if (a.is_primary) return -1
    if (b.is_primary) return 1
    return a.display_order - b.display_order
  }) || []

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Back Button */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <Link
            href="/sales"
            className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
          >
            <ChevronLeftIcon className="h-4 w-4 mr-1" />
            Back to all sales
          </Link>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2">
            {/* Image Gallery */}
            <div className="bg-white rounded-lg shadow-sm overflow-hidden mb-6">
              {sortedImages.length > 0 ? (
                <div className="relative">
                  <div className="aspect-w-16 aspect-h-9 bg-gray-900">
                    <img
                      src={sortedImages[currentImageIndex].image_url}
                      alt={`${sale.title} - Image ${currentImageIndex + 1}`}
                      className="w-full h-96 object-contain"
                    />
                  </div>

                  {sortedImages.length > 1 && (
                    <>
                      <button
                        onClick={prevImage}
                        className="absolute left-4 top-1/2 -translate-y-1/2 bg-white/90 hover:bg-white p-2 rounded-full shadow-lg"
                      >
                        <ChevronLeftIcon className="h-6 w-6 text-gray-900" />
                      </button>
                      <button
                        onClick={nextImage}
                        className="absolute right-4 top-1/2 -translate-y-1/2 bg-white/90 hover:bg-white p-2 rounded-full shadow-lg"
                      >
                        <ChevronRightIcon className="h-6 w-6 text-gray-900" />
                      </button>
                      <div className="absolute bottom-4 left-1/2 -translate-x-1/2 bg-black/60 text-white px-3 py-1 rounded-full text-sm">
                        {currentImageIndex + 1} / {sortedImages.length}
                      </div>
                    </>
                  )}
                </div>
              ) : (
                <div className="h-96 bg-gradient-to-br from-blue-100 to-blue-50 flex items-center justify-center">
                  <MapPinIcon className="h-24 w-24 text-blue-300" />
                </div>
              )}

              {/* Thumbnail Strip */}
              {sortedImages.length > 1 && (
                <div className="flex gap-2 p-4 overflow-x-auto">
                  {sortedImages.map((image, index) => (
                    <button
                      key={image.id}
                      onClick={() => setCurrentImageIndex(index)}
                      className={`flex-shrink-0 w-20 h-20 rounded-md overflow-hidden border-2 transition-all ${
                        index === currentImageIndex
                          ? 'border-blue-600 ring-2 ring-blue-200'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <img
                        src={image.image_url}
                        alt={`Thumbnail ${index + 1}`}
                        className="w-full h-full object-cover"
                      />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Sale Details */}
            <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
              <div className="flex items-start justify-between mb-4">
                <div>
                  {getSaleTypeBadge(sale.sale_type)}
                  {sale.featured && (
                    <span className="ml-2 inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800">
                      FEATURED
                    </span>
                  )}
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={handleSave}
                    className="p-2 rounded-full hover:bg-gray-100 transition-colors"
                  >
                    {isSaved ? (
                      <HeartSolidIcon className="h-6 w-6 text-red-500" />
                    ) : (
                      <HeartIcon className="h-6 w-6 text-gray-600" />
                    )}
                  </button>
                  <button
                    onClick={handleShare}
                    className="p-2 rounded-full hover:bg-gray-100 transition-colors"
                  >
                    <ShareIcon className="h-6 w-6 text-gray-600" />
                  </button>
                </div>
              </div>

              <h1 className="text-3xl font-bold text-gray-900 mb-4">{sale.title}</h1>

              <div className="flex items-center text-sm text-gray-500 mb-6">
                <EyeIcon className="h-4 w-4 mr-1" />
                <span>{sale.view_count} views</span>
              </div>

              <div className="prose max-w-none">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">Description</h3>
                <p className="text-gray-700 whitespace-pre-wrap">{sale.description}</p>
              </div>

              {sale.driving_directions && (
                <div className="mt-6 pt-6 border-t">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">Driving Directions</h3>
                  <p className="text-gray-700 whitespace-pre-wrap">{sale.driving_directions}</p>
                </div>
              )}

              {sale.parking_info && (
                <div className="mt-6 pt-6 border-t">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">Parking Information</h3>
                  <p className="text-gray-700 whitespace-pre-wrap">{sale.parking_info}</p>
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow-sm p-6 sticky top-4">
              <h2 className="text-xl font-bold text-gray-900 mb-6">Sale Information</h2>

              <div className="space-y-6">
                {/* Date & Time */}
                <div>
                  <div className="flex items-start mb-3">
                    <CalendarIcon className="h-5 w-5 text-gray-400 mr-3 mt-0.5 flex-shrink-0" />
                    <div>
                      <p className="text-sm font-medium text-gray-900">Date</p>
                      <p className="text-sm text-gray-600">{formatDate(sale.start_date)}</p>
                      {sale.start_date !== sale.end_date && (
                        <p className="text-sm text-gray-600">to {formatDate(sale.end_date)}</p>
                      )}
                    </div>
                  </div>

                  <div className="flex items-start">
                    <ClockIcon className="h-5 w-5 text-gray-400 mr-3 mt-0.5 flex-shrink-0" />
                    <div>
                      <p className="text-sm font-medium text-gray-900">Time</p>
                      <p className="text-sm text-gray-600">
                        {formatTime(sale.start_date)} - {formatTime(sale.end_date)}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Location */}
                <div className="pt-6 border-t">
                  <div className="flex items-start mb-3">
                    <MapPinIcon className="h-5 w-5 text-gray-400 mr-3 mt-0.5 flex-shrink-0" />
                    <div>
                      <p className="text-sm font-medium text-gray-900">Location</p>
                      <p className="text-sm text-gray-600">{sale.address_line1}</p>
                      {sale.address_line2 && (
                        <p className="text-sm text-gray-600">{sale.address_line2}</p>
                      )}
                      <p className="text-sm text-gray-600">
                        {sale.city}, {sale.state} {sale.zip_code}
                      </p>
                    </div>
                  </div>

                  <a
                    href={getDirectionsUrl()}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center justify-center w-full px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  >
                    Get Directions
                  </a>
                </div>

                {/* Contact Seller */}
                <div className="pt-6 border-t">
                  <button
                    className="w-full px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    onClick={() => alert('Contact feature coming soon! This will require authentication.')}
                  >
                    Contact Seller
                  </button>
                </div>

                {/* Report */}
                <div className="pt-6 border-t">
                  <button
                    className="text-sm text-gray-500 hover:text-gray-700"
                    onClick={() => alert('Report feature coming soon!')}
                  >
                    Report this listing
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
